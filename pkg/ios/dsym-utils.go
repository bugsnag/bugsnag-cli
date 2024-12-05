package ios

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/bugsnag/bugsnag-cli/pkg/log"
	"github.com/bugsnag/bugsnag-cli/pkg/utils"
)

// DwarfInfo stores the UUID, architecture, name, and location of a DWARF file.
// This information is extracted from dSYM files during processing.
type DwarfInfo struct {
	UUID     string
	Arch     string
	Name     string
	Location string
}

// FindDsymsInPath locates dSYM files within a specified path, processes them,
// and retrieves DWARF information for further use.
//
// Parameters:
// - path: The directory or file path to search for dSYM files.
// - ignoreEmptyDsym: If true, skips empty dSYM files without raising an error.
// - ignoreMissingDwarf: If true, skips invalid DWARF files without raising an error.
// - logger: Logger instance for informational and debug messages.
//
// Returns:
// - A slice of DwarfInfo structs containing details of found DWARF files.
// - A temporary directory if a ZIP file was extracted during processing.
// - An error if any issues occur during the process.
func FindDsymsInPath(path string, ignoreEmptyDsym, ignoreMissingDwarf bool, logger log.Logger) ([]*DwarfInfo, string, error) {
	var tempDir string
	var dsymLocations []string
	var dwarfInfo []*DwarfInfo

	// Process the provided path
	if utils.IsDir(path) {
		dsymLocations = findDsyms(path)
	} else if strings.HasSuffix(strings.ToLower(path), ".zip") {
		logger.Debug(fmt.Sprintf("Unzipping %s to search for dSYM files", filepath.Base(path)))
		var err error
		tempDir, err = utils.ExtractFile(path, "dsym")
		if err != nil {
			return nil, tempDir, fmt.Errorf("failed to unzip %s: %w", filepath.Base(path), err)
		}
		logger.Debug(fmt.Sprintf("Extracted %s to %s", filepath.Base(path), tempDir))
		dsymLocations = findDsyms(tempDir)
	} else {
		dsymLocations = append(dsymLocations, path)
	}

	// Ensure dwarfdump is available
	if len(dsymLocations) > 0 && !isDwarfDumpInstalled() {
		return nil, tempDir, fmt.Errorf("dwarfdump is not installed or unavailable")
	}

	// Process each located dSYM
	for _, dsymLocation := range dsymLocations {
		if info, err := processDsymLocation(dsymLocation, ignoreEmptyDsym, ignoreMissingDwarf, logger); err != nil {
			return nil, tempDir, err
		} else {
			dwarfInfo = append(dwarfInfo, info...)
		}
	}

	return dwarfInfo, tempDir, nil
}

// processDsymLocation extracts DWARF information from a specific dSYM file or directory.
//
// Parameters:
// - dsymLocation: The path to the dSYM file or directory.
// - ignoreEmptyDsym: If true, skips empty dSYM files without raising an error.
// - ignoreMissingDwarf: If true, skips invalid DWARF files without raising an error.
// - logger: Logger instance for informational and debug messages.
//
// Returns:
// - A slice of DwarfInfo structs containing details of DWARF files.
// - An error if the location cannot be processed or if invalid files are found.
func processDsymLocation(dsymLocation string, ignoreEmptyDsym, ignoreMissingDwarf bool, logger log.Logger) ([]*DwarfInfo, error) {
	var dwarfInfo []*DwarfInfo
	files, err := os.ReadDir(dsymLocation)

	if err != nil && strings.Contains(err.Error(), "not a directory") {
		// Process a single file
		fileName := filepath.Base(dsymLocation)
		return getDwarfFileInfo(filepath.Dir(dsymLocation), fileName), nil
	} else if err != nil {
		return nil, fmt.Errorf("failed to read dSYM location %s: %w", dsymLocation, err)
	}

	// Process all files in the directory
	for _, file := range files {
		filePath := filepath.Join(dsymLocation, file.Name())
		if fileInfo, _ := os.Stat(filePath); fileInfo != nil && fileInfo.Size() > 0 {
			info := getDwarfFileInfo(dsymLocation, file.Name())
			if len(info) == 0 && !ignoreMissingDwarf {
				return nil, fmt.Errorf("%s is not a valid DWARF file", fileInfo.Name())
			}
			dwarfInfo = append(dwarfInfo, info...)
		} else if fileInfo == nil || fileInfo.Size() == 0 {
			if ignoreEmptyDsym {
				logger.Info(fmt.Sprintf("%s is empty, skipping", file.Name()))
			} else {
				return nil, fmt.Errorf("%s is empty", file.Name())
			}
		}
	}

	return dwarfInfo, nil
}

// isDwarfDumpInstalled checks if the `dwarfdump` utility is available on the system.
//
// Returns:
// - `true` if the `dwarfdump` command is found in the system's executable path.
// - `false` otherwise.
func isDwarfDumpInstalled() bool {
	return utils.LocationOf(utils.DWARFDUMP) != ""
}

// getDwarfFileInfo retrieves DWARF file information from the output of the `dwarfdump` utility.
//
// Parameters:
// - path: The directory path containing the DWARF file.
// - fileName: The name of the DWARF file to be analyzed.
//
// Returns:
// - A slice of DwarfInfo structs containing extracted DWARF information.
func getDwarfFileInfo(path, fileName string) []*DwarfInfo {
	var dwarfInfo []*DwarfInfo
	cmd := exec.Command(utils.DWARFDUMP, "-u", fileName)
	cmd.Dir = path

	output, _ := cmd.Output()
	for _, line := range strings.Split(strings.TrimSpace(string(output)), "\n") {
		if strings.Contains(line, "UUID: ") {
			parts := strings.Fields(strings.ReplaceAll(strings.ReplaceAll(line, "(", ""), ")", ""))
			if len(parts) >= 4 {
				dwarfInfo = append(dwarfInfo, &DwarfInfo{
					UUID:     parts[1],
					Arch:     parts[2],
					Name:     strings.Join(parts[3:], " "),
					Location: path,
				})
			}
		}
	}
	return dwarfInfo
}

// findDsyms recursively searches a directory for dSYM files.
//
// Parameters:
// - root: The root directory to search.
//
// Returns:
// - A slice of strings representing the paths to the located dSYM files.
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
