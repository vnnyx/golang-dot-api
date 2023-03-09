package infrastructure

import (
	"github.com/spf13/viper"
	"github.com/vnnyx/golang-dot-api/exception"
)

type Config struct {
	AppPort                string `mapstructure:"APP_PORT"`
	MysqlHostSlave         string `mapstructure:"MYSQL_HOST_SLAVE"`
	MysqlPoolMin           int    `mapstructure:"MYSQL_POOL_MIN"`
	MysqlPoolMax           int    `mapstructure:"MYSQL_POOL_MAX"`
	MysqlIdleMax           int    `mapstructure:"MYSQL_IDLE_MAX"`
	MysqlMaxIdleTimeMinute int    `mapstructure:"MYSQL_MAX_IDLE_TIME_MINUTE"`
	MysqlMaxLifeTimeMinute int    `mapstructure:"MYSQL_MAX_LIFE_TIME_MINUTE"`
	JWTPublicKey           string `mapstructure:"JWT_PUBLIC_KEY"`
	JWTSecretKey           string `mapstructure:"JWT_SECRET_KEY"`
	JWTMinute              int    `mapstructure:"JWT_MINUTE"`
	RedisHost              string `mapstructure:"REDIS_HOST"`
	RedisPassword          string `mapstructure:"REDIS_PASSWORD"`
	MailHost               string `mapstructure:"MAIL_HOST"`
	MailPort               int    `mapstructure:"MAIL_PORT"`
	MailUsername           string `mapstructure:"MAIL_USERNAME"`
	MailPassword           string `mapstructure:"MAIL_PASSWORD"`
	BrokerHost             string `mapstructure:"BROKER_HOST"`
}

func NewConfig(configName string) *Config {
	config := &Config{}
	viper.AddConfigPath(".")
	viper.SetConfigName(configName)
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	exception.PanicIfNeeded(err)
	err = viper.Unmarshal(&config)
	exception.PanicIfNeeded(err)
	return config
}
