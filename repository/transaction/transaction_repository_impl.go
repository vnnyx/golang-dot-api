package transaction

import (
	"github.com/vnnyx/golang-dot-api/model/entity"
	"gorm.io/gorm"
)

type TransactionRepositoryImpl struct {
	*gorm.DB
}

func NewTransactionRepository(DB *gorm.DB) TransactionRepository {
	return &TransactionRepositoryImpl{DB: DB}
}

func (repository *TransactionRepositoryImpl) InsertTransaction(transaction entity.Transaction) (entity.Transaction, error) {
	err := repository.DB.Create(&transaction).Error
	return transaction, err
}

func (repository *TransactionRepositoryImpl) FindTransactionByID(transactionId string) (transaction entity.Transaction, err error) {
	err = repository.DB.Where("transaction_id", transactionId).First(&transaction).Error
	return transaction, err
}

func (repository *TransactionRepositoryImpl) FindAllTransaction() (transactions []entity.Transaction, err error) {
	err = repository.DB.Find(&transactions).Error
	return transactions, err
}

func (repository *TransactionRepositoryImpl) FindTransactionByUserId(userId string) (transactions []entity.Transaction, err error) {
	err = repository.DB.Where("user_id", userId).Find(&transactions).Error
	return transactions, err
}

func (repository *TransactionRepositoryImpl) UpdateTransaction(transaction entity.Transaction) (entity.Transaction, error) {
	err := repository.DB.Where("transaction_id", transaction.TransactionID).Updates(&transaction).Error
	return transaction, err
}

func (repository *TransactionRepositoryImpl) DeleteTransaction(transactionId string) error {
	return repository.DB.Where("transaction_id", transactionId).Delete(&entity.Transaction{}).Error
}

func (repository *TransactionRepositoryImpl) DeleteTransactionByUserId(userId string) error {
	return repository.DB.Where("user_id", userId).Delete(&entity.Transaction{}).Error
}
