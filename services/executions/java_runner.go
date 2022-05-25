package executions

import (
	"Licenta_Processing_Service/dtos"
	"Licenta_Processing_Service/repositories"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/udhos/equalfile"
	"io"
	"io/ioutil"
	"os"
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
	err := javaSubmissionRunner.FilesRepository.SaveFile(solutionReq.Submission.ProblemID.String(), "Solution.java", solutionReq.File)

	if err != nil {
		return nil, err
	}

	defer func() {
		err := javaSubmissionRunner.FilesRepository.DeleteFile(solutionReq.Submission.ProblemID.String(), "Solution.java")
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
		err := javaSubmissionRunner.FilesRepository.DeleteFile(solutionReq.Submission.ProblemID.String(), "Solution.class")
		if err != nil {
			logrus.WithError(err).Warnf("Could not delete file")
		}
	}()

	var results []*dtos.TestResult
	for _, test := range solutionReq.Tests {

		result, err := javaSubmissionRunner.RunTest(&dtos.RunTestRequest{
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
	return results, nil
}

func (javaSubmissionRunner *JavaSubmissionRunner) RunTest(request *dtos.RunTestRequest) (*dtos.TestResult, error) {
	inputFile, err := javaSubmissionRunner.FilesRepository.OpenFile(request.Submission.ProblemID.String(), request.Test.InputFileName)
	if err != nil {
		return nil, err
	}

	defer inputFile.Close()

	outputFile, err := javaSubmissionRunner.FilesRepository.CreateFile(request.Submission.ProblemID.String(), request.OutputFileName)

	if err != nil {
		return nil, err
	}
	defer outputFile.Close()

	_, err = javaSubmissionRunner.executeProgram(request.Submission, inputFile, outputFile)

	defer func() {
		if err := javaSubmissionRunner.FilesRepository.DeleteFile(request.Submission.ProblemID.String(), request.OutputFileName); err != nil {
			logrus.WithError(err).Errorf("Could not delete output file: %s", request.OutputFileName)
		}

	}()
	if err != nil {
		return nil, err
	}

	javaSubmissionRunner.compareOutput(request.Submission.ProblemID.String(), request.Test.ExpectedOutputFileName, request.OutputFileName)

	return nil, err
}

func (javaSubmissionRunner *JavaSubmissionRunner) executeProgram(submission dtos.Submission, stDin io.ReadCloser, stdOut io.WriteCloser) (*dtos.SolutionResult, error) {

	defer stdOut.Close()
	defer stDin.Close()

	parentDirectory, err := os.Getwd()

	if err != nil {
		return nil, err
	}

	if err := os.Chdir(javaSubmissionRunner.FilesRepository.GetDirPath(submission.ProblemID.String())); err != nil {
		return nil, err
	}

	defer func() {
		if err = os.Chdir(parentDirectory); err != nil {
			logrus.WithError(err).Errorf("Could not go back to parent directory %s", parentDirectory)
		}
	}()

	cmdConfig := dtos.CommandConfig{
		CommandName: "java",
		CommandArgs: []string{"Solution"},
		TimeOut:     2,
		StdIn:       stDin,
		StdOut:      stdOut,
	}
	return javaSubmissionRunner.ExecutionRunner.RunCommand(cmdConfig)
}

func (javaSubmissionRunner *JavaSubmissionRunner) compileSolution(request *dtos.SolutionRequest) (*dtos.SolutionResult, error) {
	solutionPath := javaSubmissionRunner.FilesRepository.GetFilePath(request.Submission.ProblemID.String(), "Solution.java")

	cmdConfig := dtos.CommandConfig{
		CommandName: "javac",
		CommandArgs: []string{solutionPath},
		TimeOut:     30000,
		StdIn:       nil,
		StdOut:      ioutil.Discard,
	}

	return NewExecutionRunner().RunCommand(cmdConfig)
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
