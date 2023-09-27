# Changelog

## 1.3.0 (TBD)

### Enhancements
 
- Add `--version` flag to the command to retrieve the version of the installed CLI. [51](https://github.com/bugsnag/bugsnag-cli/pull/51)
- Remove deprecated CLI options - `appVersion`, `appVersionCode` and `appBundleVersion`. [52](https://github.com/bugsnag/bugsnag-cli/pull/52)
- Android Build IDs can be calculated automatically for `.aab` files when none are specified in the `AndroidManifest` or on the command-line. [54](https://github.com/bugsnag/bugsnag-cli/pull/54)  
- Add `--dry-run` flag to all upload commands to validate but not upload source maps. [54](https://github.com/bugsnag/bugsnag-cli/pull/54)
- Add support for Unity Android symbol files. [56](https://github.com/bugsnag/bugsang-cli/pull/56)

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
