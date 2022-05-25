package language_runners

import (
	"Licenta_Processing_Service/dtos"
	"Licenta_Processing_Service/repositories"
	"Licenta_Processing_Service/services/executions"
	"fmt"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/udhos/equalfile"
	"io"
	"os"
)

var PY_FILE_NAME = "Solution.py"

type PythonSubmissionRunner struct {
	ExecutionRunner executions.ExecutionRunner
	FilesRepository *repositories.FilesRepository
}

func NewPythonSubmissionRunner(repository *repositories.FilesRepository) *PythonSubmissionRunner {
	return &PythonSubmissionRunner{
		ExecutionRunner: *executions.NewExecutionRunner(),
		FilesRepository: repository,
	}
}

func (PythonSubmissionRunner *PythonSubmissionRunner) RunSubmission(solutionReq *dtos.SolutionRequest) ([]*dtos.TestResult, error) {
	/* Salveaza fisierul primit ca parametru, care e luat din s3 si apoi da-i defer sa il stergi. Pe fisierul asta o sa rulez*/
	err := PythonSubmissionRunner.FilesRepository.SaveFile(solutionReq.Submission.ProblemID.String(), solutionReq.Submission.Id.String()+".py", solutionReq.File)

	if err != nil {
		return nil, err
	}

	defer func() {
		err := PythonSubmissionRunner.FilesRepository.DeleteFile(solutionReq.Submission.ProblemID.String(), "Solution.py")
		if err != nil {
			logrus.WithError(err).Warnf("Could not delete file")
		}
	}()

	var results []*dtos.TestResult
	for _, test := range solutionReq.Tests {

		result, err := PythonSubmissionRunner.RunTest(&dtos.RunTestRequest{
			Submission:     solutionReq.Submission,
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

func (PythonSubmissionRunner *PythonSubmissionRunner) RunTest(request *dtos.RunTestRequest) (*dtos.TestResult, error) {
	inputFile, err := PythonSubmissionRunner.FilesRepository.OpenFile(request.Submission.ProblemID.String(), request.Test.InputFileName)
	if err != nil {
		return nil, err
	}

	defer inputFile.Close()

	outputFile, err := PythonSubmissionRunner.FilesRepository.CreateFile(request.Submission.ProblemID.String(), request.OutputFileName)

	if err != nil {
		return nil, err
	}
	defer outputFile.Close()

	testRunDetails, err := PythonSubmissionRunner.executeProgram(request.Submission, inputFile, outputFile)

	defer func() {
		if err := PythonSubmissionRunner.FilesRepository.DeleteFile(request.Submission.ProblemID.String(), request.OutputFileName); err != nil {
			logrus.WithError(err).Errorf("Could not delete output file: %s", request.OutputFileName)
		}

	}()
	if err != nil {
		return nil, err
	}

	areTheSame, err := PythonSubmissionRunner.compareOutput(request.Submission.ProblemID.String(), request.Test.ExpectedOutputFileName, request.OutputFileName)

	if err != nil {
		return nil, err
	}

	return &dtos.TestResult{
		Id:           uuid.New(),
		Correct:      areTheSame,
		TimeElapsed:  testRunDetails.ExecutionTime,
		MemoryUsed:   testRunDetails.MemoryUsage,
		ErrorMessage: "nil",
	}, nil
}

func (PythonSubmissionRunner *PythonSubmissionRunner) executeProgram(submission dtos.Submission, stDin io.ReadCloser, stdOut io.WriteCloser) (*dtos.SolutionResult, error) {

	defer stdOut.Close()
	defer stDin.Close()

	parentDirectory, err := os.Getwd()

	if err != nil {
		return nil, err
	}

	if err := os.Chdir(PythonSubmissionRunner.FilesRepository.GetDirPath(submission.ProblemID.String())); err != nil {
		return nil, err
	}

	defer func() {
		if err = os.Chdir(parentDirectory); err != nil {
			logrus.WithError(err).Errorf("Could not go back to parent directory %s", parentDirectory)
		}
	}()

	cmdConfig := dtos.CommandConfig{
		CommandName: "py",
		CommandArgs: []string{PY_FILE_NAME},
		TimeOut:     2,
		StdIn:       stDin,
		StdOut:      stdOut,
	}
	return PythonSubmissionRunner.ExecutionRunner.RunCommand(cmdConfig)
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
