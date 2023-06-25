package user

import (
	"github.com/Nekodigi/charge-framework-backend/handler/user/payment"
	"github.com/Nekodigi/charge-framework-backend/handler/user/quota"
	infraFirestore "github.com/Nekodigi/charge-framework-backend/infrastructure/firestore"
	"github.com/Nekodigi/charge-framework-backend/infrastructure/stripe"
	"github.com/gin-gonic/gin"
)

type (
	User struct {
		StripeSecret string
		Fs           *infraFirestore.Firestore
		St           *stripe.Stripe
	}
)

func (u *User) Handle(e *gin.Engine) {
	(&payment.Payment{u.StripeSecret, u.Fs, u.St}).Handle(e)
	(&quota.Quota{u.Fs}).Handle(e)
	u.HandleUpdateRedirect(e)
	u.HandleGetRedirect(e)
	u.HandleGetUser(e)
	u.HandleGetUniqueId(e)
}
