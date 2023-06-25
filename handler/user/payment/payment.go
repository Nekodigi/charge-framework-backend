package payment

import (
	infraFirestore "github.com/Nekodigi/charge-framework-backend/infrastructure/firestore"
	infraStripe "github.com/Nekodigi/charge-framework-backend/infrastructure/stripe"
	"github.com/gin-gonic/gin"
)

type (
	Payment struct {
		StripeSecret string
		Fs           *infraFirestore.Firestore
		St           *infraStripe.Stripe
	}
)

func (p *Payment) Handle(e *gin.Engine) {
	p.HandleAfterPay(e)
	p.HandlePaymentUrl(e)
	p.HandleSubscribeUrl(e)
	p.HandleCancel(e)
	p.HandleResume(e)
}
