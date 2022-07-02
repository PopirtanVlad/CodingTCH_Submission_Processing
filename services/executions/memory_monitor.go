package executions

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

// MemoryMonitor monitors memory usage of a process
type MemoryMonitor struct {
	monitorInterval time.Duration
}

// NewMemoryMonitor returns a new MemoryMonitor object
func NewMemoryMonitor(interval time.Duration) *MemoryMonitor {
	return &MemoryMonitor{
		monitorInterval: interval,
	}
}

// StartMonitor starts monitoring memory usage of a process by pid and sends
// an error through a channel in case the maxMemory threshold is exceeded
func (m *MemoryMonitor) StartMonitor(ctx context.Context, pid int, maxMemory uint64) (<-chan uint64, <-chan error) {
	memoryChanRes := make(chan uint64, 1)
	memoryChanErr := make(chan error, 1)
	memoryResults := m.checkMemoryAtInterval(ctx, pid)
	go func() {
		defer close(memoryChanErr)
		errSent := false
		for {
			select {
			case <-ctx.Done():
				return
			case result, ok := <-memoryResults:
				if ok {
					memoryChanRes <- result
					if result > maxMemory {
						if len(memoryChanErr) < cap(memoryChanErr) && !errSent {
							memoryChanErr <- errors.New("memory limit exceeded")
							errSent = true
						}
					}
				}
			}
		}
	}()
	return memoryChanRes, memoryChanErr
}

func (m *MemoryMonitor) checkMemoryAtInterval(ctx context.Context, pid int) <-chan uint64 {
	resultsChan := make(chan uint64, 1)
	ticker := time.NewTicker(m.monitorInterval)

	usedMemory, err := m.getMemoryForProcess(pid)
	logrus.Debugf("checking memory usage for pid %d %d KB", pid, usedMemory)
	if err != nil {
		logrus.WithError(err).
			Debugf("could not check memory for pid %d", pid)
	}
	if len(resultsChan) < cap(resultsChan) {
		resultsChan <- usedMemory
	}

	go func() {
		defer close(resultsChan)
		for {
			select {
			case <-ctx.Done():
				logrus.Debugf("memory monitor stopped for pid %d", pid)
				ticker.Stop()
				return
			case <-ticker.C:
				usedMemory, err := m.getMemoryForProcess(pid)
				logrus.Debugf("checking memory usage for pid %d %d KB", pid, usedMemory)
				if err != nil {
					logrus.WithError(err).
						Debugf("could not check memory for pid %d", pid)
				}
				if len(resultsChan) < cap(resultsChan) {
					resultsChan <- usedMemory
				}
			}
		}
	}()
	return resultsChan
}

func (m *MemoryMonitor) getMemoryForProcess(pid int) (uint64, error) {
	file, err := os.Open(fmt.Sprintf("/proc/%d/smaps", pid))
	if err != nil {
		return 0, err
	}
	defer func() {
		err := file.Close()
		if err != nil {
			logrus.WithError(err).
				Errorf("could not close smaps file for process %d", pid)
		}
	}()

	totalMemory := uint64(0)
	prefix := []byte("Pss:")
	fileScanner := bufio.NewScanner(file)
	for fileScanner.Scan() {
		line := fileScanner.Bytes()
		if bytes.HasPrefix(line, prefix) {
			var size uint64
			_, err := fmt.Sscanf(string(line[4:]), "%d", &size)
			if err != nil {
				return 0, err
			}
			totalMemory += size
		}
	}
	if err := fileScanner.Err(); err != nil {
		return 0, err
	}

	return totalMemory, nil
}
