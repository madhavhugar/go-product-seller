package seller

import (
	"coding-challenge-go/pkg/store"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"net/http"
)

// API provides the seller APIs
type API interface {
	List(c *gin.Context)
	ListTopSellersByProductCount(c *gin.Context)
}

func NewController(sellerStore store.Seller) API {
	return &controller{
		sellerStore: sellerStore,
	}
}

type controller struct {
	sellerStore store.Seller
}

func (pc *controller) List(c *gin.Context) {
	sellers, err := pc.sellerStore.List()

	if err != nil {
		log.Error().Err(err).Msg("Fail to query seller list")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Fail to query seller list"})
		return
	}

	sellersJson, err := json.Marshal(sellers)

	if err != nil {
		log.Error().Err(err).Msg("Fail to marshal sellers")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Fail to marshal sellers"})
		return
	}

	c.Data(http.StatusOK, "application/json; charset=utf-8", sellersJson)
}

func (pc *controller) ListTopSellersByProductCount(c *gin.Context) {
	sellers, err := pc.sellerStore.ListTopSellersByProductCount()

	if err != nil {
		log.Error().Err(err).Msg("Fail to query top sellers list")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Fail to query top sellers list"})
		return
	}

	sellersJson, err := json.Marshal(sellers)

	if err != nil {
		log.Error().Err(err).Msg("Fail to marshal sellers")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Fail to marshal sellers"})
		return
	}

	c.Data(http.StatusOK, "application/json; charset=utf-8", sellersJson)
}
