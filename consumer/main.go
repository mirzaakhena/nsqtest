package main

import (
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/mirzaakhena/nsqtest/config"
	"github.com/mirzaakhena/nsqtest/logger"
	"github.com/mirzaakhena/nsqtest/messagebroker"
)

// MBConsumer is
type MBConsumer struct {
}

func main() {

	mb := MBConsumer{}

	cf := config.NewRealtimeConfig("config", "$GOPATH/src/bitbucket.org/nsqtest/")

	log.GetLog().Info("Consumer.main", "address: %s", cf.GetNSQdURL())

	x := messagebroker.NewConsumer([]messagebroker.ConsumerHandler{
		messagebroker.ConsumerHandler{
			Topic:               "Test_Topic",
			FunctionHandler:     mb.TestRequest,
			NumberOfConcurrency: 1,
		},
	})

	x.Run(cf.GetNSQdURL())

	log.GetLog().Info("main", "waiting message")

	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})
	router.Run(fmt.Sprintf(":%s", "8081"))

}

// TestRequest is
func (o *MBConsumer) TestRequest(m *messagebroker.Context) error {

	var req map[string]string
	if err := json.Unmarshal(m.Message, &req); err != nil {
		log.GetLog().Error("TestRequest", "Fail convert json. %s", err.Error())
		return nil
	}

	log.GetLog().Info("TestRequest", "Received: %v", req)

	return nil
}
