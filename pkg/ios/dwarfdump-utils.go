package ios

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"

	"github.com/bugsnag/bugsnag-cli/pkg/log"
	"github.com/bugsnag/bugsnag-cli/pkg/utils"
)

// DsymFile stores the UUID, architecture and name of a dSYM file
type DsymFile struct {
	UUID string
	Arch string
	Name string
}

// isDwarfDumpInstalled checks if dwarfdump is installed by checking if there is a path returned for it
func isDwarfDumpInstalled() bool {
	return utils.LocationOf(utils.DWARFDUMP) != ""
}

// GetDsymsForUpload returns information on the valid dSYM files found in a given path
func GetDsymsForUpload(path string, ignoreEmptyDsym bool) (*[]*DsymFile, error) {
	filesFound, _ := os.ReadDir(path)
	var dsymFiles []*DsymFile

	switch len(filesFound) {
	case 0:
		return nil, errors.Errorf("No dSYM files found in expected location '%s'", path)
	default:
		if isDwarfDumpInstalled() {
			for _, file := range filesFound {
				if strings.HasSuffix(file.Name(), ".zip") {
					log.Info("Attempting to unzip " + file.Name() + " before proceeding to upload")
					path, _ = utils.ExtractFile(filepath.Join(path, file.Name()), "dsym")
					defer func(path string) {
						_ = os.RemoveAll(path)
					}(path)

					if path != "" {
						log.Info("Unzipped " + file.Name() + " to " + path + " for uploading")
					}
				}

				fileInfo, _ := os.Stat(filepath.Join(path, strings.TrimSuffix(file.Name(), ".zip")))

				if fileInfo.Size() == 0 {
					if ignoreEmptyDsym {
						log.Warn("Skipping empty file: " + file.Name())
					} else {
						log.Error("Skipping empty file: "+file.Name(), 0)
					}

				} else {
					cmd := exec.Command("dwarfdump", "-u", filepath.Join(path, file.Name()))
					output, _ := cmd.Output()

					if len(output) > 0 {
						outputStr := string(output)

						outputStr = strings.TrimSuffix(outputStr, "\n")
						outputStr = strings.ReplaceAll(outputStr, "(", "")
						outputStr = strings.ReplaceAll(outputStr, ")", "")

						outputSlice := strings.Split(outputStr, "\n")

						for _, str := range outputSlice {
							if strings.Contains(str, "UUID: ") {
								info := strings.Split(str, " ")
								dsymFile := &DsymFile{}
								dsymFile.UUID = info[1]
								dsymFile.Arch = info[2]
								dsymFile.Name = filepath.Base(info[3])
								dsymFiles = append(dsymFiles, dsymFile)
							}
						}
					} else {
						log.Info("Skipping file without UUID: " + file.Name())
					}
				}
			}
		}
	}

	return &dsymFiles, nil
}
