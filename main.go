package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"benginga/components"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	fmt.Println("Benginga...!")

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	router := gin.Default()

	router.GET("/healthz", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	router.POST("/log", components.HandleLog)

	port := os.Getenv("PORT")

	if err != nil {
		log.Fatal(err)
	}

	if port == "" {
		port = "8080"
	}

	address := fmt.Sprintf(":%s", port)

	log.Fatal(router.Run(address))

}
