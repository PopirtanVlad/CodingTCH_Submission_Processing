package executions

import (
	"Licenta_Processing_Service/dtos"
	"Licenta_Processing_Service/repositories"
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
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

func (javaSubmissionRunner *JavaSubmissionRunner) RunSubmission(solutionReq dtos.SolutionRequest) ([]*dtos.TestResult, error) {
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

	_, err = javaSubmissionRunner.compileSolution("Solution.java")
	if err != nil {
		return nil, err
	}

	var results []*dtos.TestResult
	for _, _ = range solutionReq.Tests {

		result, err := runTest()
		if err != nil {
			logrus.WithError(err).Errorf("test execution failed")
			return nil, err
		}

		results = append(results, result)

	}

	return results, nil
}

func runTest() (*dtos.TestResult, error) {
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

func (javaSubmissionRunner *JavaSubmissionRunner) compileSolution(fileName string) (*dtos.SolutionResult, error) {
	cmdConfig := dtos.CommandConfig{
		CommandName: "javac",
		CommandArgs: []string{fileName},
		TimeOut:     10000,
		StdIn:       nil,
		StdOut:      ioutil.Discard,
	}

	return NewExecutionRunner().RunCommand(cmdConfig)
}
