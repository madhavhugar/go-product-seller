package api

import (
	"coding-challenge-go/pkg/api/product"
	"coding-challenge-go/pkg/api/seller"
	"database/sql"
	"github.com/gin-gonic/gin"
)

// CreateAPIEngine creates engine instance that serves API endpoints,
// consider it as a router for incoming requests.
func CreateAPIEngine(db *sql.DB) (*gin.Engine, error) {
	r := gin.New()
	productRepository := product.NewRepository(db)
	sellerRepository := seller.NewRepository(db)
	emailProvider := seller.NewEmailProvider()
	productController := product.NewController(productRepository, sellerRepository, emailProvider)
	sellerController := seller.NewController(sellerRepository)

	// TODO: define application level configuration to send apt notification
	v1 := r.Group("api/v1")
	v1.GET("products", productController.List)
	v1.GET("product", productController.Get)
	v1.POST("product", productController.Post)
	v1.PUT("product", productController.Put)
	v1.DELETE("product", productController.Delete)
	v1.GET("sellers", sellerController.List)

	v2 := r.Group("api/v2")
	v2.GET("product", productController.GetWithSellerLinks)
	v2.GET("products", productController.ListWithSellerLinks)
	v2.POST("product", productController.Post)
	v2.PUT("product", productController.Put)
	v2.DELETE("product", productController.Delete)
	v2.GET("sellers", sellerController.List)
	v2.GET("sellers/top10", sellerController.ListTopSellersByProductCount)

	return r, nil
}
