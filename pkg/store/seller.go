package store

import (
	"coding-challenge-go/pkg/model"
)

// Seller provides all the functionality to access model.Seller
type Seller interface {
	List() ([]*model.Seller, error)
	FindByUUID(uuid string) (*model.Seller, error)
	ListTopSellersByProductCount() ([]*model.Seller, error)
}
