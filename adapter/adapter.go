package adapter

import (
	"fmt"

	"github.com/Nishad4140/product_service/entities"
	"gorm.io/gorm"
)

type ProductAdapter struct {
	DB *gorm.DB
}

func NewProductAdapter(db *gorm.DB) *ProductAdapter {
	return &ProductAdapter{
		DB: db,
	}
}

func (product *ProductAdapter) AddProduct(req entities.Products) (entities.Products, error) {
	var res entities.Products

	query := "INSERT INTO products (name,price,quantity) VALUES ($1, $2, $3) RETURNING id, name, price, quantity"

	return res, product.DB.Raw(query, req.Name, req.Price, req.Quantity).Scan(&res).Error
}

func (product *ProductAdapter) GetAllProducts() ([]entities.Products, error) {
	var res []entities.Products

	query := "SELECT * FROM products"

	if err := product.DB.Raw(query).Scan(res).Error; err != nil {
		return nil, fmt.Errorf("error while getting products : %v", err)
	}

	return res, nil
}
