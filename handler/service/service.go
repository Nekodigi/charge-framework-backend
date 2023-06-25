package service

import (
	infraFirestore "github.com/Nekodigi/charge-framework-backend/infrastructure/firestore"
	"github.com/gin-gonic/gin"
)

type (
	Service struct {
		Fs *infraFirestore.Firestore
	}
)

func (s *Service) Handle(e *gin.Engine) {
	s.HandleGetPlan(e)
	s.HandleGetService(e)
	s.HandleGetServiceList(e)
}
