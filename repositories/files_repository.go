package repositories

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type FilesRepository struct {
	BaseDirectory string
}

func NewFileRepository(baseDirectory string) (*FilesRepository, error) {
	//if err := os.RemoveAll(baseDirectory); err != nil {
	//	return nil, err
	//}

	//if err := os.Mkdir(baseDirectory, os.ModePerm); err != nil {
	//	return nil, err
	//}
	return &FilesRepository{
		BaseDirectory: baseDirectory,
	}, nil
}

func (fileRepository *FilesRepository) GetFilePath(problemdir, filename string) string {
	switch filepath.Ext(filename) {
	case ".in":
		return fmt.Sprintf("%s/%s/inputs/%s\n", fileRepository.BaseDirectory, problemdir, filename)
	case ".ref":
		return fmt.Sprintf("%s/%s/expected/%s\n", fileRepository.BaseDirectory, problemdir, filename)
	default:
		return fmt.Sprintf("%s/%s/%s", fileRepository.BaseDirectory, problemdir, filename)
	}
}

func (fileRepository *FilesRepository) OpenFile(problemDir, fileName string) (*os.File, error) {
	filePath := fileRepository.GetFilePath(problemDir, fileName)
	return os.Open(filePath)
}

func (fileRepository *FilesRepository) SaveFile(problemDir, fileName string, sourceFile io.Reader) error {
	filepath := fmt.Sprintf("%s/%s/%s", fileRepository.BaseDirectory, problemDir, fileName)
	destFile, err := os.Create(filepath)
	if err != nil {
		return err
	}
	_, err = io.Copy(destFile, sourceFile)

	if err != nil {
		return err
	}
	return nil
	//if err != nil {
	//	logrus.WithFields(logrus.Fields{
	//		"file path": filepath,
	//	}).WithError(err).Debug("error trying to create the file")
	//}
}

func (fileRepository *FilesRepository) DeleteFile(problemDir, fileName string) error {
	filepath := fmt.Sprintf("%s/%s/%s", fileRepository.BaseDirectory, problemDir, fileName)
	return os.Remove(filepath)
	//if err != nil {
	//	logrus.WithFields(logrus.Fields{
	//		"file path": filepath,
	//	}).WithError(err).Debug("error trying to delete the file")
	//}
}

func main() {
	repo, _ := NewFileRepository("Florin")
	file, err := repo.OpenFile("123", "Solution.c")

	if err != nil {
		fmt.Println("first", err)
	}
	err = repo.SaveFile("123", "fisiernou.c", file)
	fmt.Println(err)
}
