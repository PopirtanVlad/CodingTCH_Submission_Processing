package language_runners

import (
	"Licenta_Processing_Service/entities"
	"Licenta_Processing_Service/repositories"
	"Licenta_Processing_Service/services/executions"
	"fmt"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/udhos/equalfile"
	"io"
	"io/ioutil"
	"os"
	"time"
)

const (
	JavaFileName         = "Solution.java"
	JavaClassName        = "Solution.class"
	CompiledJavaFileName = "Solution"
)

type JavaSubmissionRunner struct {
	ExecutionRunner executions.ExecutionRunner
	FilesRepository *repositories.FilesRepository
}

func NewJavaSubmissionRunner(repository *repositories.FilesRepository) *JavaSubmissionRunner {
	return &JavaSubmissionRunner{
		ExecutionRunner: *executions.NewExecutionRunner(50),
		FilesRepository: repository,
	}
}

func (javaSubmissionRunner *JavaSubmissionRunner) RunSubmission(solutionReq *entities.SolutionRequest) ([]*entities.TestResult, error) {
	err := javaSubmissionRunner.FilesRepository.SaveFile(solutionReq.Problem.ProblemTitle, JavaFileName, solutionReq.File)

	if err != nil {
		return nil, err
	}

	defer func() {
		err := javaSubmissionRunner.FilesRepository.DeleteFile(solutionReq.Problem.ProblemTitle, JavaFileName)
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
		err := javaSubmissionRunner.FilesRepository.DeleteFile(solutionReq.Problem.ProblemTitle, JavaClassName)
		if err != nil {
			logrus.WithError(err).Warnf("Could not delete file")
		}
	}()
	var results []*entities.TestResult
	for _, test := range solutionReq.Tests {
		result, err := javaSubmissionRunner.RunTest(&entities.RunTestRequest{
			Submission:     solutionReq.Submission,
			Problem:        solutionReq.Problem,
			Test:           test,
			OutputFileName: uuid.New().String(),
		})
		if err != nil {
			logrus.WithError(err).Errorf("test execution failed")
			return nil, err
		}

		results = append(results, result)

	}
	for _, result := range results {
		fmt.Println(result)
	}
	return results, nil
}

func (javaSubmissionRunner *JavaSubmissionRunner) RunTest(request *entities.RunTestRequest) (*entities.TestResult, error) {
	inputFile, err := javaSubmissionRunner.FilesRepository.OpenFile(request.Problem.ProblemTitle, request.Test.InputFileName)
	if err != nil {
		return nil, err
	}

	defer inputFile.Close()

	outputFile, err := javaSubmissionRunner.FilesRepository.CreateFile(request.Problem.ProblemTitle, request.OutputFileName)

	if err != nil {
		return nil, err
	}
	defer outputFile.Close()

	testRunDetails, err := javaSubmissionRunner.executeProgram(request.Problem, inputFile, outputFile)

	defer func() {
		if err := javaSubmissionRunner.FilesRepository.DeleteFile(request.Problem.ProblemTitle, request.OutputFileName); err != nil {
			logrus.WithError(err).Errorf("Could not delete output file: %s", request.OutputFileName)
		}

	}()
	if err != nil {
		return nil, err
	}
	areTheSame, err := javaSubmissionRunner.compareOutput(request.Problem.ProblemTitle, request.Test.ExpectedOutputFileName, request.OutputFileName)

	if err != nil {
		return nil, err
	}

	testResult := &entities.TestResult{
		Id:           uuid.New().String(),
		Correct:      areTheSame,
		TimeElapsed:  testRunDetails.ExecutionTime,
		MemoryUsed:   testRunDetails.MemoryUsage,
		ErrorMessage: testRunDetails.StdErr,
		SubmissionId: request.Submission.Id,
	}
	if testResult.ErrorMessage != "" {
		testResult.Correct = false
	}
	if testResult.ErrorMessage == "" && !areTheSame {
		testResult.ErrorMessage = "Wrong Answer"
	}
	return testResult, nil
}

func (javaSubmissionRunner *JavaSubmissionRunner) executeProgram(problem entities.Problem, stDin io.ReadCloser, stdOut io.WriteCloser) (*entities.SolutionResult, error) {

	defer stdOut.Close()
	defer stDin.Close()

	parentDirectory, err := os.Getwd()

	if err != nil {
		return nil, err
	}

	if err := os.Chdir(javaSubmissionRunner.FilesRepository.GetDirPath(problem.ProblemTitle)); err != nil {
		return nil, err
	}
	defer func() {
		if err = os.Chdir(parentDirectory); err != nil {
			logrus.WithError(err).Errorf("Could not go back to parent directory %s", parentDirectory)
		}
	}()
	cmdConfig := entities.CommandConfig{
		CommandName: "java",
		CommandArgs: []string{CompiledJavaFileName},
		TimeOut:     2,
		StdIn:       stDin,
		StdOut:      stdOut,
	}
	return javaSubmissionRunner.ExecutionRunner.RunCommand(cmdConfig, problem.TimeLimit, problem.MemoryLimit), nil
}

func (javaSubmissionRunner *JavaSubmissionRunner) compileSolution(request *entities.SolutionRequest) (*entities.SolutionResult, error) {
	solutionPath := javaSubmissionRunner.FilesRepository.GetFilePath(request.Problem.ProblemTitle, JavaFileName)

	cmdConfig := entities.CommandConfig{
		CommandName: "javac",
		CommandArgs: []string{solutionPath},
		TimeOut:     30000,
		StdIn:       nil,
		StdOut:      ioutil.Discard,
	}

	return executions.NewExecutionRunner(50).RunCommand(cmdConfig, time.Second, 100000), nil
}

func (javaSubmissionRunner *JavaSubmissionRunner) compareOutput(pathDir, outPutFileName, refFileName string) (bool, error) {
	outputPath, _ := javaSubmissionRunner.FilesRepository.OpenFile(pathDir, outPutFileName)
	refPath, _ := javaSubmissionRunner.FilesRepository.OpenFile(pathDir, refFileName)

	defer outputPath.Close()
	defer refPath.Close()

	equal, err := equalfile.New(nil, equalfile.Options{}).CompareReader(outputPath, refPath)
	if err != nil {
		return false, err
	}

	return equal, nil

}
