package stripe

import (
	"fmt"
	"github.com/stripe/stripe-go/v79"
	"github.com/stripe/stripe-go/v79/customersession"
	"github.com/stripe/stripe-go/v79/paymentintent"
	"weihu_server/library/config"
)

type StripeClient struct {
	secretKey string
}

var NewStripeClient *StripeClient

func init() {
	NewStripeClient = &StripeClient{
		secretKey: config.GetString("stripe.secretKey"),
	}
}

// CreateCustomer 创建客户
func (s *StripeClient) CreateCustomer(customerId string) (*stripe.CustomerSession, error) {
	stripe.Key = s.secretKey

	params := &stripe.CustomerSessionParams{
		Customer: stripe.String(customerId),
		Components: &stripe.CustomerSessionComponentsParams{
			PricingTable: &stripe.CustomerSessionComponentsPricingTableParams{
				Enabled: stripe.Bool(true),
			},
		},
	}
	result, err := customersession.New(params)

	if err != nil {
		return nil, fmt.Errorf("failed to create customer session: %w", err)
	}

	return result, nil
}

// CreatePaymentIntent 创建支付意图
func (s *StripeClient) CreatePaymentIntent(amount int64, currency string) (*stripe.PaymentIntent, error) {
	stripe.Key = s.secretKey

	params := &stripe.PaymentIntentParams{
		Amount:      stripe.Int64(amount),
		Currency:    stripe.String(currency),
		Description: stripe.String("Payment for service"),
		Confirm:     stripe.Bool(true), // 自动确认支付意图
		AutomaticPaymentMethods: &stripe.PaymentIntentAutomaticPaymentMethodsParams{
			Enabled: stripe.Bool(true),
		},
	}

	paymentIntent, err := paymentintent.New(params)
	if err != nil {
		return nil, fmt.Errorf("failed to create payment intent: %v", err)
	}

	//paymentIntent.ClientSecret 客户密钥,用于后续的支付操作
	//paymentIntent.ID 支付意图ID
	return paymentIntent, nil
}

// UpdatePaymentIntentMetadata 更新支付意图元数据
func (s *StripeClient) UpdatePaymentIntentMetadata(clientSecret string, metadata map[string]string) (*stripe.PaymentIntent, error) {
	stripe.Key = s.secretKey

	params := &stripe.PaymentIntentParams{
		Metadata: metadata,
	}

	paymentIntent, err := paymentintent.Update(clientSecret, params)
	if err != nil {
		return nil, fmt.Errorf("failed to update payment intent: %v", err)
	}

	return paymentIntent, nil
}

// ConfirmPaymentIntent 确认支付意图
func (s *StripeClient) ConfirmPaymentIntent(clientSecret string, paymentMethodID string) (*stripe.PaymentIntent, error) {
	stripe.Key = s.secretKey

	params := &stripe.PaymentIntentConfirmParams{
		PaymentMethod: stripe.String(paymentMethodID),
	}

	paymentIntent, err := paymentintent.Confirm(clientSecret, params)
	if err != nil {
		return nil, fmt.Errorf("failed to confirm payment intent: %v", err)
	}

	return paymentIntent, nil
}

// CancelPaymentIntent 取消支付意图
func (s *StripeClient) CancelPaymentIntent(clientSecret string) (*stripe.PaymentIntent, error) {
	stripe.Key = s.secretKey

	params := &stripe.PaymentIntentCancelParams{}
	paymentIntent, err := paymentintent.Cancel(clientSecret, params)
	if err != nil {
		return nil, fmt.Errorf("failed to cancel payment intent: %v", err)
	}

	return paymentIntent, nil
}

// CapturePaymentIntent 捕获支付意图
func (s *StripeClient) CapturePaymentIntent(clientSecret string) (*stripe.PaymentIntent, error) {
	stripe.Key = s.secretKey

	paymentIntent, err := paymentintent.Capture(clientSecret, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to capture payment intent: %v", err)
	}

	return paymentIntent, nil
}

// SearchPaymentIntents 搜索支付意图
func (s *StripeClient) SearchPaymentIntents(customerID string) (*paymentintent.SearchIter, error) {
	stripe.Key = s.secretKey

	params := &stripe.PaymentIntentSearchParams{
		SearchParams: stripe.SearchParams{Query: "amount>1000"},
	}
	result := paymentintent.Search(params)

	return result, nil
}

// ListPaymentIntentsForCustomer 列出与特定客户关联的支付意图
func (s *StripeClient) ListPaymentIntentsForCustomer(customerID string) (*paymentintent.Iter, error) {
	stripe.Key = s.secretKey

	params := &stripe.PaymentIntentListParams{
		Customer: stripe.String(customerID),
	}

	params.Limit = stripe.Int64(3)
	result := paymentintent.List(params)

	return result, nil
}
