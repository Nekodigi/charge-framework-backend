package subscription

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/checkout/session"
)

func (s *Subscription) HandleSubscribe(e *gin.Engine) {
	// our basic charge API route
	stripe.Key = s.StripeSecret
	e.GET("/subscribe/:service_id/:user_id/:plan_id", func(c *gin.Context) {
		serviceId := c.Param("service_id")
		planId := c.Param("plan_id")
		userId := c.Param("user_id")
		successUrl := c.Query("success_url")
		cancelUrl := c.Query("cancel_url")

		if serviceId == "" || userId == "" || planId == "" {
			c.Status(http.StatusBadRequest)
			return
		}

		user := s.Fs.GetUserById(serviceId, userId)
		service := s.Fs.GetServiceById(serviceId)

		mode := "subscription"
		priceId := service.Plan[planId].PriceId
		if user.Subscription != "" {
			c.Status(http.StatusBadRequest)
			return
		}
		quota := service.Plan[planId].Quota
		params := &stripe.CheckoutSessionParams{
			LineItems: []*stripe.CheckoutSessionLineItemParams{
				{
					Price:    stripe.String(priceId),
					Quantity: stripe.Int64(int64(quota)),
				},
			},
			AllowPromotionCodes: stripe.Bool(true),
			Mode:                stripe.String(mode),
			SuccessURL:          stripe.String(successUrl),
			CancelURL:           stripe.String(cancelUrl),
		}
		params.AddMetadata("service_id", serviceId)
		params.AddMetadata("plan_id", planId)
		params.AddMetadata("user_id", userId)
		params.AddMetadata("mode", mode)
		params.AddExpand("payment_intent") // be careful

		s, _ := session.New(params)
		if s.PaymentStatus == "paid" ||
			s.PaymentStatus == "unpaid" {
			fmt.Println(`register:${name}`)
		}
		c.Redirect(303, s.URL)
		// c.JSON(200, gin.H{
		// 	"url": s.URL,
		// })
	})
}
