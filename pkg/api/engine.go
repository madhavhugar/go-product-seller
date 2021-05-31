package api

import (
	"coding-challenge-go/pkg/api/product"
	"coding-challenge-go/pkg/api/seller"
	"coding-challenge-go/pkg/config"
	database "coding-challenge-go/pkg/db"
	"github.com/gin-gonic/gin"
)

const (
	APIv1 = "api/v1"
	APIv2 = "api/v2"
)

// WrapConfig wraps configuration to a given handler's context
// Thereby making config.Config accessible to the handler method
func WrapConfig(handler func(*config.GinContext)) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := config.NewGinContext(c)
		handler(ctx)
	}
}

// CreateAPIEngine creates engine instance that serves API endpoints,
// consider it as a router for incoming requests.
func CreateAPIEngine(db database.Adapter) (*gin.Engine, error) {
	r := gin.New()
	productStore := product.NewStore(db)
	sellerStore := seller.NewStore(db)
	emailProvider := seller.NewEmailProvider()
	smsProvider := seller.NewSmsProvider()
	productController := product.NewController(productStore, sellerStore, emailProvider, smsProvider)
	sellerController := seller.NewController(sellerStore)

	v1 := r.Group(APIv1)
	v1.GET("products", productController.List)
	v1.GET("product", productController.Get)
	v1.POST("product", productController.Post)
	v1.PUT("product", WrapConfig(productController.Put))
	v1.DELETE("product", productController.Delete)
	v1.GET("sellers", sellerController.List)

	v2 := r.Group(APIv2)
	v2.GET("product", productController.GetWithSellerLinks)
	v2.GET("products", productController.ListWithSellerLinks)
	v2.POST("product", productController.Post)
	v2.PUT("product", WrapConfig(productController.Put))
	v2.DELETE("product", productController.Delete)
	v2.GET("sellers", sellerController.List)
	v2.GET("sellers/top10", sellerController.ListTopSellersByProductCount)

	return r, nil
}
