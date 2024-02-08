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

	if err := product.DB.Raw(query).Scan(&res).Error; err != nil {
		return nil, fmt.Errorf("error while getting products : %v", err)
	}

	return res, nil
}

func (product *ProductAdapter) GetProduct(id uint) (entities.Products, error) {
	var res entities.Products

	query := "SELECT * FROM products WHERE id = $1"

	return res, product.DB.Raw(query, id).Scan(&res).Error
}

func (product *ProductAdapter) IncrementStock(id uint, stock int) (entities.Products, error) {
	var res entities.Products

	query := "UPDATE products SET quantity = quantity + $1 WHERE id = $2"

	return res, product.DB.Raw(query, stock, id).Scan(&res).Error
}

func (product *ProductAdapter) DecrementStock(id uint, stock int) (entities.Products, error) {
	var res entities.Products

	query := "UPDATE products SET quantity = quantity - $1 WHERE id = $2"

	tx := product.DB.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Raw(query, stock, id).Scan(&res).Error; err != nil {
		tx.Rollback()
		return res, err
	}

	if res.Quantity < 0 {
		tx.Rollback()
		return res, fmt.Errorf("quantity cannot be negative")
	}

	if err := tx.Commit().Error; err != nil {
		return res, err
	}

	return res, nil
}
