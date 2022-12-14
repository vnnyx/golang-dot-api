package transaction

import (
	"context"

	"github.com/vnnyx/golang-dot-api/model/entity"
	"gorm.io/gorm"
)

type TransactionRepository interface {
	InsertTransaction(ctx context.Context, transaction entity.Transaction) (entity.Transaction, error)
	FindTransactionByID(ctx context.Context, transactionId string) (transaction entity.Transaction, err error)
	FindAllTransaction(ctx context.Context) (transactions []entity.Transaction, err error)
	FindTransactionByUserId(ctx context.Context, userId string) (transactions []entity.Transaction, err error)
	UpdateTransaction(ctx context.Context, transaction entity.Transaction) (entity.Transaction, error)
	DeleteTransaction(ctx context.Context, transactionId string) error
	DeleteTransactionByUserId(ctx context.Context, tx *gorm.DB, userId string) error
	DeleteAllTransaction(ctx context.Context) error
}
