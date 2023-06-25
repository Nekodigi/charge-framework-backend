package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/subscription"
)

var ()

func (u *User) HandleGetUser(e *gin.Engine) {
	// our basic charge API route
	stripe.Key = u.StripeSecret
	e.GET("/user/:service_id/:user_id", func(c *gin.Context) {
		serviceId := c.Param("service_id")
		userId := c.Param("user_id")
		if serviceId == "" || userId == "" {
			c.Status(http.StatusBadRequest)
			return
		}

		user := u.Fs.GetUserById(serviceId, userId)

		s, _ := subscription.Get(user.Subscription, nil)

		//*user should be removed after testing
		c.JSON(200, gin.H{
			"user":                 user,
			"plan":                 user.Plan,
			"cancel_at_period_end": s.CancelAtPeriodEnd,
		})
	})
}
