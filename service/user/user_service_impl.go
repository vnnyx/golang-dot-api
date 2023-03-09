package user

import (
	"bytes"
	"context"
	"embed"
	"encoding/json"
	"errors"
	"math/rand"
	"strconv"
	"sync"
	"text/template"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/vnnyx/golang-dot-api/model/entity"
	"github.com/vnnyx/golang-dot-api/model/message"
	"github.com/vnnyx/golang-dot-api/model/web"
	"github.com/vnnyx/golang-dot-api/repository/transaction"
	"github.com/vnnyx/golang-dot-api/repository/user"
	"github.com/vnnyx/golang-dot-api/util"
	"github.com/vnnyx/golang-dot-api/validation"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserServiceImpl struct {
	*kafka.Producer
	user.UserRepository
	transaction.TransactionRepository
	*gorm.DB
}

func NewUserService(userRepository user.UserRepository, transactionRepository transaction.TransactionRepository, DB *gorm.DB, producer *kafka.Producer) UserService {
	return &UserServiceImpl{UserRepository: userRepository, TransactionRepository: transactionRepository, DB: DB, Producer: producer}
}

func (service *UserServiceImpl) CreateUser(ctx context.Context, request web.UserCreateRequest) (response web.UserEmailVerification, err error) {
	logger := logrus.New()
	validation.CreateUserValidation(request)

	if request.Password != request.PasswordConfirmation {
		return response, errors.New("PASSWORD_NOT_MATCH")
	}

	password, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Error(err)
		return response, err
	}

	user := entity.User{
		UserID:    uuid.NewString(),
		Username:  request.Username,
		Email:     request.Email,
		Handphone: request.Handphone,
		Password:  string(password),
	}

	response = web.UserEmailVerification{UserID: user.UserID, ExpiredAt: time.Until(time.Now().Add(5 * time.Minute))}

	payload := message.Message{
		User: user,
		OTP:  &response,
	}
	encodedPayload, err := json.Marshal(payload)
	if err != nil {
		logger.Error(err)
		return response, err
	}

	service.Producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &message.USER_OTP_TOPIC, Partition: kafka.PartitionAny},
		Value:          encodedPayload,
	}, nil)

	return response, nil
}

func (service *UserServiceImpl) GetUserById(ctx context.Context, userId string) (response web.UserResponse, err error) {
	logger := logrus.New()
	user, err := service.UserRepository.FindUserByID(ctx, userId)
	if err != nil {
		logger.Error(err)
		return response, errors.New("USER_NOT_FOUND")
	}

	response = web.UserResponse{
		UserID:    user.UserID,
		Username:  user.Username,
		Email:     user.Email,
		Handphone: user.Handphone,
	}

	return response, nil
}

func (service *UserServiceImpl) GetAllUser(ctx context.Context, p web.Pagination) (response *web.Pagination, err error) {
	logger := logrus.New()
	users, err := service.UserRepository.FindAllUser(ctx, &p)
	if err != nil {
		logger.Error(err)
		return response, err
	}

	rows := make([]web.UserResponse, 0)
	for _, user := range users {
		rows = append(rows, web.UserResponse{
			UserID:    user.UserID,
			Username:  user.Username,
			Email:     user.Email,
			Handphone: user.Handphone,
		})
	}

	return &web.Pagination{
		Limit:      p.GetLimit(),
		Page:       p.GetPage(),
		Sort:       p.GetSort(),
		TotalRows:  p.GetTotalRows(),
		TotalPages: p.GetTotalPages(),
		Rows:       rows,
	}, nil
}

func (service *UserServiceImpl) UpdateUserProfile(ctx context.Context, request web.UserUpdateProfileRequest) (response web.UserResponse, err error) {
	logger := logrus.New()
	validation.UpdateUserProfileValidation(request)

	user, err := service.UserRepository.FindUserByID(ctx, request.UserID)
	if err != nil {
		logger.Error(err)
		return response, errors.New("USER_NOT_FOUND")
	}

	user, err = service.UserRepository.UpdateUser(ctx, entity.User{
		UserID:    user.UserID,
		Username:  request.Username,
		Email:     request.Email,
		Handphone: request.Handphone,
	})

	if err != nil {
		logger.Error(err)
		return response, err
	}

	response = web.UserResponse{
		UserID:    user.UserID,
		Username:  user.Username,
		Email:     user.Email,
		Handphone: user.Handphone,
	}

	return response, nil
}

func (service *UserServiceImpl) RemoveUser(ctx context.Context, userId string) error {
	logger := logrus.New()
	user, err := service.UserRepository.FindUserByID(ctx, userId)
	if err != nil {
		logger.Error(err)
		return errors.New("USER_NOT_FOUND")
	}

	tx := service.DB.Begin()
	err = tx.Error
	if err != nil {
		logger.Error(err)
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit().Error
		}
	}()

	err = service.TransactionRepository.DeleteTransactionByUserId(ctx, tx, user.UserID)
	if err != nil {
		logger.Error(err)
		return err
	}

	err = service.UserRepository.DeleteUser(ctx, tx, user.UserID)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (service *UserServiceImpl) GetAllUserWithLastTransaction(ctx context.Context, p web.Pagination) (response *web.Pagination, err error) {
	var wg sync.WaitGroup
	chUser := make(chan web.UserResponse, 100)
	chUserWithLastTransaction := make(chan web.UserResponseWithLastTransaction, 100)
	wg.Add(2)
	go func() {
		var u web.UserResponse
		defer wg.Done()
		users, _ := service.UserRepository.FindAllUser(ctx, &p)
		for _, user := range users {
			u = web.UserResponse{
				UserID:    user.UserID,
				Username:  user.Username,
				Email:     user.Email,
				Handphone: user.Handphone,
			}
			chUser <- u
		}
		close(chUser)
	}()
	go func() {
		defer wg.Done()
		for user := range chUser {
			transactions, _ := service.TransactionRepository.FindTransactionByUserId(ctx, user.UserID)
			if len(transactions) > 0 {
				t := web.TransactionResponse{
					TransactionID: transactions[0].TransactionID,
					Name:          transactions[0].Name,
					UserID:        transactions[0].UserID,
				}
				ut := web.UserResponseWithLastTransaction{
					UserID:      user.UserID,
					Username:    user.Username,
					Email:       user.Email,
					Handphone:   user.Handphone,
					Transaction: t,
				}
				chUserWithLastTransaction <- ut
			}
		}
		close(chUserWithLastTransaction)
	}()
	wg.Wait()
	rows := make([]web.UserResponseWithLastTransaction, 0)
	for data := range chUserWithLastTransaction {
		rows = append(rows, web.UserResponseWithLastTransaction{
			UserID:      data.UserID,
			Username:    data.Username,
			Email:       data.Email,
			Handphone:   data.Handphone,
			Transaction: data.Transaction,
		})
	}
	return &web.Pagination{
		Limit:      p.GetLimit(),
		Page:       p.GetPage(),
		Sort:       p.GetSort(),
		TotalRows:  p.GetTotalRows(),
		TotalPages: p.GetTotalPages(),
		Rows:       rows,
	}, nil
}

//go:embed templates/*.gohtml
var templates embed.FS

func (service *UserServiceImpl) SendOTP(ctx context.Context, user entity.User, verify *web.UserEmailVerification) error {
	logger := logrus.New()
	OTP := rand.Intn(9999-1000) + 1000
	verify.UserID = user.UserID
	verify.OTP = OTP
	verify.ExpiredAt = time.Until(time.Now().Add(5 * time.Minute))

	err := service.UserRepository.StoreToRedis(ctx, *verify, user)
	if err != nil {
		logger.Error(err)
		return err
	}

	t, err := template.ParseFS(templates, "templates/*.gohtml")
	if err != nil {
		logger.Error(err)
		return err
	}
	buff := new(bytes.Buffer)
	err = t.ExecuteTemplate(buff, "otp.gohtml", map[string]interface{}{
		"Username": user.Username,
		"Otp":      OTP,
		"UserID":   user.UserID,
	})
	if err != nil {
		logger.Error(err)
		return err
	}
	err = util.SendEmailTo(user.Email, buff)
	if err != nil {
		return err
	}
	return nil
}

func (service *UserServiceImpl) ValidateOTP(ctx context.Context, check web.UserEmailVerification) (response web.UserResponse, err error) {
	logger := logrus.New()
	otp, got, err := service.UserRepository.GetDataToVerify(ctx, check.UserID)
	if err != nil {
		logger.Error(err)
		return response, errors.New("FAILED_TO_VERIFY")
	}
	if strconv.Itoa(check.OTP) == otp {
		user, err := service.UserRepository.InsertUser(ctx, got)
		response = web.UserResponse{
			UserID:    user.UserID,
			Username:  user.Username,
			Email:     user.Email,
			Handphone: user.Handphone,
		}
		if err != nil {
			logger.Error(err)
			return response, err
		}
		err = service.UserRepository.DeleteCache(ctx, user.UserID)
		if err != nil {
			logger.Error(err)
			return response, err
		}
		return response, nil
	}
	return response, errors.New("FAILED_TO_VERIFY")
}
