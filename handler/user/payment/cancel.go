package payment

import (
	"fmt"
	"net/http"

	"github.com/Nekodigi/charge-framework-backend/consts"
	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/subscription"
)

func (p *Payment) HandleCancel(e *gin.Engine) {
	e.PUT("/subscription/cancel/:service_id/:user_id", func(c *gin.Context) {
		serviceId := c.Param("service_id")
		userId := c.Param("user_id")
		if serviceId == "" || userId == "" {
			c.Status(http.StatusBadRequest)
			return
		}

		user := p.Fs.GetUserById(serviceId, userId)

		if user.Subscription == "" {
			c.Status(http.StatusBadRequest)
			return
		}

		s, err := subscription.Get(user.Subscription, nil)

		if s.CancelAtPeriodEnd == true {
			c.Status(http.StatusBadRequest)
			return
		}

		params := &stripe.SubscriptionParams{CancelAtPeriodEnd: stripe.Bool(true)}
		_, err = subscription.Update(user.Subscription, params)

		if err != nil {
			fmt.Errorf("Error canceling subscription: %v", err)
			c.JSON(400, gin.H{
				"status": consts.FAILED,
			})
		} else {
			fmt.Printf("Subscription canceled!")
			c.JSON(200, gin.H{
				"status": consts.SUBSCRIPTION_CANCEL_AT_PERIOD_END,
			})
		}
	})
}
