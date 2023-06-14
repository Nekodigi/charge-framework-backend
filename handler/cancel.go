package handler

import (
	"fmt"

	infraFirestore "github.com/Nekodigi/charge-framework-backend/infrastructure/firestore"
	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/sub"
)

type (
	Cancel struct {
		fs *infraFirestore.Firestore
	}
)

func (q *Cancel) Handle(e *gin.Engine) {
	e.POST("/cancel", func(c *gin.Context) {
		serviceId := c.Query("service_id")
		userId := c.Query("user_id")

		user := q.fs.GetUserById(serviceId, userId)

		_, err := sub.Cancel(
			user.Subscription,
			nil,
		)
		if err != nil {
			fmt.Errorf("Error canceling subscription: %v", err)
			c.JSON(400, gin.H{
				"message": "FAILED",
			})
		} else {
			fmt.Printf("Subscription canceled!")
			c.JSON(200, gin.H{
				"message": "OK",
			})
		}
	})
}
