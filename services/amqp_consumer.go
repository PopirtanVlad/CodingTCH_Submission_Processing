package services

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"log"
)

type RabbitMQConfig struct {
	Username  string
	Password  string
	HostName  string
	Port      int
	QueueName string
}

type RabbitMQConsumer struct {
	username   string
	password   string
	hostName   string
	port       int
	queueName  string
	channel    *amqp.Channel
	connection *amqp.Connection
}

func NewRabbitMQConsumer(rabbbitMQConf *RabbitMQConfig) *RabbitMQConsumer {
	return &RabbitMQConsumer{
		username:  rabbbitMQConf.Username,  //username
		password:  rabbbitMQConf.Password,  //password
		hostName:  rabbbitMQConf.HostName,  //host
		port:      rabbbitMQConf.Port,      //port
		queueName: rabbbitMQConf.QueueName, //port
	}
}

func (rabbitMQConsumer *RabbitMQConsumer) StartConnection(ctx context.Context) error {
	//connectionURI := amqp.URI{Username: "guest",
	//	Password: "guest",
	//	Host:     "localhost",
	//	Port:     15672,
	//	Scheme:   "http",
	//}
	connection, err := amqp.Dial("amqps://pgisyrij:Knh0-TwtXSPvv_1lqJoC-u92ZfHzFaVk@roedeer.rmq.cloudamqp.com/pgisyrij")

	if err != nil {
		return err
	}

	rabbitMQConsumer.connection = connection

	logrus.WithFields(logrus.Fields{
		"host": rabbitMQConsumer.hostName,
		"port": rabbitMQConsumer.port,
	}).Info("RabbitMQ connection established successfully")

	return nil
}

func (rabbitMQConsumer *RabbitMQConsumer) StopConnection(ctx context.Context) error {
	err := rabbitMQConsumer.connection.Close()

	if err != nil {
		return err
	}
	logrus.Info("Connection cancelled successfully")
	return nil
}

func (rabbitMQConsumer *RabbitMQConsumer) AcceptMessages(ctx context.Context) error {
	ch, err := rabbitMQConsumer.connection.Channel()

	ch.Qos(1, 0, false)
	if err != nil {
		return err
	}

	ch.QueueDeclare(
		rabbitMQConsumer.queueName,
		true,
		false,
		false,
		false,
		nil,
	)

	messages, err := ch.Consume(
		rabbitMQConsumer.queueName, //queue name
		"",                         // consumer
		true,                       //auto-ack
		false,                      //exclusive
		false,                      //no local
		false,                      //no wait
		nil,                        //arguments
	)

	if err != nil {
		return nil
	}

	logrus.Info("Succesfully connected ro RabbitMQ")

	infiniteChannel := make(chan bool)

	go func() {
		for message := range messages {
			HandleMessageReceived(string(message.Body))
		}
	}()

	<-infiniteChannel
	return nil
}

func HandleMessageReceived(message string) {

	log.Printf("RECEIVED MESSAGE: %s", message)
}
