package store

import (
	"coding-challenge-go/pkg/model"
)

// Product provides all the functionality to access model.Product
type Product interface {
	Insert(product *model.Product) error
	Update(product *model.Product) error
	Delete(product *model.Product) error
	FindByUUID(uuid string) (*model.Product, error)
	List(offset int, limit int) ([]*model.Product, error)
}
