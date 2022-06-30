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
)

const (
	CFileName         = "Solution.c"
	CompiledCFileName = "Solution.exe"
)

type CSubmissionRunner struct {
	ExecutionRunner executions.ExecutionRunner
	FilesRepository *repositories.FilesRepository
}

func NewCSubmissionRunner(repository *repositories.FilesRepository) *CSubmissionRunner {
	return &CSubmissionRunner{
		ExecutionRunner: *executions.NewExecutionRunner(50),
		FilesRepository: repository,
	}
}

func (CSubmissionRunner *CSubmissionRunner) RunSubmission(solutionReq *entities.SolutionRequest) ([]*entities.TestResult, error) {
	/* Salveaza fisierul primit ca parametru, care e luat din s3 si apoi da-i defer sa il stergi. Pe fisierul asta o sa rulez*/
	err := CSubmissionRunner.FilesRepository.SaveFile(solutionReq.Problem.ProblemTitle, CFileName, solutionReq.File)

	if err != nil {
		return nil, err
	}

	defer func() {
		err := CSubmissionRunner.FilesRepository.DeleteFile(solutionReq.Problem.ProblemTitle, CFileName)
		if err != nil {
			logrus.WithError(err).Warnf("Could not delete file")
		}
	}()

	_, err = CSubmissionRunner.compileSolution(solutionReq)
	if err != nil {
		panic(err)
		return nil, err
	}

	defer func() {
		err := CSubmissionRunner.FilesRepository.DeleteFile(solutionReq.Problem.ProblemTitle, CompiledCFileName)
		if err != nil {
			logrus.WithError(err).Warnf("Could not delete file")
		}
	}()

	var results []*entities.TestResult
	for _, test := range solutionReq.Tests {

		result, err := CSubmissionRunner.RunTest(&entities.RunTestRequest{
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

func (CSubmissionRunner *CSubmissionRunner) RunTest(request *entities.RunTestRequest) (*entities.TestResult, error) {
	inputFile, err := CSubmissionRunner.FilesRepository.OpenFile(request.Problem.ProblemTitle, request.Test.InputFileName)
	if err != nil {
		return nil, err
	}

	defer inputFile.Close()

	outputFile, err := CSubmissionRunner.FilesRepository.CreateFile(request.Problem.ProblemTitle, request.OutputFileName)

	if err != nil {
		return nil, err
	}
	defer outputFile.Close()

	testRunDetails, err := CSubmissionRunner.executeProgram(request.Problem, inputFile, outputFile)

	defer func() {
		if err := CSubmissionRunner.FilesRepository.DeleteFile(request.Problem.ProblemTitle, request.OutputFileName); err != nil {
			logrus.WithError(err).Errorf("Could not delete output file: %s", request.OutputFileName)
		}

	}()
	if err != nil {
		return nil, err
	}

	areTheSame, err := CSubmissionRunner.compareOutput(request.Problem.ProblemTitle, request.Test.ExpectedOutputFileName, request.OutputFileName)

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

func (CSubmissionRunner *CSubmissionRunner) executeProgram(problem entities.Problem, stDin io.ReadCloser, stdOut io.WriteCloser) (*entities.SolutionResult, error) {

	defer stdOut.Close()
	defer stDin.Close()

	parentDirectory, err := os.Getwd()

	if err != nil {
		return nil, err
	}

	if err := os.Chdir(CSubmissionRunner.FilesRepository.GetDirPath(problem.ProblemTitle)); err != nil {
		return nil, err
	}

	defer func() {
		if err = os.Chdir(parentDirectory); err != nil {
			logrus.WithError(err).Errorf("Could not go back to parent directory %s", parentDirectory)
		}
	}()

	cmdConfig := entities.CommandConfig{
		CommandName: "./" + CompiledCFileName,
		CommandArgs: []string{},
		TimeOut:     2,
		StdIn:       stDin,
		StdOut:      stdOut,
	}
	return CSubmissionRunner.ExecutionRunner.RunCommand(cmdConfig, 1, 1), nil
}

func (CSubmissionRunner *CSubmissionRunner) compileSolution(request *entities.SolutionRequest) (*entities.SolutionResult, error) {
	solutionPath := CSubmissionRunner.FilesRepository.GetFilePath(request.Problem.ProblemTitle, CFileName)

	cmdConfig := entities.CommandConfig{
		CommandName: "gcc",
		CommandArgs: []string{solutionPath, "-o", CompiledCFileName},
		TimeOut:     1400000,
		StdIn:       nil,
		StdOut:      ioutil.Discard,
	}

	return executions.NewExecutionRunner(50).RunCommand(cmdConfig, 1, 1), nil
}

func (CSubmissionRunner *CSubmissionRunner) compareOutput(pathDir, outPutFileName, refFileName string) (bool, error) {
	outputPath, _ := CSubmissionRunner.FilesRepository.OpenFile(pathDir, outPutFileName)
	refPath, _ := CSubmissionRunner.FilesRepository.OpenFile(pathDir, refFileName)

	defer outputPath.Close()
	defer refPath.Close()

	equal, err := equalfile.New(nil, equalfile.Options{}).CompareReader(outputPath, refPath)
	if err != nil {
		return false, err
	}

	return equal, nil
}
