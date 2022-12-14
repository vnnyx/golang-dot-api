package transaction

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/vnnyx/golang-dot-api/model/entity"
	"github.com/vnnyx/golang-dot-api/model/web"
	"github.com/vnnyx/golang-dot-api/repository/transaction"
	"github.com/vnnyx/golang-dot-api/repository/user"
	"github.com/vnnyx/golang-dot-api/validation"
)

type TransactionServiceImpl struct {
	transaction.TransactionRepository
	user.UserRepository
}

func NewTransactionService(transactionRepository transaction.TransactionRepository, userRepository user.UserRepository) TransactionService {
	return &TransactionServiceImpl{TransactionRepository: transactionRepository, UserRepository: userRepository}
}

func (service *TransactionServiceImpl) CreateTransaction(ctx context.Context, request web.TransactionCreateRequest) (response web.TransactionResponse, err error) {
	validation.CreateTransactionValidation(request)

	user, err := service.UserRepository.FindUserByID(ctx, request.UserID)
	if err != nil {
		return response, errors.New("USER_NOT_FOUND")
	}

	transaction, err := service.TransactionRepository.InsertTransaction(ctx, entity.Transaction{
		TransactionID: uuid.NewString(),
		Name:          request.Name,
		UserID:        user.UserID,
	})

	if err != nil {
		return response, err
	}

	response = web.TransactionResponse{
		TransactionID: transaction.TransactionID,
		Name:          transaction.Name,
		UserID:        transaction.UserID,
	}

	return response, nil
}

func (service *TransactionServiceImpl) GetTransactionById(ctx context.Context, transactionId string) (response web.TransactionResponse, err error) {
	transaction, err := service.TransactionRepository.FindTransactionByID(ctx, transactionId)
	if err != nil {
		return response, errors.New("TRANSACTION_NOT_FOUND")
	}

	response = web.TransactionResponse{
		TransactionID: transaction.TransactionID,
		Name:          transaction.Name,
		UserID:        transaction.UserID,
	}

	return response, nil
}

func (service *TransactionServiceImpl) GetAllTransaction(ctx context.Context) (response []web.TransactionResponse, err error) {
	transactions, err := service.TransactionRepository.FindAllTransaction(ctx)
	if err != nil {
		return response, err
	}

	for _, transaction := range transactions {
		response = append(response, web.TransactionResponse{
			TransactionID: transaction.TransactionID,
			Name:          transaction.Name,
			UserID:        transaction.UserID,
		})
	}

	return response, nil
}

func (service *TransactionServiceImpl) GetTransactionByUserId(ctx context.Context, userId string) (response []web.TransactionResponse, err error) {
	transactions, err := service.TransactionRepository.FindTransactionByUserId(ctx, userId)
	if err != nil {
		return response, err
	}

	for _, transaction := range transactions {
		response = append(response, web.TransactionResponse{
			TransactionID: transaction.TransactionID,
			Name:          transaction.Name,
			UserID:        transaction.UserID,
		})
	}

	return response, nil
}

func (service *TransactionServiceImpl) UpdateTransaction(ctx context.Context, request web.TransactionUpdateRequest) (response web.TransactionResponse, err error) {
	validation.UpdateTransactionValidation(request)

	transaction, err := service.TransactionRepository.FindTransactionByID(ctx, request.TransactionID)
	if err != nil {
		return response, errors.New("TRANSACTION_NOT_FOUND")
	}

	transaction, err = service.TransactionRepository.UpdateTransaction(ctx, entity.Transaction{
		TransactionID: transaction.TransactionID,
		Name:          request.Name,
		UserID:        transaction.UserID,
	})

	if err != nil {
		return response, err
	}

	response = web.TransactionResponse{
		TransactionID: transaction.TransactionID,
		Name:          transaction.Name,
		UserID:        transaction.UserID,
	}

	return response, nil
}

func (service *TransactionServiceImpl) RemoveTransaction(ctx context.Context, transactionId string) error {
	transaction, err := service.TransactionRepository.FindTransactionByID(ctx, transactionId)
	if err != nil {
		return errors.New("TRANSACTION_NOT_FOUND")
	}
	return service.TransactionRepository.DeleteTransaction(ctx, transaction.TransactionID)
}
