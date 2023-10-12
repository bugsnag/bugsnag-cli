# Upgrading Guide

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