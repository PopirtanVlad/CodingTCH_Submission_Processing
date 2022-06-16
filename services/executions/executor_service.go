package executions

import (
	"Licenta_Processing_Service/dtos"
	"bytes"
	_ "fmt"
	"github.com/sirupsen/logrus"
	_ "os"
	"os/exec"
	"time"
)

type ExecutionRunner struct {
}

func NewExecutionRunner() *ExecutionRunner {
	return &ExecutionRunner{}
}

func (executionRunner *ExecutionRunner) RunCommand(cmdConfig dtos.CommandConfig, timeLimit time.Duration, memoryLimit uint64) *dtos.SolutionResult {
	cmd := exec.Command(cmdConfig.CommandName, cmdConfig.CommandArgs...)
	startTime := time.Now()
	var errBuff bytes.Buffer

	cmd.Stderr = &errBuff
	cmd.Stdout = cmdConfig.StdOut
	cmd.Stdin = cmdConfig.StdIn
	if err := cmd.Run(); err != nil {
		return &dtos.SolutionResult{
			ExecutionTime: 0,
			MemoryUsage:   0,
			StdErr:        err.Error(),
			ExitCode:      0,
		} //errors.Wrapf(err, "Could not run the command: %s \nError: %s", cmd.String(), errBuff.String())
	}

	err := cmd.Wait()
	endTime := time.Since(startTime)
	var exitCode int

	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		}
	}

	if endTime > timeLimit {
		logrus.WithFields(logrus.Fields{
			"Time limit":  timeLimit,
			"Actual time": endTime,
		}).Debugf("Time limit exceeded")
		return executionRunner.returnTLE(endTime, 0)
	}

	return &dtos.SolutionResult{
		StdErr:        "NaN",
		MemoryUsage:   0,
		ExecutionTime: endTime,
		ExitCode:      exitCode,
	}
}

func (e *ExecutionRunner) returnTLE(timeElapsed time.Duration, maxMemory uint64) *dtos.SolutionResult {
	return &dtos.SolutionResult{
		StdErr:        "Time limit exceeded",
		ExecutionTime: timeElapsed,
		MemoryUsage:   maxMemory,
		ExitCode:      0,
	}
}
