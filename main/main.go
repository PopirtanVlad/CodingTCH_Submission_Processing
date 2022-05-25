package main

import (
	"Licenta_Processing_Service/dtos"
	"Licenta_Processing_Service/language_runners"
	"Licenta_Processing_Service/repositories"
	"Licenta_Processing_Service/services"
)

//func RunJava() {
//	repository, err := repositories.NewFileRepository("java_test")
//
//	if err != nil {
//		panic(err)
//	}
//
//	file, err := repository.OpenFile("test_dir", "Solution.java")
//	if err != nil {
//		panic(err)
//	}
//
//	language_runners.NewJavaSubmissionRunner(repository).RunSubmission(&dtos.SolutionRequest{
//		File: file,
//		Submission: dtos.Submission{
//			Id:                  uuid.New(),
//			ProblemID:           uuid.MustParse("9994ba64-a1ff-44ca-afc0-7410da8bf48e"),
//			UserId:              1,
//			ProgrammingLanguage: "Java",
//			TestResults:         nil,
//		},
//		Tests: []dtos.TestCase{
//			{
//				Id:                     uuid.New(),
//				InputFileName:          "inputs/test1",
//				ExpectedOutputFileName: "expected/ref1",
//			},
//			{
//				Id:                     uuid.New(),
//				InputFileName:          "inputs/test2",
//				ExpectedOutputFileName: "expected/ref2",
//			},
//			{
//				Id:                     uuid.New(),
//				InputFileName:          "inputs/test3",
//				ExpectedOutputFileName: "expected/ref3",
//			},
//		},
//	})
//}
//
//func RunC() {
//	repository, err := repositories.NewFileRepository("java_test")
//
//	if err != nil {
//		panic(err)
//	}
//
//	file, err := repository.OpenFile("test_dir", "Solution.c")
//	if err != nil {
//		panic(err)
//	}
//
//	language_runners.NewCSubmissionRunner(repository).RunSubmission(&dtos.SolutionRequest{
//		File: file,
//		Submission: dtos.Submission{
//			Id:                  uuid.New(),
//			ProblemID:           uuid.MustParse("9994ba64-a1ff-44ca-afc0-7410da8bf48e"),
//			UserId:              1,
//			ProgrammingLanguage: "Java",
//			TestResults:         nil,
//		},
//		Tests: []dtos.TestCase{
//			{
//				Id:                     uuid.New(),
//				InputFileName:          "inputs/test1",
//				ExpectedOutputFileName: "expected/ref1",
//			},
//			{
//				Id:                     uuid.New(),
//				InputFileName:          "inputs/test2",
//				ExpectedOutputFileName: "expected/ref2",
//			},
//			{
//				Id:                     uuid.New(),
//				InputFileName:          "inputs/test3",
//				ExpectedOutputFileName: "expected/ref3",
//			},
//		},
//	})
//}
//
//func RunPy() {
//	repository, err := repositories.NewFileRepository("java_test")
//
//	if err != nil {
//		panic(err)
//	}
//
//	file, err := repository.OpenFile("test_dir", "Solution.py")
//	if err != nil {
//		panic(err)
//	}
//
//	language_runners.NewPythonSubmissionRunner(repository).RunSubmission(&dtos.SolutionRequest{
//		File: file,
//		Submission: dtos.Submission{
//			Id:                  uuid.New(),
//			ProblemID:           uuid.MustParse("9994ba64-a1ff-44ca-afc0-7410da8bf48e"),
//			UserId:              1,
//			ProgrammingLanguage: "Java",
//			TestResults:         nil,
//		},
//		Tests: []dtos.TestCase{
//			{
//				Id:                     uuid.New(),
//				InputFileName:          "inputs/test1",
//				ExpectedOutputFileName: "expected/ref1",
//			},
//			{
//				Id:                     uuid.New(),
//				InputFileName:          "inputs/test2",
//				ExpectedOutputFileName: "expected/ref2",
//			},
//			{
//				Id:                     uuid.New(),
//				InputFileName:          "inputs/test3",
//				ExpectedOutputFileName: "expected/ref3",
//			},
//		},
//	})
//}

func main() {

	submissionRunner := language_runners.NewSubmissionWrapper(&language_runners.SubmissionWrapperConf{
		FileRepository:  nil,
		DbRepo:          nil,
		S3Repo:          repositories.NewS3Repository(),
		LanguageRunners: initLanguageRunners(),
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

func initLanguageRunners() map[dtos.ProgrammingLanguage]language_runners.LanguageRunner {
	fileRepository, _ := repositories.NewFileRepository("test")

	return map[dtos.ProgrammingLanguage]language_runners.LanguageRunner{
		dtos.C:       language_runners.NewCSubmissionRunner(fileRepository),
		dtos.Java:    language_runners.NewJavaSubmissionRunner(fileRepository),
		dtos.Python3: language_runners.NewPythonSubmissionRunner(fileRepository),
	}
}
