package transaction

import (
	"github.com/vnnyx/golang-dot-api/model/entity"
	"gorm.io/gorm"
)

type TransactionRepository interface {
	InsertTransaction(transaction entity.Transaction) (entity.Transaction, error)
	FindTransactionByID(transactionId string) (transaction entity.Transaction, err error)
	FindAllTransaction() (transactions []entity.Transaction, err error)
	FindTransactionByUserId(userId string) (transactions []entity.Transaction, err error)
	UpdateTransaction(transaction entity.Transaction) (entity.Transaction, error)
	DeleteTransaction(transactionId string) error
	DeleteTransactionByUserId(tx *gorm.DB, userId string) error
}
