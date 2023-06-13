package ops

import (
	"os"

	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/coupon"
	"github.com/stripe/stripe-go/v74/promotioncode"
)

func IssuePromo() string {
	stripeSecret := os.Getenv("SK_TEST_KEY")
	stripe.Key = stripeSecret

	//* EDIT param; currently 100% off forever
	cparam := &stripe.CouponParams{
		Duration:   stripe.String(string(stripe.CouponDurationForever)),
		PercentOff: stripe.Float64(100),
	}
	c, _ := coupon.New(cparam)

	promo := &stripe.PromotionCodeParams{
		Coupon: stripe.String(c.ID),
	}
	pc, _ := promotioncode.New(promo)
	return pc.Code
}
