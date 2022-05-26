package main

import (
	"Licenta_Processing_Service/dtos"
	"Licenta_Processing_Service/language_runners"
	"Licenta_Processing_Service/repositories"
	"Licenta_Processing_Service/services"
)

func main() {

	fileRepository, _ := repositories.NewFileRepository("problems")
	dbRepo := repositories.NewPostgresSQLRepo()
	submissionRunner := language_runners.NewSubmissionWrapper(&language_runners.SubmissionWrapperConf{
		FileRepository: fileRepository,
		DbRepo:         dbRepo,
		S3Repo: repositories.NewS3Repository(dtos.AWSConfig{
			AWSRegion:     "eu-central-1",
			AWSBucketName: "vladbucket123",
			BaseLocalDir:  "problems",
		}),
		LanguageRunners: initLanguageRunners(fileRepository),
	})

	rmqConf := &services.RabbitMQConfig{
		Username:          "pgisyrij",
		Password:          "Knh0-TwtXSPvv_1lqJoC-u92ZfHzFaVk",
		HostName:          "roedeer.rmq.cloudamqp.com",
		Port:              5672,
		QueueName:         "YES",
		SubmissionWrapper: submissionRunner,
	}

	rmq := services.NewRabbitMQConsumer(rmqConf)
	err := rmq.StartConnection()
	if err != nil {
		return
	}
	err = rmq.AcceptMessages()
	if err != nil {
		return
	}
}

func initLanguageRunners(fileRepo *repositories.FilesRepository) map[dtos.ProgrammingLanguage]language_runners.LanguageRunner {

	return map[dtos.ProgrammingLanguage]language_runners.LanguageRunner{
		dtos.C:       language_runners.NewCSubmissionRunner(fileRepo),
		dtos.Java:    language_runners.NewJavaSubmissionRunner(fileRepo),
		dtos.Python3: language_runners.NewPythonSubmissionRunner(fileRepo),
	}
}
