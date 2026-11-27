package upload_testing

import (
	"github.com/bugsnag/bugsnag-cli/pkg/elf"
	"testing"

	"github.com/bugsnag/bugsnag-cli/pkg/log"
	"github.com/bugsnag/bugsnag-cli/pkg/macho"
	"github.com/bugsnag/bugsnag-cli/pkg/options"
	"github.com/bugsnag/bugsnag-cli/pkg/upload"
	"github.com/bugsnag/bugsnag-cli/pkg/utils"
	"github.com/stretchr/testify/assert"
)

const (
	elfIosSymbolFile   = "../../features/dart/fixtures/app-debug-info/app.ios-arm64.symbols"
	machoIosSymbolFile = "../../features/dart/fixtures/macho/app-debug-info/app.ios-arm64.symbols"
	machoIosAppPath    = "../../features/dart/fixtures/macho/build/ios/iphoneos/Runner.app/Frameworks/App.framework/App"
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

func TestGetArchFromMachOFile(t *testing.T) {
	t.Log("Testing getting arch from a Mach-O symbol file")
	results, err := macho.GetArch(machoIosSymbolFile)

	assert.NoError(t, err)
	assert.Equal(t, "arm64", results, "Arch should match")
}

func TestGetArchFromMachOFileRejectsElfFile(t *testing.T) {
	t.Log("Testing getting arch from an ELF symbol file as Mach-O fails")
	_, err := macho.GetArch(elfIosSymbolFile)

	assert.Error(t, err)
}

func TestGetIosSymbolFileArch(t *testing.T) {
	tests := map[string]string{
		"Mach-O symbol file": machoIosSymbolFile,
		"ELF symbol file":    elfIosSymbolFile,
	}

	for name, symbolFile := range tests {
		t.Run(name, func(t *testing.T) {
			results, err := upload.GetIosSymbolFileArch(symbolFile)

			assert.NoError(t, err)
			assert.Equal(t, "arm64", results, "Arch should match")
		})
	}
}

func TestGetIosSymbolFileArchUnsupportedFile(t *testing.T) {
	t.Log("Testing getting arch from a file that is neither Mach-O nor ELF fails")
	_, err := upload.GetIosSymbolFileArch("dart_test.go")

	assert.Error(t, err)
}

func TestDwarfDumpUuidMachO(t *testing.T) {
	t.Log("Testing getting a build ID for a Mach-O symbol file")
	results, err := upload.DwarfDumpUuid(machoIosSymbolFile, machoIosAppPath, "arm64")

	assert.NoError(t, err)
	assert.Equal(t, "031CEFFD-4B20-3251-80D3-C561BE312E67", results, "UUID should match")
}

func TestDwarfDumpUuidArchNotInAppReturnsError(t *testing.T) {
	t.Log("Testing getting a build ID when the iOS app binary does not contain the arch")
	// dwarfdump prints nothing when the binary has no slice for the requested arch
	_, err := upload.DwarfDumpUuid(machoIosSymbolFile, machoIosAppPath, "x86_64")

	assert.Error(t, err)
}

func TestDartIosDryRun(t *testing.T) {
	tests := map[string]string{
		"Mach-O symbol file": machoIosSymbolFile,
		"ELF symbol file":    elfIosSymbolFile,
	}

	for name, symbolFile := range tests {
		t.Run(name, func(t *testing.T) {
			cli := options.CLI{}
			cli.DryRun = true
			cli.ApiKey = "1234567890ABCDEF1234567890ABCDEF"
			cli.Upload.DartSymbol = options.DartSymbol{Path: utils.Paths{symbolFile}}

			err := upload.Dart(cli, log.NewLoggerWrapper("debug"))

			assert.NoError(t, err)
		})
	}
}
