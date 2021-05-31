package config

import "github.com/gin-gonic/gin"

type GinContext struct {
	*gin.Context
	Config Config
}

func NewGinContext(c *gin.Context) *GinContext {
	return &GinContext{Context: c, Config: defaultConfig}
}

type Config struct {
	NotifySellerViaEmail bool
	NotifySellerViaSms   bool
}

var defaultConfig = Config{
	NotifySellerViaEmail: true,
	NotifySellerViaSms:   false,
}
