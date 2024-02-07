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
func GetDsymsForUpload(path string) (*[]*DsymFile, error) {
	filesFound, _ := os.ReadDir(path)
	var tempDir string
	var dsymFiles []*DsymFile

	switch len(filesFound) {
	case 0:
		return nil, errors.Errorf("No dSYM files found in expected location '%s'", path)
	default:
		if isDwarfDumpInstalled() {
			for _, file := range filesFound {
				if strings.HasSuffix(file.Name(), ".zip") {
					log.Info("Attempting to unzip " + file.Name() + " before proceeding to upload")
					tempDir, _ = utils.ExtractFile(filepath.Join(path, file.Name()), "dsym")

					if tempDir != "" {
						log.Info("Unzipped " + file.Name() + " to " + tempDir + " for uploading")
						path = tempDir
					} else {
						log.Warn("Could not unzip " + file.Name() + " to a temporary directory, skipping")
						continue
					}
				}

				cmd := exec.Command("dwarfdump", "-u", strings.TrimSuffix(file.Name(), ".zip"))
				cmd.Dir = path

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
							if len(info) == 4 {
								dsymFile := &DsymFile{}
								dsymFile.UUID = info[1]
								dsymFile.Arch = info[2]
								dsymFile.Name = info[3]
								dsymFiles = append(dsymFiles, dsymFile)
							}
						}
					}
				} else {
					log.Info("Skipping upload for file without UUID: " + file.Name())
				}
			}
			if tempDir != "" {
				err := os.RemoveAll(tempDir)
				if err != nil {
					log.Warn("Could not remove temporary directory: " + tempDir)
				}
			}
		} else {
			return nil, errors.New("Unable to locate dwarfdump on this system.")
		}
	}

	return &dsymFiles, nil
}
