package transaction

import (
	"context"

	"github.com/vnnyx/golang-dot-api/model/web"
)

type TransactionService interface {
	CreateTransaction(ctx context.Context, request web.TransactionCreateRequest) (response web.TransactionResponse, err error)
	GetTransactionById(ctx context.Context, transactionId string) (response web.TransactionResponse, err error)
	GetAllTransaction(ctx context.Context) (response []web.TransactionResponse, err error)
	GetTransactionByUserId(ctx context.Context, userId string) (response []web.TransactionResponse, err error)
	UpdateTransaction(ctx context.Context, request web.TransactionUpdateRequest) (response web.TransactionResponse, err error)
	RemoveTransaction(ctx context.Context, transactionId string) error
}
