package executions

import (
	"Licenta_Processing_Service/entities"
	"bytes"
	"context"
	"fmt"
	_ "fmt"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	_ "os"
	"os/exec"
	"time"
)

type ExecutionRunner struct {
	memoryMonitor *MemoryMonitor
}

func NewExecutionRunner(monitoringInterval time.Duration) *ExecutionRunner {
	monitor := NewMemoryMonitor(monitoringInterval)
	return &ExecutionRunner{
		memoryMonitor: monitor,
	}
}

func (executionRunner *ExecutionRunner) RunCommand(cmdConfig entities.CommandConfig, timeLimit time.Duration, memoryLimit uint64) *entities.SolutionResult {
	ctx, cancel := context.WithTimeout(context.TODO(), timeLimit)
	defer cancel()
	cmd := exec.Command(cmdConfig.CommandName, cmdConfig.CommandArgs...)
	startTime := time.Now()
	var errBuff bytes.Buffer
	var endTime time.Duration
	var maxRecorderMemory uint64
	cmd.Stderr = &errBuff
	cmd.Stdout = cmdConfig.StdOut
	cmd.Stdin = cmdConfig.StdIn

	if err := cmd.Start(); err != nil {
		fmt.Println(err)
		return &entities.SolutionResult{
			StdErr:        "Execution error",
			ExecutionTime: endTime,
			MemoryUsage:   maxRecorderMemory,
			ExitCode:      -1,
		}
	}
	monitorCtx, monitorCancel := context.WithCancel(ctx)
	defer monitorCancel()
	memoryChanRes, memoryChanErr := executionRunner.memoryMonitor.StartMonitor(monitorCtx, cmd.Process.Pid, memoryLimit)
	done := make(chan error, 1)
	go func() {
		defer close(done)
		done <- cmd.Wait()
		endTime = time.Since(startTime)
	}()
	for {
		select {
		case result, ok := <-memoryChanRes:
			if ok {
				if maxRecorderMemory < result {
					maxRecorderMemory = result
				}
			}
		case <-ctx.Done():
			monitorCancel()
			_ = executionRunner.killProcess(cmd)
			return executionRunner.returnTLE(time.Since(startTime), maxRecorderMemory)
		case _, ok := <-memoryChanErr:
			if ok {
				monitorCancel()
				_ = executionRunner.killProcess(cmd)
				return executionRunner.returnMLE(time.Since(startTime), maxRecorderMemory)
			}
		case err, ok := <-done:
			if ok {
				monitorCancel()
				if err == nil && endTime > timeLimit {
					return executionRunner.returnTLE(endTime, maxRecorderMemory)
				}
				if err == nil && maxRecorderMemory > memoryLimit {
					return executionRunner.returnMLE(endTime, maxRecorderMemory)
				}

				var exitCode int
				if err != nil {
					if exitErr, ok := err.(*exec.ExitError); ok {
						exitCode = exitErr.ExitCode()
					}
					logrus.
						WithError(err).
						WithField("Exit Code", exitCode).
						WithField("Cmd", cmd.Args).
						Debugf("Finished command execution")
				}
				return &entities.SolutionResult{
					StdErr:        errBuff.String(),
					ExecutionTime: endTime,
					MemoryUsage:   maxRecorderMemory,
					ExitCode:      exitCode,
				}
			}
			return nil
		}
	}

}

func (executionRunner *ExecutionRunner) killProcess(cmd *exec.Cmd) error {
	if err := cmd.Process.Kill(); err != nil {
		logrus.WithError(err).Errorf("could not kill process %d started for command %s", cmd.Process.Pid, cmd.Args)
		return errors.Wrapf(err, "could not kill process %d started for command %s", cmd.Process.Pid, cmd.Args)
	}
	return nil
}

func (executionRunner *ExecutionRunner) returnTLE(timeElapsed time.Duration, maxRecoredMemory uint64) *entities.SolutionResult {
	return &entities.SolutionResult{
		StdErr:        "Time limit exceeded",
		ExecutionTime: timeElapsed,
		MemoryUsage:   maxRecoredMemory,
		ExitCode:      0,
	}
}

func (executionRunner *ExecutionRunner) returnMLE(timeElapsed time.Duration, maxRecoredMemory uint64) *entities.SolutionResult {
	return &entities.SolutionResult{
		StdErr:        "Memory limit exceeded",
		ExecutionTime: timeElapsed,
		MemoryUsage:   maxRecoredMemory,
		ExitCode:      0,
	}
}
