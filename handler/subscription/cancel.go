package subscription

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/sub"
)

func (q *Subscription) HandleCancel(e *gin.Engine) {
	e.POST("/cancel/:service_id/:user_id", func(c *gin.Context) {
		serviceId := c.Param("service_id")
		userId := c.Param("user_id")
		if serviceId == "" || userId == "" {
			c.Status(http.StatusBadRequest)
			return
		}

		user := q.Fs.GetUserById(serviceId, userId)

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
