package utils

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"io"
	"os"
)

func GzipCompress(file string) (string, error) {
	fileData, err := os.Open(file)

	if err != nil {
		return "", err
	}

	read := bufio.NewReader(fileData)

	if err != nil {
		return "", err
	}

	newFile := fmt.Sprintf("%s.gz", file)

	gzipFile, err := os.Create(newFile)

	if err != nil {
		return "", err
	}

	w := gzip.NewWriter(gzipFile)
	defer w.Close()
	_, err = io.Copy(w, read)
	if err != nil {
		return "", err
	}

	return newFile, nil
}
