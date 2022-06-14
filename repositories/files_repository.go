package repositories

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"path/filepath"
)

type FilesRepository struct {
	BaseDirectory string
}

func NewFileRepository(baseDirectory string) (*FilesRepository, error) {
	if err := os.RemoveAll(baseDirectory); err != nil {
		return nil, err
	}

	if err := os.Mkdir(baseDirectory, os.ModePerm); err != nil {
		return nil, err
	}
	return &FilesRepository{
		BaseDirectory: baseDirectory,
	}, nil
}

func (fileRepository *FilesRepository) GetDirPath(problemDir string) string {
	return fmt.Sprintf("%s/%s", fileRepository.BaseDirectory, problemDir)
}

func (fileRepository *FilesRepository) GetFilePath(problemdir, filename string) string {
	switch filepath.Ext(filename) {
	case ".in":
		return fmt.Sprintf("%s/problems/%s/inputs/%s\n", fileRepository.BaseDirectory, problemdir, filename)
	case ".ref":
		return fmt.Sprintf("%s/problems/%s/expected/%s\n", fileRepository.BaseDirectory, problemdir, filename)
	default:
		return fmt.Sprintf("%s/submissions/%s/%s", fileRepository.BaseDirectory, problemdir, filename)
	}
}

func (fileRepository *FilesRepository) OpenFile(problemDir, fileName string) (*os.File, error) {
	filePath := fileRepository.GetFilePath(problemDir, fileName)
	return os.Open(filePath)
}

func (fileRepository *FilesRepository) SaveFile(problemDir, fileName string, sourceFile io.Reader) error {
	dirPath := fmt.Sprintf("%s/%s", fileRepository.BaseDirectory, problemDir)

	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		logrus.Infof("The directory %s doesn't exist. Trying to create it", dirPath)
		if err = os.Mkdir(dirPath, os.ModePerm); err != nil {
			return errors.Wrapf(err, "directory %s couldn't be created", dirPath)
		}
	}

	filePath := fmt.Sprintf("%s/%s/%s", fileRepository.BaseDirectory, problemDir, fileName)
	destFile, err := os.Create(filePath)
	if err != nil {
		return err
	}
	_, err = io.Copy(destFile, sourceFile)
	defer func() { destFile.Close() }()
	if err != nil {
		return err
	}
	return nil
	//if err != nil {
	//	logrus.WithFields(logrus.Fields{
	//		"file path": filePath,
	//	}).WithError(err).Debug("error trying to create the file")
	//}
}

func (fileRepository *FilesRepository) DeleteFile(problemDir, fileName string) error {
	filepath := fmt.Sprintf("%s/%s/%s", fileRepository.BaseDirectory, problemDir, fileName)
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return errors.Wrapf(err, "The file %s does not exist", filepath)
	}
	return os.Remove(filepath)
	//if err != nil {
	//	logrus.WithFields(logrus.Fields{
	//		"file path": filepath,
	//	}).WithError(err).Debug("error trying to delete the file")
	//}
}

func (fileRepository *FilesRepository) CreateFile(problemDir, fileName string) (io.WriteCloser, error) {
	filepath := fmt.Sprintf("%s/%s/%s", fileRepository.BaseDirectory, problemDir, fileName)
	return os.Create(filepath)
}
