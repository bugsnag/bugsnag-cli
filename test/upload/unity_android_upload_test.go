package upload_testing

import (
	"archive/zip"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/bugsnag/bugsnag-cli/pkg/options"
	"github.com/bugsnag/bugsnag-cli/pkg/upload"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// A real ELF shared object with a GNU build ID, used as a stand-in for libil2cpp.so.
const elfFixture = "../../features/android/fixtures/app/build/intermediates/merged_native_libs/release/out/lib/arm64-v8a/libbugsnag-ndk.so"

// writeSymbolsZip creates a Unity-style <name>.symbols.zip in dir containing the given
// files (zip path → source file) and returns its path.
func writeSymbolsZip(t *testing.T, dir string, files map[string]string) string {
	t.Helper()
	zipPath := filepath.Join(dir, "Game-1.0-v1.symbols.zip")
	f, err := os.Create(zipPath)
	require.NoError(t, err)
	defer f.Close()

	zw := zip.NewWriter(f)
	for name, src := range files {
		data, err := os.ReadFile(src)
		require.NoError(t, err)
		w, err := zw.Create(name)
		require.NoError(t, err)
		_, err = w.Write(data)
		require.NoError(t, err)
	}
	require.NoError(t, zw.Close())
	return zipPath
}

// Regression test: a library rejected by the server (e.g. Firebase's) must not stop
// libil2cpp.so and LineNumberMappings.json from being uploaded.
func TestProcessUnityAndroid_ContinuesAfterRejectedLibrary(t *testing.T) {
	buildDir := t.TempDir()

	junk := filepath.Join(t.TempDir(), "junk.so")
	require.NoError(t, os.WriteFile(junk, []byte("not an elf file"), 0644))

	writeSymbolsZip(t, buildDir, map[string]string{
		// Sorts before libil2cpp.so, so it is uploaded first.
		"arm64-v8a/libFirebaseCppAnalytics.so": junk,
		"arm64-v8a/libil2cpp.so":               elfFixture,
	})

	mappingDir := filepath.Join(buildDir, "Library", "Bee", "artifacts", "Android", "il2cppOutput", "cpp", "Symbols")
	require.NoError(t, os.MkdirAll(mappingDir, 0755))
	require.NoError(t, os.WriteFile(filepath.Join(mappingDir, "LineNumberMappings.json"), []byte("{}"), 0644))

	var mu sync.Mutex
	var received []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Path
		_, params, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
		if err == nil {
			mr := multipart.NewReader(r.Body, params["boundary"])
			for {
				part, err := mr.NextPart()
				if err != nil {
					break
				}
				if part.FormName() == "sharedObjectName" {
					b, _ := io.ReadAll(part)
					name = string(b)
				}
			}
		}

		mu.Lock()
		received = append(received, name)
		mu.Unlock()

		if strings.Contains(name, "Firebase") {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte("invalid file format"))
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	opts := options.CLI{
		Globals: options.Globals{ApiKey: "1234567890ABCDEF1234567890ABCDEF"},
		Upload: options.Upload{
			Timeout:          10,
			UploadAPIRootUrl: server.URL,
			UnityAndroid: options.UnityAndroid{
				Path:          []string{buildDir},
				ApplicationId: "com.so.example",
				VersionCode:   "1",
				VersionName:   "1.0",
			},
		},
	}

	logger := NewMockLogger()
	err := upload.ProcessUnityAndroid(opts, logger)

	// The rejected library is still reported as a failure...
	require.Error(t, err)
	assert.Contains(t, err.Error(), "libFirebaseCppAnalytics.so")
	assert.NotContains(t, err.Error(), "libil2cpp.so")

	// ...but every other file, including the line mappings, was uploaded.
	assert.Equal(t, []string{
		"libFirebaseCppAnalytics.so",
		"libil2cpp.so",
		"/unity-line-mappings",
	}, received)
}
