package models

import (
	"fmt"

	config "github.com/johnjiangtw0804/chatbot-back-end-authentication/env"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DBWrapper struct {
	dbConnection *gorm.DB
}

func RegisterDB(env *config.Configuration) (*DBWrapper, error) {
	var err error
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=%s",
		env.DBHost,
		env.DBUser,
		env.DBPassword,
		env.DBName,
		env.DBPort,
		viper.GetString("SERVER_TIMEZONE"),
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	return &DBWrapper{dbConnection: db}, err
}
