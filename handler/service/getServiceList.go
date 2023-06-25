package service

import "github.com/gin-gonic/gin"

func (s *Service) HandleGetServiceList(e *gin.Engine) {
	e.GET("/service/list", func(c *gin.Context) {
		services := s.Fs.GetServiceList()
		var list []string
		for _, service := range services {
			list = append(list, service.Id)
		}
		c.JSON(200, gin.H{"list": list})
	})
}
