package seller

import (
	"coding-challenge-go/pkg/model"
	"fmt"
)

// SmsNotifier provides functionality to notify seller via SMS
type SmsNotifier interface {
	StockChanged(oldStock int, newStock int, seller model.Seller, productName string) string
}

type SmsProvider struct{}

func NewSmsProvider() SmsNotifier {
	return &SmsProvider{}
}

// StockChanged provides an SMS template on stock change
func (sp *SmsProvider) StockChanged(oldStock int, newStock int, seller model.Seller, productName string) string {
	return fmt.Sprintf(
		"SMS Warning sent to {%d} (Phone: {%s}): {%s} Product stock changed from %d to %d",
		seller.SellerID,
		seller.Phone,
		productName,
		oldStock,
		newStock)
}

// EmailNotifier provides functionality to notify seller via Email
type EmailNotifier interface {
	StockChanged(oldStock int, newStock int, email string) string
}

type EmailProvider struct{}

func NewEmailProvider() EmailNotifier {
	return &EmailProvider{}
}

// StockChanged provides an Email template on stock change
func (ep *EmailProvider) StockChanged(oldStock int, newStock int, email string) string {
	return fmt.Sprintf(
		"Email Warning sent to (Email: {%s}): Product stock changed from %d to %d",
		email,
		oldStock,
		newStock)
}
