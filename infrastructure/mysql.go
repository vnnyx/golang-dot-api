package infrastructure

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/vnnyx/golang-dot-api/exception"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewMySQLDatabase(configuration *Config) *gorm.DB {
	ctx, cancel := NewMySQLContext()
	defer cancel()

	sqlDB, err := sql.Open("mysql", configuration.MysqlHostSlave)
	exception.PanicIfNeeded(err)

	err = sqlDB.PingContext(ctx)
	exception.PanicIfNeeded(err)

	mysqlPoolMax := configuration.MysqlPoolMax

	mysqlIdleMax := configuration.MysqlIdleMax

	mysqlMaxLifeTime := configuration.MysqlMaxLifeTimeMinute

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(mysqlIdleMax)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(mysqlPoolMax)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(time.Duration(mysqlMaxLifeTime) * time.Minute)

	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn: sqlDB,
	}), &gorm.Config{})
	exception.PanicIfNeeded(err)
	return gormDB
}

func NewMySQLContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 10*time.Second)
}
