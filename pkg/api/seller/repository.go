package seller

import (
	database "coding-challenge-go/pkg/db"
	"coding-challenge-go/pkg/model"
	"coding-challenge-go/pkg/store"
)

func NewStore(db database.Adapter) store.Seller {
	return &Repository{db: db}
}

type Repository struct {
	db database.Adapter
}

func (r *Repository) FindByUUID(uuid string) (*model.Seller, error) {
	rows, err := r.db.Query("SELECT id_seller, name, email, phone, uuid FROM seller WHERE uuid = ?", uuid)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	if !rows.Next() {
		return nil, nil
	}

	seller := &model.Seller{}

	err = rows.Scan(&seller.SellerID, &seller.Name, &seller.Email, &seller.Phone, &seller.UUID)

	if err != nil {
		return nil, err
	}

	return seller, nil
}

func (r *Repository) List() ([]*model.Seller, error) {
	rows, err := r.db.Query("SELECT id_seller, name, email, phone, uuid FROM seller")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var sellers []*model.Seller

	for rows.Next() {
		seller := &model.Seller{}

		err := rows.Scan(&seller.SellerID, &seller.Name, &seller.Email, &seller.Phone, &seller.UUID)
		if err != nil {
			return nil, err
		}

		sellers = append(sellers, seller)
	}

	return sellers, nil
}

func (r *Repository) ListTopSellersByProductCount() ([]*model.Seller, error) {
	rows, err := r.db.Query("SELECT id_seller, name, email, phone, uuid from seller WHERE id_seller IN (SELECT fk_seller FROM (SELECT COUNT(id_product), fk_seller FROM product GROUP BY fk_seller ORDER BY 1 DESC LIMIT 10) AS AGG)")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var sellers []*model.Seller

	for rows.Next() {
		seller := &model.Seller{}

		err := rows.Scan(&seller.SellerID, &seller.Name, &seller.Email, &seller.Phone, &seller.UUID)
		if err != nil {
			return nil, err
		}

		sellers = append(sellers, seller)
	}

	return sellers, nil
}
