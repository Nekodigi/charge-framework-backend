package subscription

import (
	infraFirestore "github.com/Nekodigi/charge-framework-backend/infrastructure/firestore"
	"github.com/gin-gonic/gin"
)

type (
	Subscription struct {
		StripeSecret string
		Fs           *infraFirestore.Firestore
	}
)

func (s *Subscription) Handle(e *gin.Engine) {
	s.HandleAfterPay(e)
	s.HandlePurchase(e)
	s.HandleSubscribe(e)
	s.HandleCancel(e)
}
