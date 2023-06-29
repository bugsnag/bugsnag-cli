# Changelog

## 1.2.0 (2023-06-29)

Add support for installing the CLI via NPM - [40](https://github.com/bugsnag/bugsnag-cli/pull/40)

Move global `appVersion`, `appVersionCode` and `appBundleVersion` flags to sub commands for `dart` and `create-build` - [41](https://github.com/bugsnag/bugsnag-cli/pull/41)

Correct `buildUUID` name in server requests - [41](https://github.com/bugsnag/bugsnag-cli/pull/41)

## 1.1.1 (2023-05-25)   

Fix how we check for the AndroidManifest.xml file for Android AAB - [37](https://github.com/bugsnag/bugsnag-cli/pull/37)

## 1.1.0 (2023-05-10)

Add support for:
- React Native source maps for Android - see our [online docs](https://docs.bugsnag.com/build-integrations/bugsnag-cli/upload-rn-android/)
- Android AAB files - see our [online docs](https://docs.bugsnag.com/build-integrations/bugsnag-cli/upload-android-aab/)
- Android NDK symbol files - see our [online docs](https://docs.bugsnag.com/build-integrations/bugsnag-cli/upload-android-ndk/)
- Android Proguard mapping files - see our [online docs](https://docs.bugsnag.com/build-integrations/bugsnag-cli/upload-android-proguard/)

Add the `create-build` command to provide extra information whenever you build, release, or deploy your application. - see our [online docs](https://docs.bugsnag.com/build-integrations/bugsnag-cli/create-build/)

## 1.0.0 (2022-11-29)

Initial release with support for Dart symbol files – see our [online docs](https://docs.bugsnag.com/build-integrations/bugsnag-cli/upload-dart/).
