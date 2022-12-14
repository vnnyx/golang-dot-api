package transaction

import (
	"errors"

	"github.com/google/uuid"
	"github.com/vnnyx/golang-dot-api/model/entity"
	"github.com/vnnyx/golang-dot-api/model/web"
	"github.com/vnnyx/golang-dot-api/repository/transaction"
	"github.com/vnnyx/golang-dot-api/repository/user"
)

type TransactionServiceImpl struct {
	transaction.TransactionRepository
	user.UserRepository
}

func NewTransactionService(transactionRepository transaction.TransactionRepository, userRepository user.UserRepository) TransactionService {
	return &TransactionServiceImpl{TransactionRepository: transactionRepository, UserRepository: userRepository}
}

func (service *TransactionServiceImpl) CreateTransaction(request web.TransactionCreateRequest) (response web.TransactionResponse, err error) {
	user, err := service.UserRepository.FindUserByID(request.UserID)
	if err != nil {
		return response, errors.New("USER_NOT_FOUND")
	}

	transaction, err := service.TransactionRepository.InsertTransaction(entity.Transaction{
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

func (service *TransactionServiceImpl) GetTransactionById(transactionId string) (response web.TransactionResponse, err error) {
	transaction, err := service.TransactionRepository.FindTransactionByID(transactionId)
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

func (service *TransactionServiceImpl) GetAllTransaction() (response []web.TransactionResponse, err error) {
	transactions, err := service.TransactionRepository.FindAllTransaction()
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

func (service *TransactionServiceImpl) GetTransactionByUserId(userId string) (response []web.TransactionResponse, err error) {
	transactions, err := service.TransactionRepository.FindTransactionByUserId(userId)
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

func (service *TransactionServiceImpl) UpdateTransaction(request web.TransactionUpdateRequest) (response web.TransactionResponse, err error) {
	transaction, err := service.TransactionRepository.FindTransactionByID(request.TransactionID)
	if err != nil {
		return response, errors.New("TRANSACTION_NOT_FOUND")
	}

	transaction, err = service.TransactionRepository.UpdateTransaction(entity.Transaction{
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

func (service *TransactionServiceImpl) RemoveTransaction(transactionId string) error {
	transaction, err := service.TransactionRepository.FindTransactionByID(transactionId)
	if err != nil {
		return errors.New("TRANSACTION_NOT_FOUND")
	}
	return service.TransactionRepository.DeleteTransaction(transaction.TransactionID)
}
