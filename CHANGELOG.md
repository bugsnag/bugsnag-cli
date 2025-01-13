# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),  
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Add a wrapper for the `npm` package to interact with the BugSnag CLI.

## [2.8.0] - 2025-01-06

### Added
- Xcode Archive Command to support uploading `.xcarchive` files. [#156](https://github.com/bugsnag/bugsnag-cli/pull/156)

### Changed
- Rename the `dsym` upload command to `xcode-build` to better reflect the command's purpose. The `dsym` command will be removed in the next major release. [#156](https://github.com/bugsnag/bugsnag-cli/pull/156)

## [2.7.0] - 2024-11-26

### Added
- `--configuration` option to the `upload dsym` command. [#154](https://github.com/bugsnag/bugsnag-cli/pull/154)

## [2.6.3] - 2024-11-26

### Added
- Default the `--project-root` to the current working directory for the `upload dsym` command. [#148](https://github.com/bugsnag/bugsnag-cli/pull/148)

### Fixed
- Add the `--code-bundle-id` option to the `upload js` command. [#150](https://github.com/bugsnag/bugsnag-cli/pull/150)

## [2.6.2] - 2024-10-17

### Fixed
- Ensure the Node package is configured correctly to run `npx @bugsnag/cli` and `yarn bugsnag-cli`. [#144](https://github.com/bugsnag/bugsnag-cli/pull/144)
- Replace the axios dependency with fetch to reduce package size. [#145](https://github.com/bugsnag/bugsnag-cli/pull/145)

## [2.6.1] - 2024-09-18

### Fixed
- Ensure only one of `--code-bundle-id` or `--version-code`/`--version-name`/`--bundle-version` is passed to the upload API. [#140](https://github.com/bugsnag/bugsnag-cli/pull/140)

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
- Remove `--upload-api-root-url` and `--build-api-root-url` flags from the general help output. [#115](https://github.com/bugsnag/bugsnag-cli/pull/115)

## [2.4.0] - 2024-07-08

### Added
- Restrict input for the `--provider` option in the `create-build` command. [#102](https://github.com/bugsnag/bugsnag-cli/pull/102)

### Fixed
- Ensure the binary installs correctly when using PNPM and Yarn. [#109](https://github.com/bugsnag/bugsnag-cli/pull/109)

## [2.3.0] - 2024-06-04

### Added
- Set the log level via the `--log-level` flag. [#103](https://github.com/bugsnag/bugsnag-cli/pull/103)
- More flexible path searching for NDK symbol file uploads. [#98](https://github.com/bugsnag/bugsnag-cli/pull/98)

### Fixed
- Correct error message when `--version-name` is missing. [#103](https://github.com/bugsnag/bugsnag-cli/pull/103)

## [2.2.0] - 2024-04-17

### Added
- Auto-locate `classes.dex` files for `upload android-proguard` when no build-UUID or dex-files are specified. [#92](https://github.com/bugsnag/bugsnag-cli/pull/92)
- `--no-build-uuid` option for `upload android-*` commands. [#92](https://github.com/bugsnag/bugsnag-cli/pull/92)
- Support for `Windows_NT` in `supported-platforms.yml`. [#95](https://github.com/bugsnag/bugsnag-cli/pull/95)

## [2.1.1] - 2023-03-22

### Fixed
- Ensure the `--retries` flag is passed to the Unity Android upload API. [#91](https://github.com/bugsnag/bugsnag-cli/pull/91)

## [2.1.0] - 2023-03-18

### Deprecated
- `--fail-on-upload-error` now has no effect. Upload commands will return non-zero exit codes on unsuccessful uploads. [#90](https://github.com/bugsnag/bugsnag-cli/pull/90)

### Added
- Support for React Native iOS source maps. [Docs](https://docs.bugsnag.com/build-integrations/bugsnag-cli/upload-rn-ios/)
- Support for iOS dSYM uploads. [Docs](https://docs.bugsnag.com/build-integrations/bugsnag-cli/upload-dsym/)

### Fixed
- Validate `--ios-app-path` for the `upload dart` CLI. [#67](https://github.com/bugsnag/bugsnag-cli/pull/67)

## [2.0.0] - 2023-10-17

### Breaking Changes
- Remove deprecated CLI options: `--version`, `--app-version`, `--app-version-code`, and `--app-bundle-version`. [#52](https://github.com/bugsnag/bugsnag-cli/pull/52)

See the [Upgrading Guide](./UPGRADING.md) for details.

### Added
- Support for Unity Android symbol files. [#56](https://github.com/bugsnag/bugsnag-cli/pull/56)
- `--dry-run` flag for validating uploads without executing them. [#54](https://github.com/bugsnag/bugsnag-cli/pull/54)

## [1.0.0] - 2022-11-29

### Added
- Initial release with support for Dart symbol files. [Docs](https://docs.bugsnag.com/build-integrations/bugsnag-cli/upload-dart/)
