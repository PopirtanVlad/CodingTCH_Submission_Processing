package main

import (
	"Licenta_Processing_Service/dtos"
	"Licenta_Processing_Service/services/executions"
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
)

var FILE_ID = "123"
var FILE_NAME = "Solution.c"

//func main() {
//	compileSolution("Solution.c")
//}

type JavaSubmissionRunner struct {
	ExecutionRunner executions.ExecutionRunner
}

func NewJavaSubmissionRunner() *JavaSubmissionRunner {
	return &JavaSubmissionRunner{}
}
=

func (javaSubmissionRunner *JavaSubmissionRunner) RunSolution(solutionReq dtos.SolutionRequest) ([]*dtos.TestResult, error) {
	/* Salveaza fisierul primit ca parametru, care e luat din s3 si apoi da-i defer sa il stergi. Pe fisierul asta o sa rulez*/

	_, err = javaSubmissionRunner.compileSolution("Solution.c")

	if err != nil {
		return nil, err
	}

	var results []*dtos.TestResult
	for _, test := range solutionReq.Tests {

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
	return nil,nil
}

func (javaSubmissionRunner *JavaSubmissionRunner) executeProgram(submission dtos.Submission, stDin io.ReadCloser, stdOut io.WriteCloser) (*dtos.SolutionResult, error) {

	defer stdOut.Close()
	defer stDin.Close()

	//Get problem dir path err = os.Chdir(submission.SubmissionID)

	if err != nil {
		return nil, err
	}

	cmdConfig := dtos.CommandConfig{
		CommandName: "java",
		CommandArgs: []string{FILE_NAME},
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

	return executions.NewExecutionRunner().RunCommand(cmdConfig)
}
