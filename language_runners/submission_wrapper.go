package language_runners

import (
	"Licenta_Processing_Service/dtos"
	"Licenta_Processing_Service/repositories"

	"fmt"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type SubmissionWrapperConf struct {
	FileRepository  *repositories.FilesRepository
	DbRepo          *repositories.PostgresSQLRepo
	S3Repo          *repositories.S3Repository
	LanguageRunners map[dtos.ProgrammingLanguage]LanguageRunner
}

type SubmissionWrapper struct {
	FileRepository  *repositories.FilesRepository
	DbRepo          *repositories.PostgresSQLRepo
	S3Repo          *repositories.S3Repository
	LanguageRunners map[dtos.ProgrammingLanguage]LanguageRunner
}

func NewSubmissionWrapper(config *SubmissionWrapperConf) *SubmissionWrapper {
	return &SubmissionWrapper{
		FileRepository:  config.FileRepository,
		DbRepo:          config.DbRepo,
		S3Repo:          config.S3Repo,
		LanguageRunners: config.LanguageRunners,
	}
}

//var submission = dtos.Submission{
//	Id:                  uuid.New(),
//	ProblemID:           uuid.MustParse("b1218de1-f0a8-4552-8b28-29009882ac63"),
//	UserId:              1,
//	ProgrammingLanguage: "Java",
//	TestResults:         nil,
//}

//var problem = dtos.Problem{
//	Id:                uuid.New(),
//	ProblemDifficulty: "Hard",
//	ProblemStatement:  "Da",
//	ProblemTitle:      "Yes",
//	TimeLimit:         3000000,
//	MemoryLimit:       1000000,
//	TestCases: []dtos.TestCase{
//		{
//			Id:                     uuid.New(),
//			InputFileName:          "inputs/test1",
//			ExpectedOutputFileName: "expected/ref1",
//		},
//		{
//			Id:                     uuid.New(),
//			InputFileName:          "inputs/test2",
//			ExpectedOutputFileName: "expected/ref2",
//		},
//		{
//			Id:                     uuid.New(),
//			InputFileName:          "inputs/test3",
//			ExpectedOutputFileName: "expected/ref3",
//		},
//	},
//}

func (submissionWrapper *SubmissionWrapper) RunSubmission(submissionId uuid.UUID) error {
	//We take the submission from the database
	submission, err := submissionWrapper.DbRepo.GetSubmission(submissionId.String())
	if err != nil {
		return errors.Wrapf(err, "Could not get the submission: %s from the database", submissionId)
	}
	//fmt.Println(submission)
	//We get the problem from the database
	problem, err := submissionWrapper.DbRepo.GetProblem(submission.ProblemID.String())
	if err != nil {
		return errors.Wrapf(err, "Could not get the problem: %s from the database", submission.ProblemID)
	}

	//We take the submission from the aws s3 repository
	s3Submission, err := submissionWrapper.S3Repo.GetSubmission(submission.ProblemID.String(), submissionId.String())
	if err != nil {
		return errors.Wrapf(err, "Error trying to download submission: %s from s3", submissionId)
	}

	if err = submissionWrapper.S3Repo.DownloadTests(problem.Id.String()); err != nil {
		return errors.Wrapf(err, "Could sync the files for the problem: %s", problem.Id.String())
	}

	//We set the correct code runner, according to the programming language that the solution was written in
	submissionRunner, mapError := submissionWrapper.LanguageRunners[submission.ProgrammingLanguage]
	if mapError != true {
		return fmt.Errorf("%s is not supported as a programming language", submission.ProgrammingLanguage)
	}

	tests, err := submissionWrapper.DbRepo.GetTests(problem.Id.String())

	if err != nil {
		return fmt.Errorf("couldn't get tests cases for problem: %s", problem.Id.String())
	}
	solutionReq := &dtos.SolutionRequest{
		File:        s3Submission,
		Submission:  *submission,
		Tests:       tests,
		TimeOut:     problem.TimeLimit,
		MemoryLimit: problem.MemoryLimit,
	}
	testResults, err := submissionRunner.RunSubmission(solutionReq)
	if err != nil {
		return err
	}

	submissionWrapper.DbRepo.SaveTestResults(testResults)

	return nil

}
