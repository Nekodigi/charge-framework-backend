package quota

import (
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (q *Quota) HandleUseQuota(e *gin.Engine) {
	e.POST("/use_quota", func(c *gin.Context) {
		serviceId := c.Query("service_id")
		userId := c.Query("user_id")
		amount, err := strconv.ParseFloat(c.Query("amount"), 64)
		if err != nil {
			log.Println(err)
		}

		user := q.Fs.GetUserById(serviceId, userId)
		service := q.Fs.GetServiceById(serviceId)
		if amount > user.RemainQuota { //update quota if possible
			if !q.Fs.UpdateUserQuota(&user) {
				c.JSON(402, gin.H{
					"message": "QUOTA_NOT_ENOUGH",
				})
				return
			}
		} else if user.Plan == "free" && amount > service.RemainQuota {
			if !q.Fs.UpdateServiceQuota(&service) {
				c.JSON(402, gin.H{
					"message": "GLOBAL_QUOTA_NOT_ENOUGH",
				})
				return
			}
		}
		user.RemainQuota -= amount
		if user.Plan == "free" {
			service.RemainQuota -= amount
		}
		q.Fs.UpdateUser(user)
		q.Fs.UpdateService(service)
		c.JSON(200, gin.H{
			"message": "OK",
		})

	})
}
