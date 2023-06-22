package subscription

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v74"
)

func (a *Subscription) HandleAfterPay(e *gin.Engine) {
	e.POST("/afterpay", func(c *gin.Context) {
		stripe.Key = a.StripeSecret

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
			if serviceId == "" || userId == "" || planId == "" {
				log.Fatalf("Invalid subscribe: %v", metadata)
				c.Status(http.StatusBadRequest)
				return
			}
			fmt.Println("tx ready!")
			ctx := context.Background()
			a.Fs.Client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
				fmt.Println("tx start!")
				user := a.Fs.GetUserByIdTx(tx, serviceId, userId)
				user.Subscription = event.Data.Object["subscription"].(string)
				service := a.Fs.GetServiceByIdTx(tx, serviceId)
				user.Plan = planId
				fmt.Println(serviceId, service)
				plan := service.Plan[user.Plan]
				user.AllocQuota = plan.Quota
				a.Fs.UpdateUserTx(tx, user)
				fmt.Println(user, service)
				fmt.Println("Subscription was created!")
				return nil
			})

		// ... handle other event types
		case "invoice.payment_succeeded": //quota update for free plan will be delegated
			ctx := context.Background()
			a.Fs.Client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
				fmt.Println(event.Data.Object["metadata"])
				user := a.Fs.GetUserBySubIdTx(tx, event.Data.Object["subscription"].(string))
				service := a.Fs.GetServiceByIdTx(tx, user.ServiceId)
				user.RemainQuota = user.AllocQuota
				service.RemainQuota += user.AllocQuota * service.Plan[user.Plan].QuotaLeak
				a.Fs.UpdateUserTx(tx, user)
				a.Fs.UpdateServiceTx(tx, service)
				fmt.Println("PaymentIntent was successful!")
				return nil
			})

		case "customer.subscription.deleted":
			user := a.Fs.GetUserBySubId(event.Data.Object["id"].(string))
			fmt.Println(user)
			service := a.Fs.GetServiceById(user.ServiceId)
			user.Plan = "free"
			user.AllocQuota = service.Plan["free"].Quota
			user.RemainQuota = user.AllocQuota
			a.Fs.UpdateUser(user)
			fmt.Println("Subscription was canceled")
		default:
			fmt.Fprintf(os.Stderr, "Unhandled event type: %s\n", event.Type)
		}

		c.Status(http.StatusOK)
	})
}
