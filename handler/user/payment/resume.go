package payment

import (
	"fmt"
	"net/http"

	"github.com/Nekodigi/charge-framework-backend/consts"
	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/subscription"
)

func (p *Payment) HandleResume(e *gin.Engine) {
	e.PUT("/subscription/resume/:service_id/:user_id", func(c *gin.Context) {
		serviceId := c.Param("service_id")
		userId := c.Param("user_id")
		if serviceId == "" || userId == "" {
			fmt.Printf("serviceId:%s, userId:%s\n", serviceId, userId)
			c.Status(http.StatusBadRequest)
			return
		}

		user := p.Fs.GetUserById(serviceId, userId)

		if user.Subscription == "" {
			c.Status(http.StatusBadRequest)
			return
		}

		s, err := subscription.Get(user.Subscription, nil)

		if s.CancelAtPeriodEnd == false {
			c.Status(http.StatusBadRequest)
			return
		}

		params := &stripe.SubscriptionParams{CancelAtPeriodEnd: stripe.Bool(false)}
		_, err = subscription.Update(user.Subscription, params)

		if err != nil {
			fmt.Errorf("Error resuming subscription: %v", err)
			c.JSON(400, gin.H{
				"status": consts.FAILED,
			})
		} else {
			fmt.Printf("Subscription resumed!")
			c.JSON(200, gin.H{
				"status": consts.SUBSCRIPTION_RESUME_AT_PERIOD_END,
			})
		}
	})
}
