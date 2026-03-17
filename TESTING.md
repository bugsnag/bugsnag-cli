# Testing the BugSnag CLI

## Initial setup

Clone the repository and navigate to the project directory:

```sh
git clone git@github.com:bugsnag/bugsnag-cli.git
cd bugsnag-cli
```

### Prerequisites

- **Go** — Install via [Homebrew](https://brew.sh) or from [go.dev](https://go.dev/dl/):
  ```sh
  brew install golang
  ```
- **Ruby & Bundler** — Required for end-to-end tests:
  ```sh
  gem install bundler
  bundle install
  ```
- **Node.js & npm** — Required for npm linting and JavaScript fixture builds.

### Building the CLI

Build for the host platform:

```sh
make build
```

Build for all platforms (Windows, Linux, macOS):

```sh
make build-all
```

## Unit tests

Unit tests are written in Go and live under the `test/` directory, covering:

- `android/` — AAB manifest parsing, DEX build IDs, NDK paths/versions, objcopy, variant listing
- `endpoints/` — Endpoint configuration
- `ios/` — Plist parsing
- `options/` — Build metadata creation
- `upload/` — Android manifest, Android NDK, Breakpad, Dart, and JS uploads
- `utils/` — dSYM handling, file utilities, gzip compression

Run all unit tests:

```sh
make unit-test
```

This executes `go test -race ./test/...` with output formatted via [gotestfmt](https://github.com/gotesttools/gotestfmt).

## Linting

### Go

```sh
make go-lint
```

Runs [golangci-lint](https://golangci-lint.run/) across the codebase.

### Go formatting

```sh
make fmt
```

Applies `gofmt` to all Go source files.

### npm

```sh
make npm-lint
```

Runs `npm-check` against the `js/` package to detect outdated or unused dependencies.

## End-to-end tests

End-to-end tests use [Bugsnag Maze Runner](https://github.com/bugsnag/maze-runner) with feature files under `features/`. Step definitions are in `features/steps/steps.rb`.

### Test categories

| Category | Feature files | Description |
|---|---|---|
| CLI | `features/cli/` | General CLI behaviour, `create-build`, `--exclude`, installation |
| Android | `features/android/` | AAB, NDK, and ProGuard uploads |
| Xcode / dSYM | `features/xcode/` | dSYM uploads, Xcode build/archive, Swift Package Manager |
| JavaScript | `features/js/` | Source map uploads (Webpack 4 & 5, nested maps) |
| Dart | `features/dart/` | Dart symbol uploads |
| React Native | `features/react-native/` | Android & iOS React Native uploads |
| Node.js | `features/node/` | Node.js source map uploads |
| Unity | `features/unity/` | Android & iOS Unity uploads |
| Breakpad | `features/breakpad/` | Breakpad symbol uploads |

### Running end-to-end tests

Ensure the CLI is built first (`make build`), then run a specific category:

```sh
bundle exec maze-runner features/<category>
```

For example:

```sh
bundle exec maze-runner features/cli
bundle exec maze-runner features/xcode
bundle exec maze-runner features/android
```

#### Maze Runner port

The feature files reference `$MAZE_RUNNER_PORT` in upload/build API URLs. This defaults to `9339` if not set. To use a different port, export the variable before running tests:

```sh
export MAZE_RUNNER_PORT=9340
bundle exec maze-runner features/cli
```

### Building test fixtures

Some end-to-end tests require pre-built fixtures. Build them with:

```sh
make test-fixtures
```

Individual fixtures can also be built separately:

```sh
make features/base-fixtures/android       # Requires Android SDK / Gradle
make features/base-fixtures/dart          # Requires Flutter SDK
make features/base-fixtures/dsym          # Requires Xcode
make features/base-fixtures/js-webpack4   # Requires Node.js
make features/base-fixtures/js-webpack5   # Requires Node.js
```

## CI

The project uses **Buildkite** for continuous integration. The pipeline (`.buildkite/pipeline.yml`) runs:

1. License audit
2. Build
3. Go and npm linting
4. Unit tests
5. End-to-end tests across all categories

CI runs on macOS agents. Unity tests use isolated macOS agents.
