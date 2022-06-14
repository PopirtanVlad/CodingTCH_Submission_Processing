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

func (submissionWrapper *SubmissionWrapper) RunSubmission(submissionId uuid.UUID) error {
	//We take the submission from the database
	submission, err := submissionWrapper.DbRepo.GetSubmission(submissionId.String())
	if err != nil {
		return errors.Wrapf(err, "Could not get the submission: %s from the database", submissionId)
	}
	//We get the problem from the database
	problem, err := submissionWrapper.DbRepo.GetProblem(submission.ProblemID.String())
	if err != nil {
		return errors.Wrapf(err, "Could not get the problem: %s from the database", submission.ProblemID)
	}

	//We take the submission from the aws s3 repository
	s3Submission, err := submissionWrapper.S3Repo.GetSubmission(problem.ProblemTitle, submissionId.String())
	if err != nil {
		return errors.Wrapf(err, "Error trying to download submission: %s from s3", submissionId)
	}

	if err = submissionWrapper.S3Repo.DownloadTests(problem.ProblemTitle); err != nil {
		return errors.Wrapf(err, "Could sync the files for the problem: %s", problem.Id.String())
	}

	//We set the correct code runner, according to the programming language that the solution was written in
	submissionRunner, mapError := submissionWrapper.LanguageRunners[submission.ProgrammingLanguage]
	if mapError != true {
		return fmt.Errorf("%s is not supported as a programming language", submission.ProgrammingLanguage)
	}

	tests, err := submissionWrapper.DbRepo.GetTests(problem.ProblemTitle)

	if err != nil {
		return fmt.Errorf("couldn't get tests cases for problem: %s", problem.ProblemTitle)
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
