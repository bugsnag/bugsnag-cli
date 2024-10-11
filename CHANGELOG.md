# Changelog
## TBD

### Fixes
- Ensure that the node package is configured correctly so that you can run `npx @bugsnag/cli` and `yarn bugsnag-cli`. [144](https://github.com/bugsnag/bugsnag-cli/pull/144)

## 2.6.1 (2024-09-18)

### Fixes
- Ensure that we only pass either `--code-bundle-id` or `--version-code`/`--version-name`/`--bundle-version` to the upload API. [140](https://github.com/bugsnag/bugsnag-cli/pull/140)

## 2.6.0 (2024-09-09)

### Enhancements

- Add React Native super command [127](https://github.com/bugsnag/bugsnag-cli/pull/127)

### Fixes

- Allow spaces when processing and uploading dSYM files [135](https://github.com/bugsnag/bugsnag-cli/pull/135)

## 2.5.0 (2024-07-31)

### Enhancements

- Add support for JavaScript source maps [121](https://github.com/bugsnag/bugsnag-cli/pull/121)

## 2.4.1 (2024-07-17)

### Fixes

- Ensure that extracted `.aab` files can be processed by the Android AAB upload function [114](https://github.com/bugsnag/bugsnag-cli/pull/114)
- Dont show `--upload-api-root-url` and `--build-api-root-url` flags in the general help output [115](https://github.com/bugsnag/bugsnag-cli/pull/115)

## 2.4.0 (2024-07-08)

### Enhancements

- Restrict input for the `--provider` option for `create-build` [102](https://github.com/bugsnag/bugsnag-cli/pull/102)

### Fixes

- Ensure that the binary is correctly installed when installing via PNPM and Yarn [109](https://github.com/bugsnag/bugsnag-cli/pull/109)

## 2.3.0 (2024-06-04)

### Enhancements

- Add the ability to set the log leveL via the `--log-level` flag [103](https://github.com/bugsnag/bugsnag-cli/pull/103)
- Allow more flexible path searching when uploading NDK symbol files [98](https://github.com/bugsnag/bugsnag-cli/pull/98)

### Fixes

- Fix the error message when `--version-name` is missing [103](https://github.com/bugsnag/bugsnag-cli/pull/103)

## 2.2.0 (2024-04-17)

### Enhancements

- `upload android-proguard` will now attempt to automatically locate the `classes.dex` files if no build-uuid or dex-files are found or specified [92](https://github.com/bugsnag/bugsnag-cli/pull/92)
- Added the `--no-build-uuid` option to the `upload android-*` options [92](https://github.com/bugsnag/bugsnag-cli/pull/92)
- Added `Windows_NT` to `supported-platforms.yml` [95](https://github.com/bugsnag/bugsnag-cli/pull/95)

## 2.1.1 (2023-03-22)

### Fixes

- Ensure that the `--retries` flag is correctly passed to the Unity Android upload API. [91](https://github.com/bugsnag/bugsnag-cli/pull/91)

## 2.1.0 (2023-03-18)

### Deprecations
- The `--fail-on-upload-error` option now has no affect: upload commands will now all return a non-zero exit code if the upload is unsuccessful. All 4xx and 5xx status codes from the upload API are treated as errors apart from duplicate files (409), which the command will not treat as an error case to allow re-run commands to succeed. [95](https://github.com/bugsnag/bugsnag-cli/pull/90)

### Enhancements

- Add support for React Native source maps for iOS [online docs](https://docs.bugsnag.com/build-integrations/bugsnag-cli/upload-rn-ios/)
- Add support for dSYM uploads for iOS [online docs](https://docs.bugsnag.com/build-integrations/bugsnag-cli/upload-dsym/)
- Allow `create build` to extract relevant information from a given Android manifest or AAB file.[65](https://github.com/bugsnag/bugsnag-cli/pull/65)

### Fixes

- Ensure that `--ios-app-path` exists when passed as an option via the `upload dart` CLI. [67](https://github.com/bugsnag/bugsnag-cli/pull/67)
- Ensure that uploads are retried when passing the `--retries=x` argument to the CLI. [70](https://github.com/bugsnag/bugsnag-cli/pull/70)

## 2.0.0 (2023-10-17)

### Breaking Changes

- Remove deprecated (renamed) CLI options - `--version`, `--app-version`, `--app-version-code` and `--app-bundle-version`. [52](https://github.com/bugsnag/bugsnag-cli/pull/52)

See [Upgrading Guide](./UPGRADING.md) for full details.

### Enhancements

- Add support for Unity Android symbol files. [56](https://github.com/bugsnag/bugsnag-cli/pull/56)
- Add `--version` flag to the command to retrieve the version of the installed CLI. [51](https://github.com/bugsnag/bugsnag-cli/pull/51)
- Add `--dry-run` flag to all upload commands to validate but not upload source maps. [54](https://github.com/bugsnag/bugsnag-cli/pull/54)
- Automatically generate a unique value for the `buildUUID` parameter from `.aab` files when not specified in the `AndroidManifest` or `--build-uuid` option. [54](https://github.com/bugsnag/bugsnag-cli/pull/54)
- Add `--dex-files` flag to `upload android-proguard` to generate a unique value for the `buildUUID` from `classes.dex` files when uploading a `mapping.txt` [61](https://github.com/bugsnag/bugsnag-cli/pull/61)

## 1.2.2 (2023-07-11)

### Enhancements

- Do not modify the projects package.json when installing the CLI via NPM. [50](https://github.com/bugsnag/bugsnag-cli/pull/50)

- Adjust `index.android.bundle` path checking for React Native Android to ensure that paths are tested correctly. [49](https://github.com/bugsnag/bugsnag-cli/pull/49)

## 1.2.1 (2023-07-03)

### Enhancements

- Allow non-standard variants when not providing the bundle path as a flag to the CLI. [44](https://github.com/bugsnag/bugsnag-cli/pull/44)

- Add bundle path support for React Native 0.72. [46](https://github.com/bugsnag/bugsnag-cli/pull/46)

## 1.2.0 (2023-06-29)

### Enhancements

- Add support for installing the CLI via NPM - [39](https://github.com/bugsnag/bugsnag-cli/pull/39)

- Move global `appVersion`, `appVersionCode` and `appBundleVersion` flags to sub commands for `dart` and `create-build` - [41](https://github.com/bugsnag/bugsnag-cli/pull/41)

- Get values from Android AAB manifest via resource ID - [41](https://github.com/bugsnag/bugsnag-cli/pull/41)

### Fixes

- Correct `buildUUID` name in server requests for Android Proguard - [41](https://github.com/bugsnag/bugsnag-cli/pull/41)

## 1.1.1 (2023-05-25)

### Fixes

- Fix how we check for the AndroidManifest.xml file for Android AAB - [37](https://github.com/bugsnag/bugsnag-cli/pull/37)

## 1.1.0 (2023-05-10)

### Enhancements

Add support for:
- React Native source maps for Android - see our [online docs](https://docs.bugsnag.com/build-integrations/bugsnag-cli/upload-rn-android/)
- Android AAB files - see our [online docs](https://docs.bugsnag.com/build-integrations/bugsnag-cli/upload-android-aab/)
- Android NDK symbol files - see our [online docs](https://docs.bugsnag.com/build-integrations/bugsnag-cli/upload-android-ndk/)
- Android Proguard mapping files - see our [online docs](https://docs.bugsnag.com/build-integrations/bugsnag-cli/upload-android-proguard/)

Add the `create-build` command to provide extra information whenever you build, release, or deploy your application. - see our [online docs](https://docs.bugsnag.com/build-integrations/bugsnag-cli/create-build/)

## 1.0.0 (2022-11-29)

- Initial release with support for Dart symbol files – see our [online docs](https://docs.bugsnag.com/build-integrations/bugsnag-cli/upload-dart/).
