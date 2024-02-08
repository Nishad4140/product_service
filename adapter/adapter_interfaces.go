package adapter

import (
	"github.com/Nishad4140/product_service/entities"
)

type AdapterInterface interface {
	AddProduct(req entities.Products) (entities.Products, error)
	GetProduct(id uint) (entities.Products, error)
	GetAllProducts() ([]entities.Products, error)
	IncrementStock(id uint, stock int) (entities.Products, error)
	DecrementStock(id uint, stock int) (entities.Products, error)
}
