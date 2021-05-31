package seller

import (
	"coding-challenge-go/pkg/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEmailProvider_StockChanged(t *testing.T) {
	t.Run("should return correct email template", func(t *testing.T) {
		ep := &EmailProvider{}

		got := ep.StockChanged(2, 1, "test@mail.com")
		expected := "Email Warning sent to (Email: {test@mail.com}): Product stock changed from 2 to 1"

		assert.Equal(t, expected, got)
	})
}

func TestSmsProvider_StockChanged(t *testing.T) {
	t.Run("should return correct sms template", func(t *testing.T) {
		sp := &SmsProvider{}

		got := sp.StockChanged(2, 1, model.Seller{
			SellerID: 123,
			UUID:     "123",
			Name:     "seller-name",
			Email:    "seller-email",
			Phone:    "seller-phone",
		}, "product-name")
		expected := "SMS Warning sent to {123} (Phone: {seller-phone}): {product-name} Product stock changed from 2 to 1"

		assert.Equal(t, expected, got)
	})
}
