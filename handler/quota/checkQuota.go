package quota

import (
	"github.com/gin-gonic/gin"
)

func (q *Quota) HandleCheckQuota(e *gin.Engine) {
	e.POST("/check_quota", func(c *gin.Context) {
		serviceId := c.Query("service_id")
		userId := c.Query("user_id")

		user := q.Fs.GetUserById(serviceId, userId)
		service := q.Fs.GetServiceById(serviceId)

		if 1 >= user.RemainQuota { //update quota if possible
			if !q.Fs.UpdateUserQuota(&user) {
				c.JSON(402, gin.H{
					"allocQuota":  user.AllocQuota,
					"remainQuota": user.RemainQuota,
					"message":     "QUOTA_NOT_ENOUGH",
				})
				return
			}
		} else if user.Plan == "free" && 1 > service.RemainQuota {
			if !q.Fs.UpdateServiceQuota(&service) {
				c.JSON(402, gin.H{
					"allocQuota":  user.AllocQuota,
					"remainQuota": user.RemainQuota,
					"message":     "GLOBAL_QUOTA_NOT_ENOUGH",
				})
				return
			}
		}
		c.JSON(200, gin.H{
			"allocQuota":  user.AllocQuota,
			"remainQuota": user.RemainQuota,
			"message":     "OK",
		})

	})
}
