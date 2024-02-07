package db

import (
	"github.com/Nishad4140/product_service/entities"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB(connect string) (*gorm.DB, error) {

	db, err := gorm.Open(postgres.Open(connect), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&entities.Products{})

	return db, nil
}
