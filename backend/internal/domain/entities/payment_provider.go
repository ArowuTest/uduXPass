package entities

// PaymentProvider represents supported payment providers
type PaymentProvider string

const (
	PaymentProviderMoMo     PaymentProvider = "momo"
	PaymentProviderPaystack PaymentProvider = "paystack"
	PaymentProviderCash     PaymentProvider = "cash"
	PaymentProviderBank     PaymentProvider = "bank_transfer"
)

// String returns the string representation of the payment provider
func (pp PaymentProvider) String() string {
	return string(pp)
}

// IsValid checks if the payment provider is valid
func (pp PaymentProvider) IsValid() bool {
	switch pp {
	case PaymentProviderMoMo, PaymentProviderPaystack, PaymentProviderCash, PaymentProviderBank:
		return true
	default:
		return false
	}
}

// GetDisplayName returns a human-readable display name for the payment provider
func (pp PaymentProvider) GetDisplayName() string {
	switch pp {
	case PaymentProviderMoMo:
		return "Mobile Money (MoMo PSB)"
	case PaymentProviderPaystack:
		return "Paystack"
	case PaymentProviderCash:
		return "Cash Payment"
	case PaymentProviderBank:
		return "Bank Transfer"
	default:
		return string(pp)
	}
}

// RequiresOnlineProcessing returns true if the provider requires online processing
func (pp PaymentProvider) RequiresOnlineProcessing() bool {
	switch pp {
	case PaymentProviderMoMo, PaymentProviderPaystack:
		return true
	case PaymentProviderCash, PaymentProviderBank:
		return false
	default:
		return true
	}
}

// SupportsRefunds returns true if the provider supports refunds
func (pp PaymentProvider) SupportsRefunds() bool {
	switch pp {
	case PaymentProviderMoMo, PaymentProviderPaystack:
		return true
	case PaymentProviderCash, PaymentProviderBank:
		return false
	default:
		return false
	}
}

