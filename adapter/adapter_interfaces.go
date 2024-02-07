package adapter

import "github.com/Nishad4140/product_service/entities"

type AdapterInterface interface {
	AddProduct(req entities.Products) (entities.Products, error)
}
