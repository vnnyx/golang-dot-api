package migration

import (
	"github.com/vnnyx/golang-dot-api/exception"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB, tables ...interface{}) {
	for _, table := range tables {
		if !db.Migrator().HasTable(table) {
			err := db.Debug().AutoMigrate(table)
			exception.PanicIfNeeded(err)
		}
	}
}
