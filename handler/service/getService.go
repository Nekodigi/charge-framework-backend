package service

import "github.com/gin-gonic/gin"

func (s *Service) HandleGetService(e *gin.Engine) {
	e.GET("/service/:service_id", func(c *gin.Context) {
		serviceId := c.Param("service_id")
		service := s.Fs.GetServiceById(serviceId)
		c.JSON(200, gin.H{"name": service.Name})
	})
}
