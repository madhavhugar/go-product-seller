package model

// Product represents the fields of product entity
type Product struct {
	ProductID  int    `json:"-"`
	UUID       string `json:"uuid"`
	Name       string `json:"name"`
	Brand      string `json:"brand"`
	Stock      int    `json:"stock"`
	SellerUUID string `json:"seller_uuid"`
}

// TODO: should we use an anonymous struct
type Self struct {
	Href string `json:"href"`
}

type Links struct {
	Self Self `json:"self"`
}

type SellerLinks struct {
	SellerID string `json:"uuid"`
	Links    Links  `json:"_links"`
}

type ProductWithSellerLinks struct {
	Product
	Seller SellerLinks `json:"seller"`
}
