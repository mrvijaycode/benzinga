package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"benginga/models"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

var (
	//batchSize     int
	//batchInterval time.Duration
	postEndpoint string
	payloads     []models.Payload
	mu           sync.Mutex
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

	router.POST("/log", handleLog)

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

func handleLog(c *gin.Context) {
	var payload models.Payload

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	mu.Lock()
	payloads = append(payloads, payload)
	mu.Unlock()

	batchSize, err := strconv.Atoi(os.Getenv("BATCH_SIZE"))
	if err != nil {
		fmt.Println("Error converting BATCH_SIZE to int or nil")
	}

	logrus.WithFields(logrus.Fields{
		"user_id":    payload.UserID,
		"batch_size": batchSize,
		"payloadLen": len(payloads),
	}).Info("payload received")

	if len(payloads) >= batchSize {
		logrus.Info("Sending logs batch size reached...")
		sendBatch()
	} else {
		logrus.Info("Batch size not reached, waiting...")
	}

	batchInterval, err := strconv.Atoi(os.Getenv("BATCH_INTERVAL"))
	if err != nil {
		fmt.Println("Error converting BATCH_SIZE to int or nil")
	}

	go func() {
		time.Sleep(time.Duration(batchInterval) * time.Minute)
		if len(payloads) > 0 {
			logrus.Info("Sending logs in given interval...")
			sendBatch()
		}
	}()
}

func sendBatch() {
	mu.Lock()
	batch := payloads
	payloads = nil
	mu.Unlock()

	if len(batch) == 0 {
		return
	}

	data, _ := json.Marshal(batch)

	dataReader := bytes.NewReader(data)

	logrus.WithFields(logrus.Fields{
		//"data": string(data),
	}).Info("Records to send ###")

	postEndpoint = os.Getenv("POST_ENDPOINT")

	for i := 0; i < 3; i++ {

		resp, err := http.Post(postEndpoint, "application/json", dataReader)
		if err != nil {
			logrus.WithError(err).Error("failed to send batch")
			time.Sleep(2 * time.Second)
			continue
		}

		logrus.WithFields(logrus.Fields{
			"batch_size":  len(batch),
			"status_code": resp.StatusCode,
		}).Info("batch sent")
		return
	}

	logrus.Fatal("failed to send batch after 3 attempts, exiting")
}
