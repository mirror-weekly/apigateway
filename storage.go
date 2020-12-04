package apigateway

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const dsn = ""

func NewDB() (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	return db, err
}
