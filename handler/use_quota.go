package handler

import (
	"log"
	"strconv"

	infraFirestore "github.com/Nekodigi/charge-framework-backend/infrastructure/firestore"
	"github.com/gin-gonic/gin"
)

type (
	Quota struct {
		fs *infraFirestore.Firestore
	}
)

func (q *Quota) Handle(e *gin.Engine) {
	e.POST("/use_quota", func(c *gin.Context) {
		serviceId := c.Query("service_id")
		userId := c.Query("user_id")
		amount, err := strconv.ParseFloat(c.Query("amount"), 64)
		if err != nil {
			log.Println(err)
		}

		user := q.fs.GetUserById(serviceId, userId)
		service := q.fs.GetServiceById(serviceId)
		if amount > user.RemainQuota { //update quota if possible
			if !q.fs.UpdateUserQuota(&user) {
				c.JSON(402, gin.H{
					"message": "QUOTA_NOT_ENOUGH",
				})
				return
			}
		} else if user.Plan == "free" && amount > service.RemainQuota {
			if !q.fs.UpdateServiceQuota(&service) {
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
		q.fs.UpdateUser(user)
		q.fs.UpdateService(service)
		c.JSON(200, gin.H{
			"message": "OK",
		})

	})
}
