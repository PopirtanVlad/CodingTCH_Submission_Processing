package main

import (
	"Licenta_Processing_Service/entities"
	"Licenta_Processing_Service/language_runners"
	"Licenta_Processing_Service/repositories"
	"Licenta_Processing_Service/services"
	"fmt"
)

func main() {

	fileRepository, _ := repositories.NewFileRepository("problems")
	dbRepo := repositories.NewPostgresSQLRepo()
	submissionRunner := language_runners.NewSubmissionWrapper(&language_runners.SubmissionWrapperConf{
		FileRepository: fileRepository,
		DbRepo:         dbRepo,
		S3Repo: repositories.NewS3Repository(entities.AWSConfig{
			AWSRegion:     "eu-central-1",
			AWSBucketName: "lictestbucket1",
			BaseLocalDir:  "problems",
		}),
		LanguageRunners: initLanguageRunners(fileRepository),
	})

	rmqConf := &services.RabbitMQConfig{
		Username:          "pgisyrij",
		Password:          "Knh0-TwtXSPvv_1lqJoC-u92ZfHzFaVk",
		HostName:          "roedeer.rmq.cloudamqp.com",
		Port:              5672,
		QueueName:         "Licenta_Queue",
		SubmissionWrapper: submissionRunner,
	}

	rmq := services.NewRabbitMQConsumer(rmqConf)
	err := rmq.StartConnection()
	if err != nil {
		return
	}
	err = rmq.AcceptMessages()
	if err != nil {
		fmt.Println(err)
		return
	}
}

func initLanguageRunners(fileRepo *repositories.FilesRepository) map[entities.ProgrammingLanguage]language_runners.LanguageRunner {

	return map[entities.ProgrammingLanguage]language_runners.LanguageRunner{
		entities.C:       language_runners.NewCSubmissionRunner(fileRepo),
		entities.Java:    language_runners.NewJavaSubmissionRunner(fileRepo),
		entities.Python3: language_runners.NewPythonSubmissionRunner(fileRepo),
	}
}
