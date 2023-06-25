package quota

import (
	"log"
	"net/http"
	"strconv"

	"github.com/Nekodigi/charge-framework-backend/consts"
	"github.com/gin-gonic/gin"
)

func (q *Quota) HandleAddQuota(e *gin.Engine) {
	e.PUT("/quota/add/:service_id/:user_id", func(c *gin.Context) {
		serviceId := c.Param("service_id")
		userId := c.Param("user_id")
		amount, err := strconv.ParseFloat(c.Query("amount"), 64)
		if err != nil {
			log.Println(err)
		}
		if serviceId == "" || c.Query("amount") == "" {
			c.Status(http.StatusBadRequest)
			return
		}

		service := q.Fs.GetServiceById(serviceId)
		if userId != "" {
			user := q.Fs.GetUserById(serviceId, userId)
			user.RemainQuota += amount
			q.Fs.UpdateUser(user)
			c.JSON(200, gin.H{
				"status": consts.USER_QUOTA_UPDATED,
			})
		} else {
			service.RemainQuota += amount
			q.Fs.UpdateService(service)
			c.JSON(200, gin.H{
				"status": consts.SERVICE_QUOTA_UPDATED,
			})
		}
	})
}
