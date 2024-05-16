package android

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/bugsnag/bugsnag-cli/pkg/utils"
)

const MagicNumberByteCount = 8
const ChecksumByteCount = 4
const SignatureStartByte = MagicNumberByteCount + ChecksumByteCount
const SignatureByteCount = 20

const HeaderSize = MagicNumberByteCount + ChecksumByteCount + SignatureByteCount

func GetDexBuildId(dexDir string) string {
	signature, err := GetAppSignature(dexDir)
	if err != nil || signature == nil {
		return ""
	}

	return fmt.Sprintf("%x", signature)
}

// GetDexFiles - Given a list of paths, return an ordered list of .dex files suitable for building a signature
func GetDexFiles(dexPaths []string) ([]string, error) {
	var dexFiles []string

	for i := range dexPaths {
		path := dexPaths[i]
		if !utils.FileExists(path) {
			return nil, fmt.Errorf("no such file or directory: %s", path)
		}

		if utils.IsDir(path) {
			classesDexFiles := GetClassesDexFromDir(path)
			dexFiles = append(dexFiles, classesDexFiles...)
		} else if strings.HasSuffix(path, ".dex") {
			dexFiles = append(dexFiles, path)
		} else {
			return nil, fmt.Errorf("not a classesN.dex file: %s", path)
		}
	}

	return dexFiles, nil
}

func GetClassesDexFromDir(dexDir string) []string {
	filename := filepath.Join(dexDir, "classes.dex")
	if !utils.FileExists(filename) {
		return []string{}
	}

	dexFiles := []string{filename}
	for dexIndex := 2; ; dexIndex++ {
		filename = filepath.Join(dexDir, fmt.Sprintf("classes%s.dex", strconv.Itoa(dexIndex)))
		if !utils.FileExists(filename) {
			break
		} else {
			dexFiles = append(dexFiles, filename)
		}
	}
	return dexFiles
}

func GetAppSignatureFromFiles(dexFiles []string) ([]byte, error) {
	buildId := make([]byte, SignatureByteCount)

	for i := range dexFiles {
		filename := dexFiles[i]
		fileSignature, err := GetDexSignature(filename)
		if err != nil {
			break
		}

		buildId = MergeSignatures(buildId, fileSignature)
	}

	return buildId, nil
}

func GetAppSignature(dexDir string) ([]byte, error) {
	files, err := GetDexFiles([]string{dexDir})
	if err != nil {
		return nil, err
	}

	return GetAppSignatureFromFiles(files)
}

func MergeSignatures(buildId []byte, signature []byte) []byte {
	output := make([]byte, SignatureByteCount)
	for index := 0; index < SignatureByteCount; index++ {
		output[index] = buildId[index] ^ signature[index]
	}
	return output
}

func GetDexSignature(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	header := make([]byte, HeaderSize)
	count, err := file.Read(header)

	if err != nil {
		return nil, err
	}

	if count != HeaderSize {
		return nil, fmt.Errorf("invalid dex file '%s' expected %d header, but could only read %d bytes", path, HeaderSize, count)
	}

	err = ValidateHeader(header)
	if err != nil {
		return nil, err
	}

	dexSignature := bytes.Clone(header[SignatureStartByte : SignatureStartByte+SignatureByteCount])

	return dexSignature, nil
}

func ValidateHeader(header []byte) error {
	fileMagicNumber := header[0:MagicNumberByteCount]

	if fileMagicNumber[0] == 0x64 &&
		fileMagicNumber[1] == 0x65 &&
		fileMagicNumber[2] == 0x78 &&
		fileMagicNumber[3] == 0x0a &&
		// skip the version bytes and check that the magic number ends in a zero
		fileMagicNumber[7] == 0 {

		return nil
	}

	return fmt.Errorf("invalid dex file, bad magic number %x", fileMagicNumber)
}
