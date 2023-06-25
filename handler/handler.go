package handler

import (
	"fmt"
	"net/http"
	"os"

	"github.com/Nekodigi/charge-framework-backend/handler/service"
	"github.com/Nekodigi/charge-framework-backend/handler/user"
	infraFirestore "github.com/Nekodigi/charge-framework-backend/infrastructure/firestore"
	infraStripe "github.com/Nekodigi/charge-framework-backend/infrastructure/stripe"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/stripe/stripe-go"
)

var (
	stripeSecret string
	fs           *infraFirestore.Firestore
	st           *infraStripe.Stripe
)

func init() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	stripeSecret = os.Getenv("SK_TEST_KEY")
	stripe.Key = stripeSecret
	fmt.Println(stripeSecret)
	fs = infraFirestore.NewFirestore()
	st = infraStripe.NewStripe(stripeSecret)
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func Router(e *gin.Engine) {
	e.Use(CORSMiddleware())
	(&user.User{stripeSecret, fs, st}).Handle(e)
	(&service.Service{fs}).Handle(e)
	e.GET("/ping", func(ctx *gin.Context) { ctx.String(http.StatusOK, "pong") })
}
