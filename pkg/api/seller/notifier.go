package seller

import "fmt"

type Notifier interface {
	StockChanged(oldStock int, newStock int, seller Seller, productName string) string
}

type SmsProvider struct{}

func (sp *SmsProvider) StockChanged(oldStock int, newStock int, seller Seller, productName string) string {
	return fmt.Sprintf(
		"SMS Warning set to {%s} (Phone: {%s}): {%s} Product stock changed from %d to %d",
		seller.SellerID,
		seller.Phone,
		productName,
		oldStock,
		newStock)
}

func NewEmailProvider() *EmailProvider {
	return &EmailProvider{}
}

type EmailProvider struct{}

func (ep *EmailProvider) StockChanged(oldStock int, newStock int, seller Seller, productName string) string {
	return fmt.Sprintf(
		"SMS Warning set to {%s} (Phone: {%s}): {%s} Product stock changed from %d to %d",
		seller.SellerID,
		seller.Phone,
		productName,
		oldStock,
		newStock)
}
