package product_test

import (
	"bytes"
	"coding-challenge-go/pkg/api/mocks"
	"coding-challenge-go/pkg/api/product"
	"coding-challenge-go/pkg/config"
	"coding-challenge-go/pkg/model"
	storemock "coding-challenge-go/pkg/store/mocks"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestController_GetWithSellerLinks(t *testing.T) {
	t.Run("should return product with seller link", func(t *testing.T) {
		productUUID := "product-uuid"
		productStoreMock := &storemock.Product{}
		productStoreMock.On("FindByUUID", productUUID).Return(&model.Product{
			ProductID:  234,
			UUID:       productUUID,
			Name:       "product-name",
			Brand:      "product-brand",
			Stock:      10,
			SellerUUID: "seller-uuid",
		}, nil)

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		apiURI := fmt.Sprintf("http://localhost/api/v2/product?id=%s", productUUID)
		ctx.Request, _ = http.NewRequest("GET", apiURI, nil)

		productController := product.NewController(productStoreMock, nil, nil, nil)
		productController.GetWithSellerLinks(ctx)

		res := w.Result()
		b, _ := ioutil.ReadAll(res.Body)
		assert.Equal(t, res.StatusCode, 200)
		assert.Equal(t, string(b), "{\"uuid\":\"product-uuid\",\"name\":\"product-name\",\"brand\":\"product-brand\",\"stock\":10,\"seller_uuid\":\"\",\"seller\":{\"uuid\":\"seller-uuid\",\"_links\":{\"self\":{\"href\":\"http://localhost:8080/api/v2/sellers/seller-uuid\"}}}}")
	})

	t.Run("should return internal server error on DB query error", func(t *testing.T) {
		productUUID := "product-uuid"
		productStoreMock := &storemock.Product{}
		productStoreMock.On("FindByUUID", productUUID).Return(nil, fmt.Errorf("unable to query"))

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		apiURI := fmt.Sprintf("hostname/api/v2/product?id=%s", productUUID)
		ctx.Request, _ = http.NewRequest("GET", apiURI, nil)

		productController := product.NewController(productStoreMock, nil, nil, nil)
		productController.GetWithSellerLinks(ctx)

		res := w.Result()
		assert.Equal(t, res.StatusCode, 500)
	})
}

func TestController_ListWithSellerLinks(t *testing.T) {
	t.Run("should return products with seller link", func(t *testing.T) {
		productUUID := "product-uuid"
		productStoreMock := &storemock.Product{}
		productStoreMock.On("List", mock.Anything, mock.Anything).Return([]*model.Product{
			{
				ProductID:  234,
				UUID:       productUUID,
				Name:       "product-name",
				Brand:      "product-brand",
				Stock:      10,
				SellerUUID: "seller-uuid",
			},
		}, nil)

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		apiURI := fmt.Sprintf("http://localhost/api/v2/products")
		ctx.Request, _ = http.NewRequest("GET", apiURI, nil)

		productController := product.NewController(productStoreMock, nil, nil, nil)
		productController.ListWithSellerLinks(ctx)

		res := w.Result()
		b, _ := ioutil.ReadAll(res.Body)
		assert.Equal(t, res.StatusCode, 200)
		assert.Equal(t, string(b), "[{\"uuid\":\"product-uuid\",\"name\":\"product-name\",\"brand\":\"product-brand\",\"stock\":10,\"seller_uuid\":\"\",\"seller\":{\"uuid\":\"seller-uuid\",\"_links\":{\"self\":{\"href\":\"http://localhost:8080/api/v2/sellers/seller-uuid\"}}}}]")
	})

	t.Run("should return internal server error on DB query error", func(t *testing.T) {
		productStoreMock := &storemock.Product{}
		productStoreMock.On("List", mock.Anything, mock.Anything).Return(nil, fmt.Errorf("unable to query"))

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		apiURI := fmt.Sprintf("http://localhost/api/v2/products")
		ctx.Request, _ = http.NewRequest("GET", apiURI, nil)

		productController := product.NewController(productStoreMock, nil, nil, nil)
		productController.ListWithSellerLinks(ctx)

		res := w.Result()
		assert.Equal(t, res.StatusCode, 500)
	})
}

func TestController_Put(t *testing.T) {
	t.Run("should notify products with seller link", func(t *testing.T) {
		productUUID := "product-uuid"
		sellerUUID := "seller-uuid"
		oldProduct := model.Product{
			ProductID:  234,
			UUID:       productUUID,
			Name:       "product-name",
			Brand:      "product-brand",
			Stock:      10,
			SellerUUID: sellerUUID,
		}

		productStoreMock := &storemock.Product{}
		productStoreMock.On("FindByUUID", productUUID).Return(&oldProduct, nil)
		productStoreMock.On("Update", mock.Anything).Return(nil)

		sellerStoreMock := &storemock.Seller{}
		sellerStoreMock.On("FindByUUID", sellerUUID).Return(&model.Seller{
			SellerID: 123,
			UUID:     "seller-uuid",
			Name:     "seller-name",
			Email:    "seller-email",
			Phone:    "seller-phone",
		}, nil)

		emailProviderMock := &mocks.EmailNotifier{}
		emailProviderMock.On("StockChanged", mock.Anything, mock.Anything, mock.Anything).Return("Email Template")

		smsProviderMock := &mocks.SmsNotifier{}
		smsProviderMock.On("StockChanged", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return("Sms Template")

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		apiURI := fmt.Sprintf("http://localhost/api/v2/product?id=%s", productUUID)
		body := []byte(`{"name": "product-name", "brand": "product-brand", "stock": 9}`)
		ctx.Request, _ = http.NewRequest("PUT", apiURI, bytes.NewBuffer(body))
		ginCtx := &config.GinContext{
			Context: ctx,
			Config: config.Config{
				NotifySellerViaEmail: true,
				NotifySellerViaSms:   false,
			},
		}

		productController := product.NewController(productStoreMock, sellerStoreMock, emailProviderMock, smsProviderMock)
		productController.Put(ginCtx)

		res := w.Result()
		b, _ := ioutil.ReadAll(res.Body)

		assert.Equal(t, res.StatusCode, 200)
		assert.Equal(t, string(b), "{\"uuid\":\"product-uuid\",\"name\":\"product-name\",\"brand\":\"product-brand\",\"stock\":9,\"seller_uuid\":\"seller-uuid\"}")
		emailProviderMock.AssertCalled(t, "StockChanged", 10, 9, "seller-email")
		smsProviderMock.AssertNotCalled(t, "StockChanged")
	})
}
