package quota

import (
	"log"
	"net/http"
	"strconv"

	"github.com/Nekodigi/charge-framework-backend/consts"
	"github.com/gin-gonic/gin"
)

func (q *Quota) HandleUseQuota(e *gin.Engine) {
	e.PUT("/quota/use/:service_id/:user_id", func(c *gin.Context) {
		serviceId := c.Param("service_id")
		userId := c.Param("user_id")
		amount, err := strconv.ParseFloat(c.Query("amount"), 64)
		if err != nil {
			log.Println(err)
		}
		if serviceId == "" || userId == "" || c.Query("amount") == "" {
			c.Status(http.StatusBadRequest)
			return
		}

		user := q.Fs.GetUserById(serviceId, userId)
		service := q.Fs.GetServiceById(serviceId)
		if amount > user.RemainQuota { //update quota if possible
			if !q.Fs.UpdateUserQuota(&user) {
				c.JSON(402, gin.H{
					"status": consts.QUOTA_NOT_ENOUGH,
				})
				return
			}
		} else if user.Plan == "free" && amount > service.RemainQuota {
			if !q.Fs.UpdateServiceQuota(&service) {
				c.JSON(402, gin.H{
					"status": consts.GLOBAL_QUOTA_NOT_ENOUGH,
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
			"status": consts.OK,
		})

	})
}
