package handler

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/checkout/session"
)

type (
	Subscribe struct {
		stripeSecret string
	}
)

func (co *Subscribe) Handle(e *gin.Engine) {
	// our basic charge API route
	e.POST("/subscribe", func(c *gin.Context) {

		stripe.Key = co.stripeSecret

		params := &stripe.CheckoutSessionParams{
			LineItems: []*stripe.CheckoutSessionLineItemParams{
				{
					Price:    stripe.String("price_1NHh0RErQLZ12HR8rh6waQDL"),
					Quantity: stripe.Int64(1),
				},
			},
			AllowPromotionCodes: stripe.Bool(true),
			Mode:                stripe.String("subscription"),
			SuccessURL:          stripe.String("https://example.com/success"),
			CancelURL:           stripe.String("https://example.com/success"),
		}
		params.AddMetadata("service_id", c.Query("service_id"))
		params.AddMetadata("plan_id", c.Query("plan_id"))
		params.AddMetadata("user_id", c.Query("user_id"))
		params.AddExpand("payment_intent") // be careful

		s, _ := session.New(params)
		if s.PaymentStatus == "paid" ||
			s.PaymentStatus == "unpaid" {
			fmt.Println(`register:${name}`)
		}
		c.Redirect(303, s.URL)
	})
}
