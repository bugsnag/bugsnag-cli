package utils

import (
	"bufio"
	"compress/gzip"
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

	newFile := file + ".gz"

	gzipFile, err := os.Create(newFile)

	if err != nil {
		return "", err
	}

	w := gzip.NewWriter(gzipFile)
	io.Copy(w, read)

	w.Close()

	return newFile, nil
}
