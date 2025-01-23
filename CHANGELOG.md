# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/), and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [2.9.1] - 2025-01-23

### Changed
- Switched `js-yaml` to `yaml` for the NodeJS package. [#168](https://github.com/bugsnag/bugsnag-cli/pull/168)

## [2.9.0] - 2025-01-20

### Added
- Wrapper for the `npm` package to interact with the BugSnag CLI. [#161](https://github.com/bugsnag/bugsnag-cli/pull/161)
- Support for `.so.*` files when processing NDK symbol files. [#163](https://github.com/bugsnag/bugsnag-cli/pull/163)
- Additional logging to the Android AAB upload command. [#165](https://github.com/bugsnag/bugsnag-cli/pull/165)

## [2.8.0] - 2025-01-06

### Added
- Xcode Archive Command to support the uploading of `.xcarchive` files. [#156](https://github.com/bugsnag/bugsnag-cli/pull/156)

### Changed
- Renamed the `dsym` upload command to `xcode-build` to better reflect its purpose. The `dsym` command will be removed in the next major release. [#156](https://github.com/bugsnag/bugsnag-cli/pull/156)

## [2.7.0] - 2024-11-26

### Added
- `--configuration` option to the `upload dsym` command. [#154](https://github.com/bugsnag/bugsnag-cli/pull/154)

## [2.6.3] - 2024-11-26

### Added
- Default the `--project-root` to the current working directory for the `upload dsym` command. [#148](https://github.com/bugsnag/bugsnag-cli/pull/148)

### Fixed
- Added the `--code-bundle-id` option to the `upload js` command. [#150](https://github.com/bugsnag/bugsnag-cli/pull/150)

## [2.6.2] - 2024-10-17

### Fixed
- Ensured the Node package is configured to run `npx @bugsnag/cli` and `yarn bugsnag-cli`. [#144](https://github.com/bugsnag/bugsnag-cli/pull/144)
- Replaced the `axios` dependency with `fetch` to reduce package size. [#145](https://github.com/bugsnag/bugsnag-cli/pull/145)

## [2.6.1] - 2024-09-18

### Fixed
- Ensure only one of `--code-bundle-id`, `--version-code`, `--version-name`, or `--bundle-version` is passed to the upload API. [#140](https://github.com/bugsnag/bugsnag-cli/pull/140)

## [2.6.0] - 2024-09-09

### Added
- React Native super command. [#127](https://github.com/bugsnag/bugsnag-cli/pull/127)

### Fixed
- Allow spaces when processing and uploading dSYM files. [#135](https://github.com/bugsnag/bugsnag-cli/pull/135)

## [2.5.0] - 2024-07-31

### Added
- Support for JavaScript source maps. [#121](https://github.com/bugsnag/bugsnag-cli/pull/121)

## [2.4.1] - 2024-07-17

### Fixed
- Ensure `.aab` files can be processed by the Android AAB upload function. [#114](https://github.com/bugsnag/bugsnag-cli/pull/114)
- Remove `--upload-api-root-url` and `--build-api-root-url` flags from general help output. [#115](https://github.com/bugsnag/bugsnag-cli/pull/115)

## [2.4.0] - 2024-07-08

### Added
- Restricted input for the `--provider` option for `create-build`. [#102](https://github.com/bugsnag/bugsnag-cli/pull/102)

### Fixed
- Ensure binary installation works correctly with PNPM and Yarn. [#109](https://github.com/bugsnag/bugsnag-cli/pull/109)

## [2.3.0] - 2024-06-04

### Added
- Ability to set the log level via the `--log-level` flag. [#103](https://github.com/bugsnag/bugsnag-cli/pull/103)
- Flexible path searching for NDK symbol file uploads. [#98](https://github.com/bugsnag/bugsnag-cli/pull/98)

### Fixed
- Correct error message when `--version-name` is missing. [#103](https://github.com/bugsnag/bugsnag-cli/pull/103)

## [2.2.0] - 2024-04-17

### Added
- `upload android-proguard` will automatically locate `classes.dex` files if none are specified. [#92](https://github.com/bugsnag/bugsnag-cli/pull/92)
- Added `--no-build-uuid` option for `upload android-*`. [#92](https://github.com/bugsnag/bugsnag-cli/pull/92)
- Added `Windows_NT` support in `supported-platforms.yml`. [#95](https://github.com/bugsnag/bugsnag-cli/pull/95)

## [2.1.1] - 2023-03-22

### Fixed
- Ensure `--retries` flag is correctly passed to the Unity Android upload API. [#91](https://github.com/bugsnag/bugsnag-cli/pull/91)

## [2.1.0] - 2023-03-18

### Deprecated
- `--fail-on-upload-error` no longer has an effect. All upload errors will return a non-zero exit code. [#90](https://github.com/bugsnag/bugsnag-cli/pull/90)

### Added
- Support for React Native source maps for iOS. [docs](https://docs.bugsnag.com/build-integrations/bugsnag-cli/upload-rn-ios/)
- Support for dSYM uploads for iOS. [docs](https://docs.bugsnag.com/build-integrations/bugsnag-cli/upload-dsym/)
- Ability for `create-build` to extract information from Android manifests and AAB files. [#65](https://github.com/bugsnag/bugsnag-cli/pull/65)

### Fixed
- Validate `--ios-app-path` in `upload dart`. [#67](https://github.com/bugsnag/bugsnag-cli/pull/67)
- Retry uploads when `--retries` is specified. [#70](https://github.com/bugsnag/bugsnag-cli/pull/70)

## [2.0.0] - 2023-10-17

### Breaking Changes
- Removed deprecated CLI options: `--version`, `--app-version`, `--app-version-code`, `--app-bundle-version`. See [Upgrading Guide](./UPGRADING.md). [#52](https://github.com/bugsnag/bugsnag-cli/pull/52)

### Added
- Support for Unity Android symbol files. [#56](https://github.com/bugsnag/bugsnag-cli/pull/56)
- `--version` flag to retrieve CLI version. [#51](https://github.com/bugsnag/bugsnag-cli/pull/51)
- `--dry-run` flag to validate source maps without uploading. [#54](https://github.com/bugsnag/bugsnag-cli/pull/54)
- Automatic `buildUUID` generation for `.aab` files. [#54](https://github.com/bugsnag/bugsnag-cli/pull/54)
- `--dex-files` flag for `upload android-proguard`. [#61](https://github.com/bugsnag/bugsnag-cli/pull/61)

## [1.2.2] - 2023-07-11

### Enhancements
- Do not modify the project's `package.json` when installing the CLI via NPM. [#50](https://github.com/bugsnag/bugsnag-cli/pull/50)
- Adjust `index.android.bundle` path checking for React Native Android to ensure that paths are tested correctly. [#49](https://github.com/bugsnag/bugsnag-cli/pull/49)

## [1.2.1] - 2023-07-03

### Enhancements
- Allow non-standard variants when not providing the bundle path as a flag to the CLI. [#44](https://github.com/bugsnag/bugsnag-cli/pull/44)
- Add bundle path support for React Native 0.72. [#46](https://github.com/bugsnag/bugsnag-cli/pull/46)

## [1.2.0] - 2023-06-29

### Enhancements
- Add support for installing the CLI via NPM. [#39](https://github.com/bugsnag/bugsnag-cli/pull/39)
- Move global `appVersion`, `appVersionCode`, and `appBundleVersion` flags to subcommands for `dart` and `create-build`. [#41](https://github.com/bugsnag/bugsnag-cli/pull/41)
- Get values from Android AAB manifest via resource ID. [#41](https://github.com/bugsnag/bugsnag-cli/pull/41)

### Fixes
- Correct `buildUUID` name in server requests for Android Proguard. [#41](https://github.com/bugsnag/bugsnag-cli/pull/41)

## [1.1.1] - 2023-05-25

### Fixes
- Fix how we check for the `AndroidManifest.xml` file for Android AAB. [#37](https://github.com/bugsnag/bugsnag-cli/pull/37)

## [1.1.0] - 2023-05-10

### Enhancements
Add support for:
- React Native source maps for Android - see our [online docs](https://docs.bugsnag.com/build-integrations/bugsnag-cli/upload-rn-android/)
- Android AAB files - see our [online docs](https://docs.bugsnag.com/build-integrations/bugsnag-cli/upload-android-aab/)
- Android NDK symbol files - see our [online docs](https://docs.bugsnag.com/build-integrations/bugsnag-cli/upload-android-ndk/)
- Android Proguard mapping files - see our [online docs](https://docs.bugsnag.com/build-integrations/bugsnag-cli/upload-android-proguard/)

Add the `create-build` command to provide extra information whenever you build, release, or deploy your application - see our [online docs](https://docs.bugsnag.com/build-integrations/bugsnag-cli/create-build/).

## [1.0.0] - 2022-11-29

- Initial release with support for Dart symbol files – see our [online docs](https://docs.bugsnag.com/build-integrations/bugsnag-cli/upload-dart/).
