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

var submission = dtos.Submission{
	Id:                  uuid.New(),
	ProblemID:           uuid.MustParse("5f7f14a6-23bd-467f-b5b3-d306bb06af44"),
	UserId:              1,
	ProgrammingLanguage: "Java",
	TestResults:         nil,
}

var problem = dtos.Problem{
	Id:                uuid.New(),
	ProblemDifficulty: "Hard",
	ProblemStatement:  "Da",
	ProblemTitle:      "Yes",
	TimeLimit:         3000000,
	MemoryLimit:       1000000,
	TestCases: []dtos.TestCase{
		{
			Id:                     uuid.New(),
			InputFileName:          "inputs/test1",
			ExpectedOutputFileName: "expected/ref1",
		},
		{
			Id:                     uuid.New(),
			InputFileName:          "inputs/test2",
			ExpectedOutputFileName: "expected/ref2",
		},
		{
			Id:                     uuid.New(),
			InputFileName:          "inputs/test3",
			ExpectedOutputFileName: "expected/ref3",
		},
	},
}

//Change it to take the problem from the db, instead of hardcoding it
func (submissionWrapper *SubmissionWrapper) RunSubmission(submissionId uuid.UUID) error {

	submissionRunner, mapError := submissionWrapper.LanguageRunners[submission.ProgrammingLanguage]

	if mapError != true {
		return fmt.Errorf("%s is not supported as a programming language", submission.ProgrammingLanguage)
	}

	s3Submission, err := submissionWrapper.S3Repo.GetSubmission(submissionId.String())

	if err != nil {
		return errors.Wrapf(err, "Error trying to download submission: %s from s3", submissionId)
	}

	solutionReq := &dtos.SolutionRequest{
		File:       s3Submission,
		Submission: submission,
		Tests:      problem.TestCases,
	}

	submission.TestResults, err = submissionRunner.RunSubmission(solutionReq)
	if err != nil {
		return err
	}

	fmt.Println(submission)

	return nil

}

func main() {

}
