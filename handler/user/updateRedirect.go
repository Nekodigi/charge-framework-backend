package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v74"
)

func (u *User) HandleUpdateRedirect(e *gin.Engine) {
	// our basic charge API route
	stripe.Key = u.StripeSecret
	e.PUT("/redirect/:service_id/:user_id", func(c *gin.Context) {
		serviceId := c.Param("service_id")
		userId := c.Param("user_id")
		destUserId := c.Query("dest_user")

		if serviceId == "" || userId == "" || destUserId == "" {
			c.Status(http.StatusBadRequest)
			return
		}

		user := u.Fs.GetUserById(serviceId, userId)
		user.Redirect = destUserId
		u.Fs.UpdateUser(user)

		c.Status(http.StatusOK)
	})
}
