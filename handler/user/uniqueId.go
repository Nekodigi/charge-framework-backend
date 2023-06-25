package user

import (
	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v74"
)

var ()

func (u *User) HandleGetUniqueId(e *gin.Engine) {
	// our basic charge API route
	stripe.Key = u.StripeSecret
	e.GET("/user/unique_id", func(c *gin.Context) {

		c.JSON(200, gin.H{
			"id": u.Fs.CreateUniqueID(),
		})
	})
}
