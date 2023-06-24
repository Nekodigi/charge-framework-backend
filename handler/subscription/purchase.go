package subscription

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/checkout/session"
)

func (s *Subscription) HandlePurchase(e *gin.Engine) {
	// our basic charge API route
	stripe.Key = s.StripeSecret
	e.GET("/purchase/:service_id/:user_id", func(c *gin.Context) {
		serviceId := c.Param("service_id")
		userId := c.Param("user_id")
		quota, _ := strconv.ParseFloat(c.Query("quota"), 64)
		successUrl := c.Query("success_url")
		cancelUrl := c.Query("cancel_url")

		if serviceId == "" || userId == "" {
			c.Status(http.StatusBadRequest)
			return
		}

		service := s.Fs.GetServiceById(serviceId)

		priceId := service.Plan["free"].PriceId
		mode := "payment"
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
		params.AddMetadata("user_id", userId)
		params.AddMetadata("mode", mode)
		params.AddMetadata("quota", strconv.FormatFloat(quota, 'f', -1, 64))
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
