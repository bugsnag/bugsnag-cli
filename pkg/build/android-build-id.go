package build

import (
	"fmt"
	"github.com/bugsnag/bugsnag-cli/pkg/android"
	"github.com/bugsnag/bugsnag-cli/pkg/utils"
)

type CreateAndroidBuildId struct {
	Path utils.Paths `arg:"" name:"path" help:"Path to the project directory" type:"path"`
}

func PrintAndroidBuildId(paths []string) error {
	dexFiles, err := android.GetDexFiles(paths)

	if err != nil {
		return err
	}

	signature, err := android.GetAppSignatureFromFiles(dexFiles)
	if err != nil {
		return err
	}

	fmt.Printf("%x", signature)
	fmt.Println()
	return nil
}
