<div align="center">
  <a href="https://docs.bugsnag.com/build-integrations/bugsnag-cli">
    <picture>
      <source media="(prefers-color-scheme: dark)" srcset="https://assets.smartbear.com/m/3dab7e6cf880aa2b/original/BugSnag-Repository-Header-Dark.svg">
      <img alt="SmartBear BugSnag logo" src="https://assets.smartbear.com/m/3945e02cdc983893/original/BugSnag-Repository-Header-Light.svg">
    </picture>
  </a>
  <h1>CLI</h1>
</div>


[![Documentation](https://img.shields.io/badge/documentation-latest-blue.svg)](https://docs.bugsnag.com/build-integrations/bugsnag-cli/)
[![Build status](https://badge.buildkite.com/4c42f3d6345b14ecdc243abcf974cad0cfd9844e1b0e5f2418.svg)](https://buildkite.com/bugsnag/bugsnag-cli)

Simplify the process of creating releases on the BugSnag dashboard and uploading files to improve the stacktraces in your errors with our command line tool.

## Installation

The binaries are available on our [GitHub releases page](https://github.com/bugsnag/bugsnag-cli/releases) for macOS, Linux and Windows.

### cURL / Wget

To install or upgrade to the latest binary for your architecture, you can also run the following `cURL` or `Wget` commands:

```sh
curl -o- https://raw.githubusercontent.com/bugsnag/bugsnag-cli/main/install.sh | bash
```
```sh
wget -qO- https://raw.githubusercontent.com/bugsnag/bugsnag-cli/main/install.sh | bash
```

The script downloads the appropriate binary and attempts to install it to `~/.local/bugsnag`.

### NPM

To install or upgrade the BugSnag CLI via `npm`, you can run the following command:

`npm install @bugsnag/cli`

## Supported commands

This tool is currently being developed. It currently supports the following commands:

### Create builds

Allows you to create a build within BugSnag to enrich releases shown in the BugSnag dashboard.

    $ bugsnag-cli create-build --api-key=YOUR_API_KEY --app-version=YOUR_APP_VERSION

See the [`create-build`](https://docs.bugsnag.com/build-integrations/bugsnag-cli/create-build/) command reference for full usage information.

### Android NDK mapping files

For apps that use the [NDK](https://developer.android.com/ndk/), this command extracts symbols from `.so` files and uploads them along with version information.

    $ bugsnag-cli upload android-ndk \
        app/build/intermediates/merged_native_libs/release/out/lib/arm64-v8a/libMyApp.so

See the [`upload android-ndk`](https://docs.bugsnag.com/build-integrations/bugsnag-cli/upload-android-ndk/) command reference for full usage information.

### Android Proguard mapping flies

If you are using [ProGuard](https://developer.android.com/studio/build/shrink-code.html), [DexGuard](https://www.guardsquare.com/en/dexguard), or [R8](https://r8.googlesource.com/r8#d8-dexer-and-r8-shrinker) to minify and optimize your app, this command uploads the mapping file along with version information from your project directory:

    $ bugsnag-cli upload android-proguard app/build/outputs/proguard/release/mapping.txt

See the [`upload android-proguard`](https://docs.bugsnag.com/build-integrations/bugsnag-cli/upload-android-proguard/) command reference for full usage information.

### Android App Bundle (AAB) files

If you distribute your app as an [Android App Bundle](https://developer.android.com/guide/app-bundle) (AAB), they contain all required files and so can be uploaded in a single command:

    $ bugsnag-cli upload android-aab app/build/outputs/bundle/release/app-release.aab

See the [`upload android-aab`](https://docs.bugsnag.com/build-integrations/bugsnag-cli/upload-android-ndk/) command reference for full usage information.

### React Native Symbols

This command is a convenience wrapper for the following commands, which upload the different file types separately with more specific options available to customize the upload:
- `upload react-native-android`
- `upload react-native-ios`
- `upload android-ndk`
- `upload android-proguard`
- `upload dsym`


    $ bugsnag-cli upload react-native

### React Native JavaScript source maps (Android only)

To get unminified stack traces for JavaScript code in your React Native app built for Android, source maps must be generated and can be uploaded to BugSnag using the following command from the root of your project:

    $ bugsnag-cli upload react-native-android

See the [`upload react-native-android`](https://docs.bugsnag.com/build-integrations/bugsnag-cli/upload-rn-android/) command reference for full usage information.

### React Native JavaScript source maps (iOS only)

To get unminified stack traces for JavaScript code in your React Native app built for iOS, source maps must be generated and uploaded to BugSnag.

    $ bugsnag-cli upload react-native-ios

See the [`upload react-native-ios`](https://docs.bugsnag.com/build-integrations/bugsnag-cli/upload-rn-ios/) command reference for full usage information.

### JavaScript source maps

To get unminified stack traces for JavaScript code on the web, source maps must be generated and uploaded to BugSnag.

    $ bugsnag-cli upload js

See the [`upload js`](https://docs.bugsnag.com/build-integrations/bugsnag-cli/upload-js/) command reference for full usage information.

### Dart symbols for Flutter

If you are stripping debug symbols from your Dart code when building your Flutter apps, you will need to upload symbol files in order to see full stacktraces using the following command:

    $ bugsnag-cli upload dart --api-key=YOUR_API_KEY app-debug-info/

### dSYM files (iOS, macOS, tvOS)

Upload dSYM files from an Xcode project - either by locating an Xcode archive (.xcarchive) or from the project build – to allow BugSnag to show human-friendly function names, file paths, and line numbers in your iOS, macOS, and tvOS stacktraces.

    $ bugsnag-cli upload dsym

### dSYM files from an Xcode build (iOS, macOS, tvOS)

Upload dSYM files generated from a Build in Xcode to allow BugSnag to show human-friendly function names, file paths, and line numbers in your iOS, macOS, and tvOS stacktraces.

    $ bugsnag-cli upload xcode-build

### dSYM files from an Xcode Archive (iOS, macOS, tvOS)

Upload dSYM files generated from an Archive in Xcode to allow BugSnag to show human-friendly function names, file paths, and line numbers in your iOS, macOS, and tvOS stacktraces.

    $ bugsnag-cli upload xcode-archive

### Unity Symbol Files (Android only) 

The unity-android command uploads the IL2CPP symbols from the .symbols.zip file produced by the Unity build (see [Unity documentation](https://docs.unity3d.com/Manual/android-symbols.html) for more information) to the [NDK symbol API](https://d1upynpnqddd6j.cloudfront.net/api/ndk-symbol-mapping-upload/).

    $ bugsnag-cli upload unity-android /path/to/build/directory

### Breakpad symbol files

Upload [Breakpad](https://chromium.googlesource.com/breakpad/breakpad/) `.sym` files, generated by the `dump_syms` tool.

    $ bugsnag-cli upload breakpad /path/to/mylib.sym

See our [Breakpad documentation](https://docs.bugsnag.com/api/breakpad-symbol-upload/#generating-symbol-files) for more information on how to generate symbol files.

## BugSnag On-Premise

If you are using BugSnag On-premise, you should use the `--build-api-root-url` and `--upload-api-root-url` options to set the URL of your [build](https://docs.bugsnag.com/on-premise/single-machine/service-ports/#bugsnag-build-api) and [upload](https://docs.bugsnag.com/on-premise/single-machine/service-ports/#bugsnag-upload-server) servers, for example:

```sh
bugsnag-cli upload \
  --upload-api-root-url https://bugsnag.my-company.com/
  # ... other options
```

## Support

* Check out the [documentation](https://docs.bugsnag.com/build-integrations/bugsnag-cli/)
* [Search open and closed issues](https://github.com/bugsnag/bugsnag-cli/issues?q=+) for similar problems
* [Report a bug or request a feature](https://github.com/bugsnag/bugsnag-cli/issues/new)

## Contributing

Most updates to this repo will be made by Bugsnag employees. We are unable to accommodate significant external PRs such as features additions or any large refactoring, however minor fixes are welcome. See [contributing](CONTRIBUTING.md) for more information.

## License

This package is free software released under the MIT License. See [license](./LICENSE) for details.
