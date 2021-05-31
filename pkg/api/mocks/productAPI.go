// Code generated by mockery 2.7.4. DO NOT EDIT.

package mocks

import (
	gin "github.com/gin-gonic/gin"
	mock "github.com/stretchr/testify/mock"
)

// API is an autogenerated mock type for the API type
type API struct {
	mock.Mock
}

// Delete provides a mock function with given fields: c
func (_m *API) Delete(c *gin.Context) {
	_m.Called(c)
}

// Get provides a mock function with given fields: c
func (_m *API) Get(c *gin.Context) {
	_m.Called(c)
}

// GetWithSellerLinks provides a mock function with given fields: c
func (_m *API) GetWithSellerLinks(c *gin.Context) {
	_m.Called(c)
}

// List provides a mock function with given fields: c
func (_m *API) List(c *gin.Context) {
	_m.Called(c)
}

// ListWithSellerLinks provides a mock function with given fields: c
func (_m *API) ListWithSellerLinks(c *gin.Context) {
	_m.Called(c)
}

// Post provides a mock function with given fields: c
func (_m *API) Post(c *gin.Context) {
	_m.Called(c)
}

// Put provides a mock function with given fields: c
func (_m *API) Put(c *gin.Context) {
	_m.Called(c)
}