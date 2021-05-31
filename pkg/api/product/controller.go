package product

import (
	"coding-challenge-go/pkg/api/seller"
	"coding-challenge-go/pkg/config"
	"coding-challenge-go/pkg/model"
	"coding-challenge-go/pkg/store"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"net/http"
)

const (
	LIST_PAGE_SIZE = 10
	// TODO: move to init func
	DomainName = "http://localhost:8080"
	APIv1      = "api/v1"
	APIv2      = "api/v2"
)

// API provides the product APIs
type API interface {
	Get(c *gin.Context)
	GetWithSellerLinks(c *gin.Context)
	List(c *gin.Context)
	ListWithSellerLinks(c *gin.Context)
	Post(c *gin.Context)
	Put(c *config.GinContext)
	Delete(c *gin.Context)
}

func NewController(
	productStore store.Product,
	sellerStore store.Seller,
	sellerEmailProvider seller.EmailNotifier,
	sellerSmsProvider seller.SmsNotifier) API {
	return &controller{
		productStore:        productStore,
		sellerStore:         sellerStore,
		sellerEmailProvider: sellerEmailProvider,
		sellerSmsProvider:   sellerSmsProvider,
	}
}

type controller struct {
	productStore        store.Product
	sellerStore         store.Seller
	sellerEmailProvider seller.EmailNotifier
	sellerSmsProvider   seller.SmsNotifier
}

func sellerLink(sellerUUID string) string {
	return fmt.Sprintf("%s/%s/sellers/%s", DomainName, APIv2, sellerUUID)
}

func (pc *controller) list(c *gin.Context) ([]*model.Product, error) {
	request := &struct {
		Page int `form:"page,default=1"`
	}{}

	if err := c.ShouldBindQuery(request); err != nil {
		log.Error().Err(err).Msg("Fail to extract id from request")
		return nil, err
	}

	products, err := pc.productStore.List((request.Page-1)*LIST_PAGE_SIZE, LIST_PAGE_SIZE)

	if err != nil {
		log.Error().Err(err).Msg("Fail to query product list")
		return nil, err
	}

	return products, nil
}

func (pc *controller) List(c *gin.Context) {
	products, err := pc.list(c)
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
	products, err := pc.list(c)
	if err != nil {
		log.Error().Err(err).Msg("Fail to query product list")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Fail to query product list"})
		return
	}

	var productsWithLinks []model.ProductWithSellerLinks
	for _, p := range products {
		productsWithLinks = append(productsWithLinks, model.ProductWithSellerLinks{
			Product: model.Product{
				ProductID: p.ProductID,
				UUID:      p.UUID,
				Name:      p.Name,
				Brand:     p.Brand,
				Stock:     p.Stock,
			},
			Seller: model.SellerLinks{
				SellerID: p.SellerUUID,
				Links: model.Links{
					Self: model.Self{Href: sellerLink(p.SellerUUID)},
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

func (pc *controller) get(c *gin.Context) (*model.Product, error) {
	request := &struct {
		UUID string `form:"id" binding:"required"`
	}{}

	if err := c.ShouldBindQuery(request); err != nil {
		return nil, err
	}

	p, err := pc.productStore.FindByUUID(request.UUID)
	if err != nil {
		log.Error().Err(err).Msg("Fail to query product")
		return nil, err
	}

	return p, nil
}

func (pc *controller) Get(c *gin.Context) {
	product, err := pc.get(c)
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
	p, err := pc.get(c)
	if err != nil {
		log.Error().Err(err).Msg("Fail to query product by uuid")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Fail to query product by uuid"})
		return
	}

	productMeta := model.ProductWithSellerLinks{
		Product: model.Product{
			ProductID: p.ProductID,
			UUID:      p.UUID,
			Name:      p.Name,
			Brand:     p.Brand,
			Stock:     p.Stock,
		},
		Seller: model.SellerLinks{
			SellerID: p.SellerUUID,
			Links: model.Links{
				Self: model.Self{Href: sellerLink(p.SellerUUID)},
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

	seller, err := pc.sellerStore.FindByUUID(request.Seller)

	if err != nil {
		log.Error().Err(err).Msg("Fail to query seller by UUID")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Fail to query seller by UUID"})
		return
	}

	if seller == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Seller is not found"})
		return
	}

	product := &model.Product{
		UUID:       uuid.New().String(),
		Name:       request.Name,
		Brand:      request.Brand,
		Stock:      request.Stock,
		SellerUUID: seller.UUID,
	}

	err = pc.productStore.Insert(product)

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

func (pc *controller) notifySeller(product model.Product, config config.Config, oldStock int) error {
	if !config.NotifySellerViaEmail && !config.NotifySellerViaSms {
		return nil
	}

	s, err := pc.sellerStore.FindByUUID(product.SellerUUID)
	if err != nil {
		log.Error().Err(err).Msg("Fail to query seller by UUID")
		return err
	}

	if config.NotifySellerViaSms {
		notification := pc.sellerSmsProvider.StockChanged(oldStock, product.Stock, *s, product.Name)
		log.Info().Msg(notification)
	}
	if config.NotifySellerViaEmail {
		notification := pc.sellerEmailProvider.StockChanged(oldStock, product.Stock, s.Email)
		log.Info().Msg(notification)
	}

	return nil
}

func (pc *controller) Put(c *config.GinContext) {
	queryRequest := &struct {
		UUID string `form:"id" binding:"required"`
	}{}

	if err := c.ShouldBindQuery(queryRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product, err := pc.productStore.FindByUUID(queryRequest.UUID)

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

	err = pc.productStore.Update(product)

	if err != nil {
		log.Error().Err(err).Msg("Fail to update product")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Fail to update product"})
		return
	}

	// TODO: identify an async way to handle notifications
	if oldStock != product.Stock {
		err = pc.notifySeller(*product, c.Config, oldStock)
		if err != nil {
			log.Error().Err(err).Msg("Failed to notify seller on stock update")
		}
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

	product, err := pc.productStore.FindByUUID(request.UUID)

	if err != nil {
		log.Error().Err(err).Msg("Fail to query product by uuid")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Fail to query product by uuid"})
		return
	}

	if product == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Product is not found"})
		return
	}

	err = pc.productStore.Delete(product)

	if err != nil {
		log.Error().Err(err).Msg("Fail to delete product")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Fail to delete product"})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}
