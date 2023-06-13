package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	infraFirestore "github.com/Nekodigi/charge-framework-backend/infrastructure/firestore"
	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v74"
)

type (
	Afterpay struct {
		stripeSecret string
		fs           *infraFirestore.Firestore
	}
)

func (a *Afterpay) Handle(e *gin.Engine) {
	e.POST("/afterpay", func(c *gin.Context) {
		stripe.Key = a.stripeSecret

		body, _ := ioutil.ReadAll(c.Request.Body)
		//log.Println("body = ", string(body))

		event := stripe.Event{}

		if err := json.Unmarshal(body, &event); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to parse webhook body json: %v\n", err.Error())
			c.Status(http.StatusBadRequest)
			return
		}

		fmt.Println(event.Type, event.Data.Object["metadata"], event.Data.Object["subscription"], event.Data.Object["id"]) //metadata will be handed when checkout.session.completed

		// Unmarshal the event data into an appropriate struct depending on its Type
		switch event.Type {
		case "checkout.session.completed":
			metadata := event.Data.Object["metadata"].(map[string]interface{})
			serviceId := metadata["service_id"].(string)
			userId := metadata["user_id"].(string)
			planId := metadata["plan_id"].(string)
			user := a.fs.GetUserById(serviceId, userId)

			user.Subscription = event.Data.Object["subscription"].(string)
			service := a.fs.GetServiceById(serviceId)
			user.Plan = planId
			fmt.Println(serviceId, service)
			plan := service.Plan[user.Plan]
			user.AllocQuota = plan.Quota
			a.fs.UpdateUser(user)
			service.RemainQuota += user.AllocQuota * plan.QuotaLeak
			fmt.Println("Checkout was successful!")
		// ... handle other event types
		case "invoice.payment_succeeded": //quota update for free plan will be delegated
			user := a.fs.GetUserBySubId(event.Data.Object["subscription"].(string))
			service := a.fs.GetServiceById(user.ServiceId)
			user.RemainQuota = user.AllocQuota
			service.RemainQuota += user.AllocQuota * service.Plan[user.Plan].QuotaLeak
			a.fs.UpdateUser(user)
			fmt.Println("PaymentIntent was successful!")
		case "customer.subscription.deleted":
			user := a.fs.GetUserBySubId(event.Data.Object["id"].(string))
			service := a.fs.GetServiceById(user.ServiceId)
			user.Plan = "free"
			user.AllocQuota = service.Plan["free"].Quota
			user.RemainQuota = user.AllocQuota
			a.fs.UpdateUser(user)
			fmt.Println("Subscription was canceled")
		default:
			fmt.Fprintf(os.Stderr, "Unhandled event type: %s\n", event.Type)
		}

		c.Status(http.StatusOK)
	})
}
