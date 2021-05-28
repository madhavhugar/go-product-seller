package product

type product struct {
	ProductID  int    `json:"-"`
	UUID       string `json:"uuid"`
	Name       string `json:"name"`
	Brand      string `json:"brand"`
	Stock      int    `json:"stock"`
	SellerUUID string `json:"seller_uuid"`
}

// TODO: should we use an anonymous struct
type self struct {
	Href string `json:"href"`
}

type links struct {
	Self self `json:"self"`
}

type sellerLinks struct {
	SellerID string `json:"uuid"`
	Links    links  `json:"_links"`
}

type productWithSellerLinks struct {
	ProductID int         `json:"-"`
	UUID      string      `json:"uuid"`
	Name      string      `json:"name"`
	Brand     string      `json:"brand"`
	Stock     int         `json:"stock"`
	Seller    sellerLinks `json:"seller"`
}
