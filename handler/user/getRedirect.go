package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v74"
)

var ()

func (u *User) HandleGetRedirect(e *gin.Engine) {
	// our basic charge API route
	stripe.Key = u.StripeSecret
	e.GET("/redirect/:service_id/:user_id", func(c *gin.Context) {
		serviceId := c.Param("service_id")
		userId := c.Param("user_id")

		if serviceId == "" || userId == "" {
			c.Status(http.StatusBadRequest)
			return
		}

		user := u.Fs.GetUserById(serviceId, userId)

		c.JSON(200, gin.H{
			"id": user.Redirect,
		})
	})
}
