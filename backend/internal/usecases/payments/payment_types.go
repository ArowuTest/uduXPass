package payments

// MoMoPaymentRequest represents a Mobile Money payment request
type MoMoPaymentRequest struct {
	Amount       float64 `json:"amount"`
	Currency     string  `json:"currency"`
	ExternalID   string  `json:"external_id"`
	Phone        string  `json:"phone"`
	PayerMessage string  `json:"payer_message"`
	PayeeNote    string  `json:"payee_note"`
}

// PaystackPaymentRequest represents a Paystack payment request
type PaystackPaymentRequest struct {
	Amount     float64           `json:"amount"`
	Currency   string            `json:"currency"`
	Email      string            `json:"email"`
	Reference  string            `json:"reference"`
	CallbackURL string           `json:"callback_url"`
	Metadata   map[string]string `json:"metadata"`
}

// PaymentCustomerInfo represents customer information for payments
type PaymentCustomerInfo struct {
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

