package main

import (
	"Licenta_Processing_Service/dtos"
	"Licenta_Processing_Service/repositories"
	"Licenta_Processing_Service/services"
	"Licenta_Processing_Service/services/executions"
	uuid "github.com/satori/go.uuid"
)

func main() {

	_ = &services.RabbitMQConfig{
		Username:  "pgisyrij",
		Password:  "Knh0-TwtXSPvv_1lqJoC-u92ZfHzFaVk",
		HostName:  "roedeer.rmq.cloudamqp.com",
		Port:      5672,
		QueueName: "YES",
	}

	repository, err := repositories.NewFileRepository("java_test")

	if err != nil {
		panic(err)
	}

	file, err := repository.OpenFile("49c6db5f-39a1-4647-8b40-a66875d6cc32", "Solution.java")
	if err != nil {
		panic(err)
	}

	executions.NewJavaSubmissionRunner(repository).RunSubmission(dtos.SolutionRequest{
		File: file,
		Solution: dtos.Submission{
			Id:                  uuid.FromStringOrNil("49c6db5f-39a1-4647-8b40-a66875d6cc32"),
			ProblemID:           uuid.FromStringOrNil("9994ba64-a1ff-44ca-afc0-7410da8bf48e"),
			UserId:              1,
			ProgrammingLanguage: "Java",
			TestResults:         nil,
		},
		Tests: []dtos.TestCase{
			{
				Id:                     uuid.FromStringOrNil("49c6db5f-39a1-4647-8b40-a66875d6cc32"),
				InputFileName:          "main/java_test/inputs",
				ExpectedOutputFileName: "main/java_test/expected",
			},
			{
				Id:                     uuid.FromStringOrNil("49c6db5f-39a1-4647-8b40-a66875d6cc32"),
				InputFileName:          "main/java_test/inputs",
				ExpectedOutputFileName: "main/java_test/expected",
			},
		},
	})

}
