package seller

import "database/sql"

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

type Repository struct {
	db *sql.DB
}

func (r *Repository) FindByUUID(uuid string) (*Seller, error) {
	rows, err := r.db.Query("SELECT id_seller, name, email, phone, uuid FROM seller WHERE uuid = ?", uuid)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	if !rows.Next() {
		return nil, nil
	}

	seller := &Seller{}

	err = rows.Scan(&seller.SellerID, &seller.Name, &seller.Email, &seller.Phone, &seller.UUID)

	if err != nil {
		return nil, err
	}

	return seller, nil
}

func (r *Repository) list() ([]*Seller, error) {
	rows, err := r.db.Query("SELECT id_seller, name, email, phone, uuid FROM seller")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var sellers []*Seller

	for rows.Next() {
		seller := &Seller{}

		err := rows.Scan(&seller.SellerID, &seller.Name, &seller.Email, &seller.Phone, &seller.UUID)
		if err != nil {
			return nil, err
		}

		sellers = append(sellers, seller)
	}

	return sellers, nil
}

func (r *Repository) listTopSellersByProductCount() ([]*Seller, error) {
	// TODO: optimize query
	rows, err := r.db.Query("SELECT id_seller, name, email, phone, uuid from seller WHERE id_seller IN (SELECT fk_seller FROM (SELECT COUNT(id_product), fk_seller FROM product GROUP BY fk_seller ORDER BY 1 DESC LIMIT 10) AS AGG)	")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var sellers []*Seller

	for rows.Next() {
		seller := &Seller{}

		err := rows.Scan(&seller.SellerID, &seller.Name, &seller.Email, &seller.Phone, &seller.UUID)
		if err != nil {
			return nil, err
		}

		sellers = append(sellers, seller)
	}

	return sellers, nil
}
