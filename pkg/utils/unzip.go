package utils

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func Unzip(path, outputPath string) error {
	archive, err := zip.OpenReader(path)
	if err != nil {
		return err
	}
	defer archive.Close()

	for _, f := range archive.File {
		filePath := filepath.Join(outputPath, f.Name)
		if !strings.HasPrefix(filePath, filepath.Clean(outputPath)+string(os.PathSeparator)) {
			return fmt.Errorf("invalid file path")
		}
		if f.FileInfo().IsDir() {
			err = os.MkdirAll(filePath, os.ModePerm)
			if err != nil {
				return err
			}
			continue
		}

		if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			return fmt.Errorf(err.Error())
		}

		outputPathFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return fmt.Errorf(err.Error())
		}

		fileInArchive, err := f.Open()
		if err != nil {
			return fmt.Errorf(err.Error())
		}

		if _, err := io.Copy(outputPathFile, fileInArchive); err != nil {
			return fmt.Errorf(err.Error())
		}

		outputPathFile.Close()
		fileInArchive.Close()
	}
	return nil
}
