package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mirzaakhena/nsqtest/config"
	"github.com/mirzaakhena/nsqtest/logger"
	"github.com/mirzaakhena/nsqtest/messagebroker"
)

// MBProducer is
type MBProducer struct {
	MessageBroker messagebroker.IMessageBrokerProducer
}

func main() {

	cf := config.NewRealtimeConfig("config", "$GOPATH/src/bitbucket.org/mirzaakhena/nsqtest/")

	log.GetLog().Info("Producer.main", "address: %s", cf.GetNSQdURL())

	x := MBProducer{
		MessageBroker: messagebroker.NewProducer(cf.GetNSQdURL()),
	}

	router := gin.Default()

	router.GET("/test", x.CallProducer)

	router.Run(":8080")
}

// CallProducer is
func (mb *MBProducer) CallProducer(c *gin.Context) {

	rawData := `{"message": "hello"}
	`
	// publish to message broker
	if err := mb.MessageBroker.Publish("Test_Topic", []byte(rawData)); err != nil {
		log.GetLog().Error("CallProducer", "Failed to publish to message broker %s", err.Error())
	}

	log.GetLog().Info("CallProducer", "raw data %v", rawData)

	c.JSON(http.StatusOK, map[string]string{"message": "done"})
}
