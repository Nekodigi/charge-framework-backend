package subscription

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/checkout/session"
)

func (co *Subscription) HandleSubscribe(e *gin.Engine) {
	// our basic charge API route
	stripe.Key = co.StripeSecret
	e.GET("/subscribe/:service_id/:user_id/:plan_id", func(c *gin.Context) {
		serviceId := c.Param("service_id")
		planId := c.Param("plan_id")
		userId := c.Param("user_id")

		if serviceId == "" || userId == "" || planId == "" {
			c.Status(http.StatusBadRequest)
			return
		}

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
		params.AddMetadata("service_id", serviceId)
		params.AddMetadata("plan_id", planId)
		params.AddMetadata("user_id", userId)
		params.AddExpand("payment_intent") // be careful

		s, _ := session.New(params)
		if s.PaymentStatus == "paid" ||
			s.PaymentStatus == "unpaid" {
			fmt.Println(`register:${name}`)
		}
		//c.Redirect(303, s.URL)
		c.JSON(200, gin.H{
			"url": s.URL,
		})
	})
}
