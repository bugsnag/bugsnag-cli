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

// DwarfInfo stores the UUID, architecture and name of a dwarf file
type DwarfInfo struct {
	UUID     string
	Arch     string
	Name     string
	Location string
}

func FindDsymsInPath(path string, ignoreEmptyDsym, ignoreMissingDwarf bool) ([]*DwarfInfo, string, error) {
	var tempDir string
	var dsymLocations []string
	var dwarfInfo []*DwarfInfo

	// If path is set and is a directory
	if utils.IsDir(path) {
		// Check for dSYMs within it
		dsymLocations = findDsyms(path)

	} else {

		// If path is pointing to a .zip file, we will extract it and look for dSYMS within it to get dsymLocations
		if strings.HasSuffix(strings.ToLower(path), ".zip") {

			fileName := filepath.Base(path)
			log.Info("Attempting to unzip " + fileName + " before proceeding to upload")

			var err error
			tempDir, err = utils.ExtractFile(path, "dsym")

			if err != nil {
				log.Error("Could not unzip "+fileName+" to a temporary directory, skipping", 1)

			} else {
				log.Info("Unzipped " + fileName + " to " + tempDir + " for uploading")
				dsymLocations = findDsyms(tempDir)
			}

		} else {
			// If path points to a file, then we will assume it is a dSYM and use it as-is
			dsymLocations = append(dsymLocations, path)
		}

	}

	// If we have found dSYMs, use dwarfdump to get the UUID etc for each dSYM
	if len(dsymLocations) > 0 {
		if !isDwarfDumpInstalled() {
			return nil, tempDir, errors.New("Unable to locate dwarfdump on this system.")
		}

		for _, dsymLocation := range dsymLocations {
			filesFound, err := os.ReadDir(dsymLocation)

			if err != nil {
				// If not a directory, then we'll assume that the path is pointing straight to a file
				if strings.Contains(err.Error(), "not a directory") {
					fileName := filepath.Base(dsymLocation)
					dsymLocation = filepath.Dir(dsymLocation)
					dwarfInfo = append(dwarfInfo, getDwarfFileInfo(dsymLocation, fileName)...)
				}
			}

			for _, file := range filesFound {
				fileInfo, _ := os.Stat(filepath.Join(dsymLocation, file.Name()))

				if fileInfo.Size() > 0 {
					info := getDwarfFileInfo(dsymLocation, file.Name())
					if len(info) == 0 {
						if ignoreMissingDwarf {
							log.Warn(fileInfo.Name() + " is not a valid DWARF file, skipping")
						} else {
							return nil, tempDir, errors.New(fileInfo.Name() + " is not a valid DWARF file, skipping")
						}
					}
					dwarfInfo = append(dwarfInfo, info...)
				} else {
					if ignoreEmptyDsym {
						log.Warn("Skipping empty file: " + file.Name())
					} else {
						return nil, tempDir, errors.New("Skipping empty file: " + file.Name())
					}
				}
			}
		}
	}

	return dwarfInfo, tempDir, nil
}

// isDwarfDumpInstalled checks if dwarfdump is installed by checking if there is a path returned for it
func isDwarfDumpInstalled() bool {
	return utils.LocationOf(utils.DWARFDUMP) != ""
}

// getDwarfFileInfo parses dwarfdump output to easier to manage/parsable DwarfInfo structs
func getDwarfFileInfo(path, fileName string) []*DwarfInfo {
	var dwarfInfo []*DwarfInfo

	cmd := exec.Command(utils.DWARFDUMP, "-u", strings.TrimSuffix(fileName, ".zip"))
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
				rawDwarfInfo := strings.Split(str, " ")
				if len(rawDwarfInfo) == 4 {
					dwarf := &DwarfInfo{}
					dwarf.UUID = rawDwarfInfo[1]
					dwarf.Arch = rawDwarfInfo[2]
					dwarf.Name = rawDwarfInfo[3]
					dwarf.Location = path
					dwarfInfo = append(dwarfInfo, dwarf)
				}
			}
		}
	}

	return dwarfInfo
}

// findDsyms walks the directory tree and returns a list of dSYM locations
func findDsyms(root string) []string {
	var dsyms []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// If the file is a dSYM, add it to the list (unless it resides within the __MACOSX directory)
		if strings.HasSuffix(strings.ToLower(info.Name()), ".dsym") && !strings.Contains(strings.ToLower(path), "__macosx") {
			dsyms = append(dsyms, filepath.Join(path, "Contents", "Resources", "DWARF"))
		}

		return nil
	})
	if err != nil {
		return nil
	}
	return dsyms
}
