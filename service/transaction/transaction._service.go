package transaction

import "github.com/vnnyx/golang-dot-api/model/web"

type TransactionService interface {
	CreateTransaction(request web.TransactionCreateRequest) (response web.TransactionResponse, err error)
	GetTransactionById(transactionId string) (response web.TransactionResponse, err error)
	GetAllTransaction() (response []web.TransactionResponse, err error)
	GetTransactionByUserId(userId string) (response []web.TransactionResponse, err error)
	UpdateTransaction(request web.TransactionUpdateRequest) (response web.TransactionResponse, err error)
	RemoveTransaction(transactionId string) error
}
