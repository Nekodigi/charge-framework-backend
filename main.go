package main

import (
	"log"
	"os"

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
		infraFirestore.NewFirestore().GetServiceById("test")
		//ops.IssuePromo()
	} else {

		// set up server
		engine := gin.Default()
		handler.Router(engine)
		engine.Run(":8080")
	}
}
