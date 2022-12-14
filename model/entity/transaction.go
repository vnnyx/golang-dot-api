package entity

type Transaction struct {
	TransactionID string `gorm:"column:transaction_id;primaryKey;type:varchar(255)"`
	Name          string `gorm:"column:name;type:varchar(50)"`
	UserID        string `gorm:"column:user_id;type:varchar(255)"`
	User          *User  `gorm:"association_foreignkey:UserID;references:UserID"`
}

func (Transaction) TableName() string {
	return "transactions"
}
