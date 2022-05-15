package dtos

type ServiceConfig struct {
	RabbitMQConfig
	AWSConfig
}

type RabbitMQConfig struct {
	RabbitMQUsername  string
	RabbitMQPassword  string
	RabbitMQHostName  string
	RabbitMQPort      int
	RabbitMQQueueName string
}

type AWSConfig struct {
	AWSRegion     string
	AWSBucketName string
}
