package dtos

type ServiceConfig struct {
	RabbitMQConfig
	AWSConfig
	PostgresSQLConfig
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

type PostgresSQLConfig struct {
	PostgresDialect  string
	PostgresHost     string
	PostgresDBport   int
	PostgresUser     string
	PostgresName     string
	PostgresPassword string
}
