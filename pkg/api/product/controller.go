package product

import (
	"coding-challenge-go/pkg/api/seller"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"net/http"
)

const (
	LIST_PAGE_SIZE = 10
)

func NewController(repository *repository, sellerRepository *seller.Repository, sellerEmailProvider *seller.EmailProvider) *controller {
	return &controller{
		repository:          repository,
		sellerRepository:    sellerRepository,
		sellerEmailProvider: sellerEmailProvider,
	}
}

type controller struct {
	repository          *repository
	sellerRepository    *seller.Repository
	sellerEmailProvider *seller.EmailProvider
	sellerSmsProvider   *seller.SmsProvider
}

func (pc *controller) List(c *gin.Context) {
	request := &struct {
		Page int `form:"page,default=1"`
	}{}

	if err := c.ShouldBindQuery(request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	products, err := pc.repository.list((request.Page-1)*LIST_PAGE_SIZE, LIST_PAGE_SIZE)

	if err != nil {
		log.Error().Err(err).Msg("Fail to query product list")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Fail to query product list"})
		return
	}

	productsJson, err := json.Marshal(products)

	if err != nil {
		log.Error().Err(err).Msg("Fail to marshal products")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Fail to marshal products"})
		return
	}

	c.Data(http.StatusOK, "application/json; charset=utf-8", productsJson)
}

func (pc *controller) ListWithSellerLinks(c *gin.Context) {
	request := &struct {
		Page int `form:"page,default=1"`
	}{}

	if err := c.ShouldBindQuery(request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	products, err := pc.repository.list((request.Page-1)*LIST_PAGE_SIZE, LIST_PAGE_SIZE)

	if err != nil {
		log.Error().Err(err).Msg("Fail to query product list")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Fail to query product list"})
		return
	}

	var productsWithLinks []productWithSellerLinks
	for _, product := range products {
		productsWithLinks = append(productsWithLinks, productWithSellerLinks{
			ProductID: product.ProductID,
			UUID:      product.UUID,
			Name:      product.Name,
			Brand:     product.Brand,
			Stock:     product.Stock,
			Seller: sellerLinks{
				SellerID: product.SellerUUID,
				Links: links{
					// TODO: replace with global variables
					Self: self{Href: fmt.Sprintf("%s/%s/sellers/%s", "hostname", "version", product.SellerUUID)},
				},
			},
		})
	}

	productsJson, err := json.Marshal(productsWithLinks)

	if err != nil {
		log.Error().Err(err).Msg("Fail to marshal products")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Fail to marshal products"})
		return
	}

	c.Data(http.StatusOK, "application/json; charset=utf-8", productsJson)
}

func (pc *controller) Get(c *gin.Context) {
	request := &struct {
		UUID string `form:"id" binding:"required"`
	}{}

	if err := c.ShouldBindQuery(request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product, err := pc.repository.findByUUID(request.UUID)

	if err != nil {
		log.Error().Err(err).Msg("Fail to query product by uuid")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Fail to query product by uuid"})
		return
	}

	productJson, err := json.Marshal(product)

	if err != nil {
		log.Error().Err(err).Msg("Fail to marshal product")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Fail to marshal product"})
		return
	}

	c.Data(http.StatusOK, "application/json; charset=utf-8", productJson)
}

func (pc *controller) GetWithSellerLinks(c *gin.Context) {
	request := &struct {
		UUID string `form:"id" binding:"required"`
	}{}

	if err := c.ShouldBindQuery(request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product, err := pc.repository.findByUUID(request.UUID)

	if err != nil {
		log.Error().Err(err).Msg("Fail to query product by uuid")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Fail to query product by uuid"})
		return
	}

	productMeta := productWithSellerLinks{
		ProductID: product.ProductID,
		UUID:      product.UUID,
		Name:      product.Name,
		Brand:     product.Brand,
		Stock:     product.Stock,
		Seller: sellerLinks{
			SellerID: product.SellerUUID,
			Links: links{
				// TODO: replace with global variables
				Self: self{Href: fmt.Sprintf("%s/%s/sellers/%s", "hostname", "version", product.SellerUUID)},
			},
		},
	}

	productJson, err := json.Marshal(productMeta)

	if err != nil {
		log.Error().Err(err).Msg("Fail to marshal product")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Fail to marshal product"})
		return
	}

	c.Data(http.StatusOK, "application/json; charset=utf-8", productJson)
}

func (pc *controller) Post(c *gin.Context) {
	request := &struct {
		Name   string `form:"name"`
		Brand  string `form:"brand"`
		Stock  int    `form:"stock"`
		Seller string `form:"seller"`
	}{}

	if err := c.ShouldBindJSON(request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	seller, err := pc.sellerRepository.FindByUUID(request.Seller)

	if err != nil {
		log.Error().Err(err).Msg("Fail to query seller by UUID")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Fail to query seller by UUID"})
		return
	}

	if seller == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Seller is not found"})
		return
	}

	product := &product{
		UUID:       uuid.New().String(),
		Name:       request.Name,
		Brand:      request.Brand,
		Stock:      request.Stock,
		SellerUUID: seller.UUID,
	}

	err = pc.repository.insert(product)

	if err != nil {
		log.Error().Err(err).Msg("Fail to insert product")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Fail to insert product"})
		return
	}

	jsonData, err := json.Marshal(product)

	if err != nil {
		log.Error().Err(err).Msg("Fail to marshal product")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Fail to marshal product"})
		return
	}

	c.Data(http.StatusOK, "application/json; charset=utf-8", jsonData)
}

func (pc *controller) Put(c *gin.Context) {
	queryRequest := &struct {
		UUID string `form:"id" binding:"required"`
	}{}

	if err := c.ShouldBindQuery(queryRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product, err := pc.repository.findByUUID(queryRequest.UUID)

	if err != nil {
		log.Error().Err(err).Msg("Fail to query product by uuid")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Fail to query product by uuid"})
		return
	}

	request := &struct {
		Name  string `form:"name"`
		Brand string `form:"brand"`
		Stock int    `form:"stock"`
	}{}

	if err := c.ShouldBindJSON(request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	oldStock := product.Stock

	product.Name = request.Name
	product.Brand = request.Brand
	product.Stock = request.Stock

	err = pc.repository.update(product)

	if err != nil {
		log.Error().Err(err).Msg("Fail to insert product")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Fail to insert product"})
		return
	}

	// TODO: identify an async way to handle notifications
	if oldStock != product.Stock {
		seller, err := pc.sellerRepository.FindByUUID(product.SellerUUID)

		if err != nil {
			log.Error().Err(err).Msg("Fail to query seller by UUID")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Fail to query seller by UUID"})
			return
		}
		// TODO: use configuration here to send Sms and/or Email
		notification := pc.sellerEmailProvider.StockChanged(oldStock, product.Stock, *seller, product.Name)
		fmt.Println(notification)
		notification = pc.sellerSmsProvider.StockChanged(oldStock, product.Stock, *seller, product.Name)
		fmt.Println(notification)
	}

	jsonData, err := json.Marshal(product)

	if err != nil {
		log.Error().Err(err).Msg("Fail to marshal product")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Fail to marshal product"})
		return
	}

	c.Data(http.StatusOK, "application/json; charset=utf-8", jsonData)
}

func (pc *controller) Delete(c *gin.Context) {
	request := &struct {
		UUID string `form:"id" binding:"required"`
	}{}

	if err := c.ShouldBindQuery(request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product, err := pc.repository.findByUUID(request.UUID)

	if err != nil {
		log.Error().Err(err).Msg("Fail to query product by uuid")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Fail to query product by uuid"})
		return
	}

	if product == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Product is not found"})
		return
	}

	err = pc.repository.delete(product)

	if err != nil {
		log.Error().Err(err).Msg("Fail to delete product")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Fail to delete product"})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}
