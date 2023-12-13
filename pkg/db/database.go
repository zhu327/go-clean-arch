package db

import (
	"fmt"
	"go-wire/pkg/config"
	"go-wire/pkg/domain"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDatabase(config config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s dbname=%s port=%s password=%s", config.DBHost, config.DBUser, config.DBName, config.DBPort, config.DBPassword)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true,
	})

	if err != nil {
		return nil, err
	}

	//migrate database
	err = db.AutoMigrate(domain.User{})
	if err != nil {
		log.Fatal("error migrating database: ", err)
	}

	return db, nil
}
