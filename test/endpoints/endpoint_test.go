package endpoints

import (
	"github.com/bugsnag/bugsnag-cli/pkg/endpoints"
	"strings"
	"testing"

	"github.com/bugsnag/bugsnag-cli/pkg/options"
)

func TestBuildEndpointURL(t *testing.T) {
	tests := []struct {
		name      string
		uri       string
		port      int
		want      string
		expectErr bool
	}{
		{
			name: "Valid URI without port and non-zero port",
			uri:  "http://localhost",
			port: 8080,
			want: "http://localhost:8080",
		},
		{
			name: "Valid URI with existing port",
			uri:  "http://localhost:1234",
			port: 9999,
			want: "http://localhost:1234",
		},
		{
			name: "Zero port does not append",
			uri:  "http://localhost",
			port: 0,
			want: "http://localhost",
		},
		{
			name:      "Invalid URI returns error",
			uri:       "://%",
			port:      8080,
			expectErr: true,
		},
		{
			name: "Path is preserved with port",
			uri:  "https://example.com/path",
			port: 443,
			want: "https://example.com:443/path",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := endpoints.BuildEndpointURL(tt.uri, tt.port)
			if (err != nil) != tt.expectErr {
				t.Fatalf("BuildEndpointURL() error = %v, expectErr %v", err, tt.expectErr)
			}
			if err == nil && got != tt.want {
				t.Errorf("BuildEndpointURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetDefaultUploadEndpoint(t *testing.T) {
	tests := []struct {
		name         string
		apiKey       string
		endpointPath string
		options      options.CLI
		expectPrefix string
		expectErr    bool
	}{
		{
			name:         "Uses InsightHub endpoint with HUB prefix",
			apiKey:       endpoints.HUB_PREFIX + "123",
			endpointPath: "/upload",
			options: options.CLI{
				Globals: options.Globals{
					Port: 9999,
				},
			},
			expectPrefix: endpoints.HUB_UPLOAD + ":9999/upload",
		},
		{
			name:         "Uses Bugsnag endpoint with non-HUB key",
			apiKey:       "abc123",
			endpointPath: "/upload",
			options: options.CLI{
				Globals: options.Globals{
					Port: 0,
				},
			},
			expectPrefix: endpoints.BUGSNAG_UPLOAD + "/upload",
		},
		{
			name:         "Uses custom UploadAPIRootUrl if provided",
			apiKey:       "abc123",
			endpointPath: "/symbols",
			options: options.CLI{
				Globals: options.Globals{
					Port: 9999,
				},
				Upload: struct {
					Retries            int                            `help:"The number of retry attempts before failing an upload request" default:"0"`
					Timeout            int                            `help:"The number of seconds to wait before failing an upload request" default:"300"`
					UploadAPIRootUrl   string                         `help:"The upload server hostname, optionally containing port number"`
					All                options.DiscoverAndUploadAny   `cmd:"" help:"Upload any symbol/mapping files"`
					AndroidAab         options.AndroidAabMapping      `cmd:"" help:"Process and upload application bundle files for Android"`
					AndroidNdk         options.AndroidNdkMapping      `cmd:"" help:"Process and upload NDK symbol files for Android"`
					AndroidProguard    options.AndroidProguardMapping `cmd:"" help:"Process and upload Proguard/R8 mapping files for Android"`
					DartSymbol         options.DartSymbol             `cmd:"" help:"Process and upload symbol files for Flutter" name:"dart"`
					XcodeBuild         options.XcodeBuild             `cmd:"" help:"Upload dSYMs for iOS from a build"`
					Dsym               options.Dsym                   `cmd:"" help:"(deprecated) Upload dSYMs for iOS"`
					XcodeArchive       options.XcodeArchive           `cmd:"" help:"Upload dSYMs for iOS from a Xcode archive"`
					Js                 options.Js                     `cmd:"" help:"Upload source maps for JavaScript"`
					ReactNative        options.ReactNative            `cmd:"" help:"Upload source maps for React Native"`
					ReactNativeAndroid options.ReactNativeAndroid     `cmd:"" help:"Upload source maps for React Native Android"`
					ReactNativeIos     options.ReactNativeIos         `cmd:"" help:"Upload source maps for React Native iOS"`
					UnityAndroid       options.UnityAndroid           `cmd:"" help:"Upload Android mappings and NDK symbol files from Unity projects"`
					UnityIos           options.UnityIos               `cmd:"" help:"Upload iOS mappings and dSYMs from Unity projects"`
					Breakpad           options.Breakpad               `cmd:"" help:"Upload breakpad .sym files"`
				}(struct {
					Retries            int
					Timeout            int
					UploadAPIRootUrl   string
					All                options.DiscoverAndUploadAny
					AndroidAab         options.AndroidAabMapping
					AndroidNdk         options.AndroidNdkMapping
					AndroidProguard    options.AndroidProguardMapping
					DartSymbol         options.DartSymbol
					XcodeBuild         options.XcodeBuild
					Dsym               options.Dsym
					XcodeArchive       options.XcodeArchive
					Js                 options.Js
					ReactNative        options.ReactNative
					ReactNativeAndroid options.ReactNativeAndroid
					ReactNativeIos     options.ReactNativeIos
					UnityAndroid       options.UnityAndroid
					UnityIos           options.UnityIos
					Breakpad           options.Breakpad
				}{
					UploadAPIRootUrl: "https://custom.bugsnag.com",
				}),
			},
			expectPrefix: "https://custom.bugsnag.com:9999/symbols",
		},
		{
			name:         "Returns error on invalid custom URL",
			apiKey:       "abc123",
			endpointPath: "/bad",
			options: options.CLI{
				Globals: options.Globals{
					Port: 8080,
				},
				Upload: struct {
					Retries            int                            `help:"The number of retry attempts before failing an upload request" default:"0"`
					Timeout            int                            `help:"The number of seconds to wait before failing an upload request" default:"300"`
					UploadAPIRootUrl   string                         `help:"The upload server hostname, optionally containing port number"`
					All                options.DiscoverAndUploadAny   `cmd:"" help:"Upload any symbol/mapping files"`
					AndroidAab         options.AndroidAabMapping      `cmd:"" help:"Process and upload application bundle files for Android"`
					AndroidNdk         options.AndroidNdkMapping      `cmd:"" help:"Process and upload NDK symbol files for Android"`
					AndroidProguard    options.AndroidProguardMapping `cmd:"" help:"Process and upload Proguard/R8 mapping files for Android"`
					DartSymbol         options.DartSymbol             `cmd:"" help:"Process and upload symbol files for Flutter" name:"dart"`
					XcodeBuild         options.XcodeBuild             `cmd:"" help:"Upload dSYMs for iOS from a build"`
					Dsym               options.Dsym                   `cmd:"" help:"(deprecated) Upload dSYMs for iOS"`
					XcodeArchive       options.XcodeArchive           `cmd:"" help:"Upload dSYMs for iOS from a Xcode archive"`
					Js                 options.Js                     `cmd:"" help:"Upload source maps for JavaScript"`
					ReactNative        options.ReactNative            `cmd:"" help:"Upload source maps for React Native"`
					ReactNativeAndroid options.ReactNativeAndroid     `cmd:"" help:"Upload source maps for React Native Android"`
					ReactNativeIos     options.ReactNativeIos         `cmd:"" help:"Upload source maps for React Native iOS"`
					UnityAndroid       options.UnityAndroid           `cmd:"" help:"Upload Android mappings and NDK symbol files from Unity projects"`
					UnityIos           options.UnityIos               `cmd:"" help:"Upload iOS mappings and dSYMs from Unity projects"`
					Breakpad           options.Breakpad               `cmd:"" help:"Upload breakpad .sym files"`
				}(struct {
					Retries            int
					Timeout            int
					UploadAPIRootUrl   string
					All                options.DiscoverAndUploadAny
					AndroidAab         options.AndroidAabMapping
					AndroidNdk         options.AndroidNdkMapping
					AndroidProguard    options.AndroidProguardMapping
					DartSymbol         options.DartSymbol
					XcodeBuild         options.XcodeBuild
					Dsym               options.Dsym
					XcodeArchive       options.XcodeArchive
					Js                 options.Js
					ReactNative        options.ReactNative
					ReactNativeAndroid options.ReactNativeAndroid
					ReactNativeIos     options.ReactNativeIos
					UnityAndroid       options.UnityAndroid
					UnityIos           options.UnityIos
					Breakpad           options.Breakpad
				}{
					UploadAPIRootUrl: "http://%",
				}),
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := endpoints.GetDefaultUploadEndpoint(tt.apiKey, tt.endpointPath, tt.options)
			if (err != nil) != tt.expectErr {
				t.Fatalf("GetDefaultUploadEndpoint() error = %v, expectErr %v", err, tt.expectErr)
			}
			if err == nil && !strings.HasPrefix(got, tt.expectPrefix) {
				t.Errorf("GetDefaultUploadEndpoint() = %v, want prefix %v", got, tt.expectPrefix)
			}
		})
	}
}

func TestGetDefaultBuildEndpoint(t *testing.T) {
	tests := []struct {
		name         string
		apiKey       string
		options      options.CLI
		expectPrefix string
		expectErr    bool
	}{
		{
			name:   "Uses InsightHub build endpoint with HUB prefix",
			apiKey: endpoints.HUB_PREFIX + "XYZ",
			options: options.CLI{
				Globals: options.Globals{
					Port: 9000,
				},
			},
			expectPrefix: endpoints.HUB_BUILD,
		},
		{
			name:   "Uses Bugsnag build endpoint with non-HUB key",
			apiKey: "nohub",
			options: options.CLI{
				Globals: options.Globals{
					Port: 0,
				},
			},
			expectPrefix: endpoints.BUGSNAG_BUILD,
		},
		{
			name:   "Uses custom BuildApiRootUrl if provided",
			apiKey: "anykey",
			options: options.CLI{
				Globals: options.Globals{
					Port: 8081,
				},
				CreateBuild: options.CreateBuild{
					BuildApiRootUrl: "https://custom.build.smartbear.com",
				},
			},
			expectPrefix: "https://custom.build.smartbear.com:8081",
		},
		{
			name:   "Returns error on invalid custom BuildApiRootUrl",
			apiKey: "badkey",
			options: options.CLI{
				Globals: options.Globals{
					Port: 9999,
				},
				CreateBuild: options.CreateBuild{
					BuildApiRootUrl: "http://%",
				},
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := endpoints.GetDefaultBuildEndpoint(tt.apiKey, tt.options)
			if (err != nil) != tt.expectErr {
				t.Fatalf("GetDefaultBuildEndpoint() error = %v, expectErr %v", err, tt.expectErr)
			}
			if err == nil && !strings.HasPrefix(got, tt.expectPrefix) {
				t.Errorf("GetDefaultBuildEndpoint() = %v, want prefix %v", got, tt.expectPrefix)
			}
		})
	}
}
