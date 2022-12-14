package transaction

import (
	"context"

	"github.com/vnnyx/golang-dot-api/model/entity"
	"gorm.io/gorm"
)

type TransactionRepositoryImpl struct {
	*gorm.DB
}

func NewTransactionRepository(DB *gorm.DB) TransactionRepository {
	return &TransactionRepositoryImpl{DB: DB}
}

func (repository *TransactionRepositoryImpl) InsertTransaction(ctx context.Context, transaction entity.Transaction) (entity.Transaction, error) {
	err := repository.DB.WithContext(ctx).Create(&transaction).Error
	return transaction, err
}

func (repository *TransactionRepositoryImpl) FindTransactionByID(ctx context.Context, transactionId string) (transaction entity.Transaction, err error) {
	err = repository.DB.WithContext(ctx).Where("transaction_id", transactionId).First(&transaction).Error
	return transaction, err
}

func (repository *TransactionRepositoryImpl) FindAllTransaction(ctx context.Context) (transactions []entity.Transaction, err error) {
	err = repository.DB.WithContext(ctx).Find(&transactions).Error
	return transactions, err
}

func (repository *TransactionRepositoryImpl) FindTransactionByUserId(ctx context.Context, userId string) (transactions []entity.Transaction, err error) {
	err = repository.DB.WithContext(ctx).Where("user_id", userId).Find(&transactions).Error
	return transactions, err
}

func (repository *TransactionRepositoryImpl) UpdateTransaction(ctx context.Context, transaction entity.Transaction) (entity.Transaction, error) {
	err := repository.DB.WithContext(ctx).Where("transaction_id", transaction.TransactionID).Updates(&transaction).Error
	return transaction, err
}

func (repository *TransactionRepositoryImpl) DeleteTransaction(ctx context.Context, transactionId string) error {
	return repository.DB.WithContext(ctx).Where("transaction_id", transactionId).Delete(&entity.Transaction{}).Error
}

func (repository *TransactionRepositoryImpl) DeleteTransactionByUserId(ctx context.Context, tx *gorm.DB, userId string) error {
	return tx.WithContext(ctx).Where("user_id", userId).Delete(&entity.Transaction{}).Error
}
