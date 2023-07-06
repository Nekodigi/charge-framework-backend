package quota

import (
	"net/http"

	"github.com/Nekodigi/charge-framework-backend/consts"
	"github.com/gin-gonic/gin"
)

func (q *Quota) HandleCheckQuota(e *gin.Engine) {
	e.GET("/quota/:service_id/:user_id", func(c *gin.Context) {
		serviceId := c.Param("service_id")
		userId := c.Param("user_id")
		if serviceId == "" || userId == "" {
			c.Status(http.StatusBadRequest)
			return
		}

		user := q.Fs.GetUserById(serviceId, userId)

		c.JSON(200, gin.H{
			"allocQuota":  user.AllocQuota,
			"remainQuota": user.RemainQuota,
			"status":      consts.OK,
		})

	})
}
