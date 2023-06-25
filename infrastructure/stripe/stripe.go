package stripe

import (
	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/customer"
)

var (
	St *Stripe
)

type (
	Stripe struct {
	}
)

func NewStripe(StripeSecret string) *Stripe {
	stripe.Key = StripeSecret
	St = &Stripe{}
	return St
}

func (st *Stripe) GetCustomerById(customerId string) *stripe.Customer {
	var cus *stripe.Customer
	if customerId == "" {
		params := &stripe.CustomerParams{}
		cus, _ = customer.New(params)
	} else {
		cus, _ = customer.Get(customerId, nil)
	}
	return cus
}

func (st *Stripe) NewCustomer() *stripe.Customer {
	params := &stripe.CustomerParams{}
	cus, _ := customer.New(params)
	return cus
}

func NewCustomer() *stripe.Customer {
	params := &stripe.CustomerParams{}
	cus, _ := customer.New(params)
	return cus
}
