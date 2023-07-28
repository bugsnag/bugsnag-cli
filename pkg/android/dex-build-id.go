package android

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

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

func GetAppSignature(dexDir string) ([]byte, error) {
	buildId, err := GetDexSignature(filepath.Join(dexDir, "classes.dex"))
	if err != nil {
		return nil, err
	}

	for dexIndex := 2; ; dexIndex++ {
		filename := filepath.Join(dexDir, "classes"+strconv.Itoa(dexIndex)+".dex")
		if !utils.FileExists(filename) {
			break
		}

		secondarySignature, err := GetDexSignature(filename)
		if err != nil {
			break
		}

		MergeSignatures(buildId, secondarySignature)
	}

	return buildId, nil
}

func MergeSignatures(buildId []byte, signature []byte) {
	for index := 0; index < SignatureByteCount; index++ {
		buildId[index] = buildId[index] ^ signature[index]
	}
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
