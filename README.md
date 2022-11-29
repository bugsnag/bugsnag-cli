# Bugsnag CLI
[![Build status](https://badge.buildkite.com/4c42f3d6345b14ecdc243abcf974cad0cfd9844e1b0e5f2418.svg)](https://buildkite.com/bugsnag/bugsnag-cli)

Upload your Flutter app's Dart symbol files to [Bugsnag](https://www.bugsnag.com/platforms/flutter) using our CLI.

## Installation

The binaries are available on our [GitHub releases page](https://github.com/bugsnag/bugsnag-cli/releases) for macOS, Linux and Windows.

To install or upgrade to the latest binary for your architecture, you can also run the following `cURL` or `Wget` commands:

```sh
curl -o- https://raw.githubusercontent.com/bugsnag/bugsnag-cli/main/install.sh | bash
```
```sh
wget -qO- https://raw.githubusercontent.com/bugsnag/bugsnag-cli/main/install.sh | bash
```

The script downloads the appropriate binary and attempts to install it to `~/.local/bugsnag`.

## Usage

See the our [online docs](https://docs.bugsnag.com/build-integrations/bugsnag-cli/) for full usage information.

```
Usage: bugsnag-cli <command>

Flags:
  -h, --help                                                Show context-sensitive help.
      --upload-api-root-url="https://upload.bugsnag.com"    Bugsnag On-Premise upload server URL. Can contain port number
      --port=443                                            Port number for the upload server
      --api-key=STRING                                      Bugsnag integration API key for this application
      --fail-on-upload-error                                Stops the upload when a mapping file fails to upload to Bugsnag successfully

Commands:
  upload all <path>
    Upload any symbol/mapping files

  upload dart <path>
    Process and upload symbol files for Flutter

Run "bugsnag-cli <command> --help" for more information on a command.
```

## Bugsnag On-Premise

If you are using Bugsnag On-premise, you should use the `--upload-api-root-url` option to set the url of your [upload server](https://docs.bugsnag.com/on-premise/single-machine/service-ports/#bugsnag-upload-server), for example:

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
