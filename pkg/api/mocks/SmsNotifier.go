// Code generated by mockery 2.7.4. DO NOT EDIT.

package mocks

import (
	model "coding-challenge-go/pkg/model"

	mock "github.com/stretchr/testify/mock"
)

// SmsNotifier is an autogenerated mock type for the SmsNotifier type
type SmsNotifier struct {
	mock.Mock
}

// StockChanged provides a mock function with given fields: oldStock, newStock, _a2, productName
func (_m *SmsNotifier) StockChanged(oldStock int, newStock int, _a2 model.Seller, productName string) string {
	ret := _m.Called(oldStock, newStock, _a2, productName)

	var r0 string
	if rf, ok := ret.Get(0).(func(int, int, model.Seller, string) string); ok {
		r0 = rf(oldStock, newStock, _a2, productName)
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}
