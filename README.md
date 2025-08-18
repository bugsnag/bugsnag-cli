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

If your project uses `npm` or `yarn`, the CLI can be installed by adding the [`@bugsnag/cli`](https://www.npmjs.com/package/@bugsnag/cli) package:

```sh
npm install @bugsnag/cli`
```

It can then be executed from your project scripts at `/node_modules/.bin/bugsnag-cli` or using `npx @bugsnag/cli`.

## Supported commands

### Create builds

Allows you to create a build within BugSnag to enrich releases shown in the BugSnag dashboard.

    $ bugsnag-cli create-build --api-key=YOUR_API_KEY --app-version=YOUR_APP_VERSION

See the [`create-build`](https://docs.bugsnag.com/build-integrations/bugsnag-cli/create-build/) command reference for full usage information.

### Symbol &amp; mapping file uploads

Simplifies the upload of the various symbol and mapping files required to make your stacktraces readable in the BugSnag dashboard. Where possible files the files to upload are located automatically and the parameters, such as API key, located in project files. However all options can be overridden to allow you to customize the command for your build system.

Supported uploads with links to online docs for the file type:

* Android (obfuscation [mapping]((https://docs.bugsnag.com/build-integrations/bugsnag-cli/upload-android-proguard/)) and [native symbol]((https://docs.bugsnag.com/build-integrations/bugsnag-cli/upload-android-ndk/)) files, from builds or [AAB]((https://docs.bugsnag.com/build-integrations/bugsnag-cli/upload-android-aab/)) files)
* iOS (`.dSYM` from [builds]((https://docs.bugsnag.com/build-integrations/bugsnag-cli/upload-xcode-build/)) or [archives](https://docs.bugsnag.com/build-integrations/bugsnag-cli/upload-xcode-archive/))
* JavaScript source maps – for [web](https://docs.bugsnag.com/build-integrations/bugsnag-cli/upload-js/) and [React Native]((https://docs.bugsnag.com/build-integrations/bugsnag-cli/upload-rn/))
* Unity ([Android symbols](https://docs.bugsnag.com/build-integrations/bugsnag-cli/upload-unity-android/) or [iOS symbols](https://docs.bugsnag.com/build-integrations/bugsnag-cli/upload-unity-ios/))
* Dart ([stripped symbols](https://docs.bugsnag.com/build-integrations/bugsnag-cli/upload-dart/))
* Breakpad ([generated symbol files](https://docs.bugsnag.com/build-integrations/bugsnag-cli/upload-breakpad/))

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
