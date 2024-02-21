package components

import (
	"benginga/models"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var (
	//postEndpoint string
	payloadsChan = make(chan []models.Payload, 100)
	payloads     = make([]models.Payload, 0, 100)
	//mu           sync.Mutex
)

func HandleLog(c *gin.Context) {
	var payload models.Payload

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	payloads = append(payloads, payload)

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

		payloadsChan <- payloads

		go sendBatch()
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
			go sendBatch()
		}
	}()
}

func sendBatch() {
	for loads := range payloadsChan {
		fmt.Println("sending batch")
		if len(loads) > 0 {

			data, _ := json.Marshal(loads)

			dataReader := bytes.NewReader(data)

			logrus.WithFields(logrus.Fields{}).Info("Records to send ###")

			postEndpoint := os.Getenv("POST_ENDPOINT")

			payloads = make([]models.Payload, 0, 100)

			for i := 1; i <= 3; i++ {

				resp, err := http.Post(postEndpoint, "application/json", dataReader)
				if err != nil {
					logrus.WithError(err).Error("failed to send batch")
					time.Sleep(2 * time.Second)
					continue
				}

				logrus.WithFields(logrus.Fields{
					"batch_size":  len(loads),
					"status_code": resp.StatusCode,
				}).Info("batch sent")
				return
			}

			logrus.Fatal("failed to send batch after 3 attempts, exiting")
		}
	}
}
