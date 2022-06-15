package services

import (
	"Licenta_Processing_Service/language_runners"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

type RabbitMQConfig struct {
	Username          string
	Password          string
	HostName          string
	Port              int
	QueueName         string
	SubmissionWrapper *language_runners.SubmissionWrapper
}

type RabbitMQConsumer struct {
	username          string
	password          string
	hostName          string
	port              int
	queueName         string
	channel           *amqp.Channel
	connection        *amqp.Connection
	SubmissionWrapper *language_runners.SubmissionWrapper
}

func NewRabbitMQConsumer(rabbbitMQConf *RabbitMQConfig) *RabbitMQConsumer {
	return &RabbitMQConsumer{
		username:          rabbbitMQConf.Username,
		password:          rabbbitMQConf.Password,
		hostName:          rabbbitMQConf.HostName,
		port:              rabbbitMQConf.Port,
		queueName:         rabbbitMQConf.QueueName,
		SubmissionWrapper: rabbbitMQConf.SubmissionWrapper,
	}
}

func (rabbitMQConsumer *RabbitMQConsumer) StartConnection() error {
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

func (rabbitMQConsumer *RabbitMQConsumer) StopConnection() error {
	err := rabbitMQConsumer.connection.Close()

	if err != nil {
		return err
	}
	logrus.Info("Connection cancelled successfully")
	return nil
}

func (rabbitMQConsumer *RabbitMQConsumer) AcceptMessages() error {
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
			rabbitMQConsumer.handleMessageReceived(string(message.Body))
		}
	}()

	<-infiniteChannel
	return nil
}

func (rabbitMQConsumer *RabbitMQConsumer) handleMessageReceived(message string) {
	err := rabbitMQConsumer.SubmissionWrapper.RunSubmission(message)
	if err != nil {
		panic(err)
	}
}
