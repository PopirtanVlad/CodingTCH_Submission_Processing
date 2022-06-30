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
	"os"
)

var PyFileName = "Solution.py"

type PythonSubmissionRunner struct {
	ExecutionRunner executions.ExecutionRunner
	FilesRepository *repositories.FilesRepository
}

func NewPythonSubmissionRunner(repository *repositories.FilesRepository) *PythonSubmissionRunner {
	return &PythonSubmissionRunner{
		ExecutionRunner: *executions.NewExecutionRunner(50),
		FilesRepository: repository,
	}
}

func (PythonSubmissionRunner *PythonSubmissionRunner) RunSubmission(solutionReq *entities.SolutionRequest) ([]*entities.TestResult, error) {
	err := PythonSubmissionRunner.FilesRepository.SaveFile(solutionReq.Problem.ProblemTitle, PyFileName, solutionReq.File)

	if err != nil {
		return nil, err
	}

	defer func() {
		err := PythonSubmissionRunner.FilesRepository.DeleteFile(solutionReq.Problem.ProblemTitle, PyFileName)
		if err != nil {
			logrus.WithError(err).Warnf("Could not delete file")
		}
	}()

	var results []*entities.TestResult
	for _, test := range solutionReq.Tests {

		result, err := PythonSubmissionRunner.RunTest(&entities.RunTestRequest{
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

func (PythonSubmissionRunner *PythonSubmissionRunner) RunTest(request *entities.RunTestRequest) (*entities.TestResult, error) {
	inputFile, err := PythonSubmissionRunner.FilesRepository.OpenFile(request.Problem.ProblemTitle, request.Test.InputFileName)
	if err != nil {
		return nil, err
	}

	defer inputFile.Close()

	outputFile, err := PythonSubmissionRunner.FilesRepository.CreateFile(request.Problem.ProblemTitle, request.OutputFileName)

	if err != nil {
		return nil, err
	}
	defer outputFile.Close()

	testRunDetails, err := PythonSubmissionRunner.executeProgram(request.Problem, inputFile, outputFile)

	defer func() {
		if err := PythonSubmissionRunner.FilesRepository.DeleteFile(request.Problem.ProblemTitle, request.OutputFileName); err != nil {
			logrus.WithError(err).Errorf("Could not delete output file: %s", request.OutputFileName)
		}

	}()
	if err != nil {
		return nil, err
	}

	areTheSame, err := PythonSubmissionRunner.compareOutput(request.Problem.ProblemTitle, request.Test.ExpectedOutputFileName, request.OutputFileName)

	if err != nil {
		return nil, err
	}

	return &entities.TestResult{
		Id:           uuid.New().String(),
		Correct:      areTheSame,
		TimeElapsed:  testRunDetails.ExecutionTime,
		MemoryUsed:   testRunDetails.MemoryUsage,
		ErrorMessage: testRunDetails.StdErr,
		SubmissionId: request.Submission.Id,
	}, nil
}

func (PythonSubmissionRunner *PythonSubmissionRunner) executeProgram(problem entities.Problem, stDin io.ReadCloser, stdOut io.WriteCloser) (*entities.SolutionResult, error) {

	defer stdOut.Close()
	defer stDin.Close()

	parentDirectory, err := os.Getwd()

	if err != nil {
		return nil, err
	}

	if err := os.Chdir(PythonSubmissionRunner.FilesRepository.GetDirPath(problem.ProblemTitle)); err != nil {
		return nil, err
	}

	defer func() {
		if err = os.Chdir(parentDirectory); err != nil {
			logrus.WithError(err).Errorf("Could not go back to parent directory %s", parentDirectory)
		}
	}()

	cmdConfig := entities.CommandConfig{
		CommandName: "py",
		CommandArgs: []string{PyFileName},
		TimeOut:     2,
		StdIn:       stDin,
		StdOut:      stdOut,
	}
	return PythonSubmissionRunner.ExecutionRunner.RunCommand(cmdConfig, 1, 1), nil
}

func (PythonSubmissionRunner *PythonSubmissionRunner) compareOutput(pathDir, outPutFileName, refFileName string) (bool, error) {
	outputPath, _ := PythonSubmissionRunner.FilesRepository.OpenFile(pathDir, outPutFileName)
	refPath, _ := PythonSubmissionRunner.FilesRepository.OpenFile(pathDir, refFileName)

	defer outputPath.Close()
	defer refPath.Close()

	equal, err := equalfile.New(nil, equalfile.Options{}).CompareReader(outputPath, refPath)
	if err != nil {
		return false, err
	}

	return equal, nil
}
