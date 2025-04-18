# Upgrading Guide

## v2.x to v3.x

### Breaking changes

When using the `upload dsym` and `upload react-native` commands, Xcode archives are now prioritized over dSYMs found inside the project from a build. This means that if you have both an archive and a build, the CLI will search the default archive location set in Xcode for the latest archive for that day and upload the dSYMs from the archive.

The `upload xcode-build` commands can be used if you want to always pick dSYM files from a project build instead.

### Deprecated options

The following options have been removed from the CLI:

| Command                       | Removed                  | Replacement        |
|-------------------------------|--------------------------| ------------------ |
| `upload`                       | `--fail-on-upload-error` |   | 


## v1.x to v2.x

### Deprecated options

The following options were marked as deprecated as they have been renamed:

| Command                       | Removed                | Replacement        |
| ----------------------------- | ---------------------- | ------------------ |
| `create-build`, `upload dart` | `--app-version`        | `--version-name`   | 
| `create-build`, `upload dart` | `--app-version-code`   | `--version-code`   | 
| `create-build`, `upload dart` | `--app-bundle-version` | `--bundle-version` | 
| `upload react-native-android` | `--version`            | `--version-name`   | 

`--version` now displays the version of the CLI.