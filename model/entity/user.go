package entity

import "time"

type User struct {
	UserID    string    `gorm:"column:id;primaryKey;type:varchar(255)"`
	Username  string    `gorm:"column:username;unique;type:varchar(50)"`
	Email     string    `gorm:"column:email;type:varchar(100);unique"`
	Handphone string    `gorm:"column:handphone;type:varchar(20)"`
	Password  string    `gorm:"column:password;type:varchar(255)"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"-"`
}

func (User) TableName() string {
	return "users"
}
