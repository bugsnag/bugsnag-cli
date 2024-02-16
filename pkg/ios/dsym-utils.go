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

func FindDsymsInPath(path string) ([]*DwarfInfo, string, error) {
	var tempDir string
	var dsymLocations []string
	var dwarfInfo []*DwarfInfo

	// If path is set and is a directory
	if utils.IsDir(path) {
		// Check for dSYMs within it
		dsymLocations = findDsyms(path)

	} else {

		// If path is pointing to a .zip file, we will extract it and look for dSYMS within it to get dsymLocations
		if strings.HasSuffix(path, ".zip") {

			fileName := filepath.Base(path)
			log.Info("Attempting to unzip " + fileName + " before proceeding to upload")

			var err error
			tempDir, err = utils.ExtractFile(path, "dsym")

			if err != nil {
				// TODO: This will be downgraded to a warning with --fail-on-upload in near future
				log.Error("Could not unzip "+fileName+" to a temporary directory, skipping", 1)
			} else {
				log.Info("Unzipped " + fileName + " to " + tempDir + " for uploading")
				dsymLocations = findDsyms(tempDir)

			}

		} else if strings.HasSuffix(path, ".dSYM") {
			// If path points to a .dSYM file, then we will use it as is
			dsymLocations = append(dsymLocations, path)
		}

	}

	// If we have found dSYMs, use dwarfdump to get the UUID etc for each dSYM
	if len(dsymLocations) > 0 {
		if !isDwarfDumpInstalled() {
			return nil, tempDir, errors.New("Unable to locate dwarfdump on this system.")
		}

		for _, dsymLocation := range dsymLocations {
			filesFound, _ := os.ReadDir(dsymLocation)

			for _, file := range filesFound {
				dwarfInfo = append(dwarfInfo, getDwarfFileInfo(dsymLocation, file.Name())...)
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
	} else {
		log.Info("Skipping file without UUID: " + fileName)
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
		if strings.HasSuffix(info.Name(), ".dSYM") {
			dsyms = append(dsyms, filepath.Join(path, "Contents", "Resources", "DWARF"))
		}
		return nil
	})
	if err != nil {
		return nil
	}
	return dsyms
}
