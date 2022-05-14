package main

import (
	"Licenta_Processing_Service/services"
	"context"
)

func main() {

	rabbitMQConf := &services.RabbitMQConfig{
		Username:  "pgisyrij",
		Password:  "Knh0-TwtXSPvv_1lqJoC-u92ZfHzFaVk",
		HostName:  "roedeer.rmq.cloudamqp.com",
		Port:      5672,
		QueueName: "YES",
	}

	rabbitMQConsumer := services.NewRabbitMQConsumer(rabbitMQConf)
	rabbitMQConsumer.StartConnection(context.TODO())
	rabbitMQConsumer.AcceptMessages(context.TODO())

}
