package upload_testing

import (
	"github.com/bugsnag/bugsnag-cli/pkg/elf"
	"testing"

	"github.com/bugsnag/bugsnag-cli/pkg/upload"
	"github.com/stretchr/testify/assert"
)

func TestReadElfBuildId(t *testing.T) {
	t.Log("Testing getting a build ID from an ELF file")
	results, err := elf.GetBuildId("../../features/dart/fixtures//app-debug-info/app.android-arm64.symbols")
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, results, "07cc131ca803c124e93268ce19322737", "Build Ids should be the same")
}

func TestGetArchFromElfFile(t *testing.T) {
	t.Log("Testing getting arch from symbol file")

	results, err := elf.GetArch("../../features/dart/fixtures//app-debug-info/app.ios-arm64.symbols")

	t.Log(results)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, results, "arm64", "Arch should match")
}

func TestDwarfDumpUuid(t *testing.T) {
	t.Log("Testing getting a build ID from a Dwarf file")
	results, err := upload.DwarfDumpUuid("../../features/dart/fixtures/app-debug-info/app.ios-arm64.symbols", "../../features/dart/fixtures/build/ios/iphoneos/Runner.app/Frameworks/App.framework/App", "arm64")

	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, results, "E30C1BE5-DEB6-373C-98B4-52D827B7FF0D", "UUID should match")
}

func TestGetIosAppPath(t *testing.T) {
	t.Log("Testing getting the IOS app path from a given symbols path")
	results, err := upload.GetIosAppPath("../../features/dart/fixtures/app-debug-info/app.android-arm64.symbols")

	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, results, "../../features/dart/fixtures/build/ios/iphoneos/Runner.app/Frameworks/App.framework/App", "They should match")
}
