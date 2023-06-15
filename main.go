package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"cloud.google.com/go/firestore"
	"github.com/Nekodigi/charge-framework-backend/handler"
	infraFirestore "github.com/Nekodigi/charge-framework-backend/infrastructure/firestore"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// load .env file
	if len(os.Args) == 2 && os.Args[1] == "test" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
		serviceId := "test"
		userId := "u0007"
		planId := "basic"
		// planId := "basic"
		fs := infraFirestore.NewFirestore()
		fmt.Println("tx ready!")
		ctx := context.Background()
		fs.Client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
			fmt.Println("tx start!")
			user := fs.GetUserByIdTx(tx, serviceId, userId)
			user.Subscription = ""
			service := fs.GetServiceByIdTx(tx, serviceId)
			user.Plan = planId
			fmt.Println(serviceId, service)
			plan := service.Plan[user.Plan]
			user.AllocQuota = plan.Quota
			fs.UpdateUserTx(tx, user)
			fmt.Println(user, service)
			fmt.Println("Subscription was created!")
			return nil
		})
		//ops.IssuePromo()
	} else {

		// set up server
		engine := gin.Default()
		handler.Router(engine)
		engine.Run(":8080")
	}
}
