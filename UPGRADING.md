# Upgrading Guide

## v2.x to v3.x

### Breaking changes

We now prioritize Xcode archives over Xcode builds when uploading dSYMs. This means that if you have both an archive and a build, the CLI will upload the dSYMs from the archive when using the `upload dsym` and `upload react-native` commands.

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