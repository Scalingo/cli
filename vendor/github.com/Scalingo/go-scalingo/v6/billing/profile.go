package billing

type PaymentMethodType string

const (
	Stripe PaymentMethodType = "stripe"
	Paypal                   = "paypal"
)

type Profile struct {
	Company           string              `json:"company"`
	VATNumber         string              `json:"vat_number"`
	PaymentMethodType PaymentMethodType   `json:"payment_method_type"`
	Stripe            StripePaymentMethod `json:"stripe"`
}

type StripePaymentMethod struct {
	Brand string `json:"brand"`
	Last4 string `json:"last4"`
	Exp   string `json:"exp"`
}
