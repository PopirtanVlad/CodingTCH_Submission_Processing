package executions

import (
	"Licenta_Processing_Service/dtos"
	"Licenta_Processing_Service/repositories"
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"time"
)

var FILE_ID = "java_test"
var FILE_NAME = "Solution.java"

type JavaSubmissionRunner struct {
	ExecutionRunner ExecutionRunner
	FilesRepository *repositories.FilesRepository
}

func NewJavaSubmissionRunner(repository *repositories.FilesRepository) *JavaSubmissionRunner {
	return &JavaSubmissionRunner{
		ExecutionRunner: *NewExecutionRunner(),
		FilesRepository: repository,
	}
}

func (javaSubmissionRunner *JavaSubmissionRunner) RunSubmission(solutionReq *dtos.SolutionRequest) ([]*dtos.TestResult, error) {
	/* Salveaza fisierul primit ca parametru, care e luat din s3 si apoi da-i defer sa il stergi. Pe fisierul asta o sa rulez*/
	err := javaSubmissionRunner.FilesRepository.SaveFile(solutionReq.Solution.ProblemID.String(), "Solution.java", solutionReq.File)

	if err != nil {
		return nil, err
	}

	defer func() {
		err := javaSubmissionRunner.FilesRepository.DeleteFile(solutionReq.Solution.ProblemID.String(), "Solution.java")
		if err != nil {
			logrus.WithError(err).Warnf("Could not delete file")
		}
	}()

	_, err = javaSubmissionRunner.compileSolution(solutionReq)
	if err != nil {
		panic(err)
		return nil, err
	}

	defer func() {
		err := javaSubmissionRunner.FilesRepository.DeleteFile(solutionReq.Solution.ProblemID.String(), "Solution.class")
		if err != nil {
			logrus.WithError(err).Warnf("Could not delete file")
		}
	}()

	var results []*dtos.TestResult
	for _, _ = range solutionReq.Tests {

		result, err := javaSubmissionRunner.RunTest(solutionReq)
		if err != nil {
			logrus.WithError(err).Errorf("test execution failed")
			return nil, err
		}

		results = append(results, result)

	}
	time.Sleep(time.Second * 5)
	return results, nil
}

func (javaSubmissionRunner *JavaSubmissionRunner) RunTest(request *dtos.SolutionRequest) (*dtos.TestResult, error) {
	return nil, nil
}

func (javaSubmissionRunner *JavaSubmissionRunner) executeProgram(submission dtos.Submission, stDin io.ReadCloser, stdOut io.WriteCloser) (*dtos.SolutionResult, error) {

	defer stdOut.Close()
	defer stDin.Close()

	args := []string{
		javaSubmissionRunner.FilesRepository.GetFilePath(submission.ProblemID.String(), "Solution.java"),
	}

	cmdConfig := dtos.CommandConfig{
		CommandName: "java",
		CommandArgs: args,
		TimeOut:     2,
		StdIn:       stDin,
		StdOut:      stdOut,
	}
	return javaSubmissionRunner.ExecutionRunner.RunCommand(cmdConfig)
}

func (javaSubmissionRunner *JavaSubmissionRunner) compileSolution(request *dtos.SolutionRequest) (*dtos.SolutionResult, error) {
	solutionPath := javaSubmissionRunner.FilesRepository.GetFilePath(request.Solution.ProblemID.String(), "Solution.java")

	cmdConfig := dtos.CommandConfig{
		CommandName: "javac",
		CommandArgs: []string{solutionPath},
		TimeOut:     30000,
		StdIn:       nil,
		StdOut:      ioutil.Discard,
	}

	return NewExecutionRunner().RunCommand(cmdConfig)
}
