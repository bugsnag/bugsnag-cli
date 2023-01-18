package utils

import (
	"bufio"
	"compress/gzip"
	"io/ioutil"
	"os"
	"strings"
)

func GzipCompress(file string) (string, error) {
	fileData, err := os.Open(file)

	if err != nil {
		return "", err
	}

	read := bufio.NewReader(fileData)

	data, err := ioutil.ReadAll(read)

	if err != nil {
		return "", err
	}

	newFile := strings.Replace(file, ".txt", ".gz", -1)

	gzipFile, err := os.Create(newFile)

	if err != nil {
		return "", err
	}

	w := gzip.NewWriter(gzipFile)
	w.Write(data)

	w.Close()

	return newFile, nil
}
