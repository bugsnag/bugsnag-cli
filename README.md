# Bugsnag CLI
[![Build status](https://badge.buildkite.com/4c42f3d6345b14ecdc243abcf974cad0cfd9844e1b0e5f2418.svg)](https://buildkite.com/bugsnag/bugsnag-cli)


## About

Tooling to process and upload symbol files to Bugsnag.

### Install & Update Script

To **install** or **update** bugsnag-cli, you should run the install script. To do that, you may either download and run the script manually, or use the following cURL or Wget command:
```sh
curl -o- https://raw.githubusercontent.com/bugsnag/bugsnag-cli/je/go-live-prep/install.sh | bash
```
```sh
wget -qO- https://raw.githubusercontent.com/bugsnag/bugsnag-cli/je/go-live-prep/install.sh | bash
```

Running either of the above commands downloads a script and runs it. The script downloads the latest release to `/usr/local/bin`.

You can also find the latest binary on the Github [releases](https://github.com/bugsnag/bugsnag-cli/releases) page.

## Usage

### Dart

Uploads symbol files generated from Flutter to Bugsnag

Example: `bugsnag-cli upload dart --api-key=YOUR_API_KEY path/to/symbol/files`

Full list of options as described by `bugsnag-cli upload dart --help`

```shell
Usage: bugsnag-cli upload dart <path>

Upload Dart symbol files

Arguments:
  <path>    Path to directory or file to upload

Flags:
  -h, --help                         Show context-sensitive help.
      --upload-api-root-url="https://upload.bugsnag.com"
                                     Bugsnag On-Premise upload server URL. Can contain port number
      --port=443                     Port number for the upload server
      --api-key=STRING               Bugsnag project API key
      --fail-on-upload-error         FailOnUploadError

      --overwrite                    ignore existing upload with same version
      --timeout=300                  seconds to wait before failing an upload request
      --retries=0                    number of retry attempts before failing a request

      --app-version=STRING           (optional) the version of the application.
      --app-version-code=STRING      (optional) the version code for the application (Android only).
      --app-bundle-version=STRING    (optional) the bundle version for the application (iOS only).
      --ios-app-path=STRING          (optional) the path to the built iOS app.
```