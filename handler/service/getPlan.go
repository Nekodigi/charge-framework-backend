package service

import "github.com/gin-gonic/gin"

func (s *Service) HandleGetPlan(e *gin.Engine) {
	e.GET("/service/plan/:service_id", func(c *gin.Context) {
		serviceId := c.Param("service_id")
		service := s.Fs.GetServiceById(serviceId)
		c.JSON(200, gin.H{"plan": service.Plan})
	})
}
