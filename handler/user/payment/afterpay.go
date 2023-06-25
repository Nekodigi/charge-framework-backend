package payment

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v74"
)

func (p *Payment) HandleAfterPay(e *gin.Engine) {
	e.POST("/afterpay", func(c *gin.Context) {
		stripe.Key = p.StripeSecret

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
			mode := metadata["mode"].(string)
			if serviceId == "" || userId == "" {
				log.Fatalf("Invalid subscribe: %v", metadata)
				c.Status(http.StatusBadRequest)
				return
			}
			if mode == "payment" {
				fmt.Println("pay tx ready!")
				ctx := context.Background()
				err := p.Fs.Client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
					user := p.Fs.GetUserByIdTx(tx, serviceId, userId)
					quota, _ := strconv.ParseFloat(metadata["quota"].(string), 64)
					user.RemainQuota += quota
					p.Fs.UpdateUserTx(tx, user)
					fmt.Println(serviceId, userId, quota)
					fmt.Println("Payment was created!")
					return nil
				})
				if err != nil {
					log.Printf("Failed to process payment: %v\n", err)
					c.Status(http.StatusBadRequest)
					return
				}
				return
			} else {
				planId := metadata["plan_id"].(string)
				fmt.Println("sub tx ready!")
				ctx := context.Background()
				err := p.Fs.Client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
					user := p.Fs.GetUserByIdTx(tx, serviceId, userId)
					user.Subscription = event.Data.Object["subscription"].(string)
					service := p.Fs.GetServiceByIdTx(tx, serviceId)
					user.CustomerId = event.Data.Object["customer"].(string)
					user.Plan = planId
					fmt.Println(serviceId, service)
					plan := service.Plan[user.Plan]
					user.AllocQuota = plan.Quota
					p.Fs.UpdateUserTx(tx, user)
					fmt.Println(user, service)
					fmt.Println("Subscription was created!")
					return nil
				})
				if err != nil {
					log.Printf("Failed to subscribe user: %v\n", err)
					c.Status(http.StatusBadRequest)
					return
				}
			}

		// ... handle other event types
		case "invoice.payment_succeeded": //quota update for free plan will be delegated
			if event.Data.Object["subscription"].(string) == "" {
				fmt.Println("Not subscription")
				return
			}
			ctx := context.Background()
			p.Fs.Client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
				fmt.Println(event.Data.Object["metadata"])
				user := p.Fs.GetUserBySubIdTx(tx, event.Data.Object["subscription"].(string))
				service := p.Fs.GetServiceByIdTx(tx, user.ServiceId)
				user.RemainQuota = user.AllocQuota
				service.RemainQuota += user.AllocQuota * service.Plan[user.Plan].QuotaLeak
				p.Fs.UpdateUserTx(tx, user)
				p.Fs.UpdateServiceTx(tx, service)
				fmt.Println("PaymentIntent was successful!")
				return nil
			})

		case "customer.subscription.deleted":
			user := p.Fs.GetUserBySubId(event.Data.Object["id"].(string))
			fmt.Println(user)
			service := p.Fs.GetServiceById(user.ServiceId)
			user.Plan = "free"
			user.AllocQuota = service.Plan["free"].Quota
			user.RemainQuota = user.AllocQuota
			user.Subscription = ""
			p.Fs.UpdateUser(user)
			fmt.Println("Subscription was canceled")
		default:
			fmt.Fprintf(os.Stderr, "Unhandled event type: %s\n", event.Type)
		}

		c.Status(http.StatusOK)
	})
}
