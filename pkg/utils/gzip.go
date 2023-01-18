package utils

import (
	"compress/gzip"
	"io"
	"os"
)

func GzipCompress(file string) error {
	fileToCompress, err := os.Open(file)
	if err != nil {
		return err
	}

	defer fileToCompress.Close()

	gzipWriter := gzip.NewWriter(fileToCompress)
	defer gzipWriter.Close()

	_, err = io.Copy(gzipWriter, fileToCompress)
	if err != nil {
		return err
	}

	return nil
}
