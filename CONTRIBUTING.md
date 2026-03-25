# Contributing

Thank you for considering contributing to the Bugsnag CLI! This document outlines the process for contributing to this project and helps ensure a smooth collaboration.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [How Can I Contribute?](#how-can-i-contribute)
  - [Reporting Bugs](#reporting-bugs)
  - [Suggesting Enhancements](#suggesting-enhancements)
  - [Pull Requests](#pull-requests)
- [Development Setup](#development-setup)
- [Coding Standards](#coding-standards)
- [Testing](#testing)
- [Commit Messages](#commit-messages)
- [License](#license)

## Code of Conduct

We are committed to providing a welcoming and inclusive environment for all contributors. Please be respectful and considerate in all interactions.

## How Can I Contribute?

### Reporting Bugs

Before submitting a bug report:

- **Check existing issues** to see if the problem has already been reported
- **Use the latest version** of the CLI to verify the bug still exists
- **Collect information** about your environment (OS, architecture, CLI version)

When submitting a bug report, include:

- A clear and descriptive title
- Steps to reproduce the issue
- Expected behavior vs. actual behavior
- Any relevant logs or error messages
- Your environment details (OS, CLI version, Go version if building from source)

### Suggesting Enhancements

Enhancement suggestions are welcome! When suggesting a feature:

- Use a clear and descriptive title
- Provide a detailed description of the proposed functionality
- Explain why this enhancement would be useful
- Include examples of how the feature would work

### Pull Requests

We actively welcome your pull requests:

1. **Fork the repository** and create your branch from `next`
2. **Follow the development setup** instructions below
3. **Make your changes** following our coding standards
4. **Add or update tests** as appropriate
5. **Ensure all tests pass** (see [Testing](#testing))
6. **Update documentation** if you're changing functionality
7. **Run linters** to ensure code quality
8. **Submit your pull request** with a clear description of the changes

## Development Setup

### Prerequisites

You'll need the following tools installed:

- **Go** (1.19 or later) — [Installation guide](https://go.dev/dl/)
- **Ruby & Bundler** — For end-to-end tests
  ```sh
  gem install bundler
  bundle install
  ```
- **Node.js & npm** — For JavaScript components and linting

### Building the CLI

Build for your current platform:

```sh
make build
```

Build for all platforms:

```sh
make build-all
```

For detailed testing instructions, see [TESTING.md](TESTING.md).

## Coding Standards

### Go

- Follow standard Go conventions as outlined in [Effective Go](https://go.dev/doc/effective_go)
- Use `gofmt` to format your code:
  ```sh
  make fmt
  ```
- Run the Go linter before committing:
  ```sh
  make go-lint
  ```

### JavaScript/TypeScript

- Follow standard JavaScript/TypeScript conventions
- Check for outdated or unused dependencies:
  ```sh
  make npm-lint
  ```

### General Guidelines

- Write clear, self-documenting code with meaningful variable and function names
- Add comments for complex logic or non-obvious behavior
- Keep functions focused and reasonably sized
- Handle errors appropriately and provide useful error messages

## Testing

All contributions should include appropriate tests. This project uses:

- **Unit tests** (Go) — Located in `test/` directory
- **End-to-end tests** — Feature files in `features/` directory using Maze Runner

### Running Tests

Run unit tests:

```sh
make unit-test
```

Run end-to-end tests (ensure CLI is built first):

```sh
bundle exec maze-runner features/<category>
```

For comprehensive testing documentation, see [TESTING.md](TESTING.md).

### Test Requirements

- New features must include unit tests
- Bug fixes should include a test demonstrating the fix
- All existing tests must pass before a PR can be merged
- End-to-end tests should be added for user-facing functionality when appropriate

## Commit Messages

Write clear and meaningful commit messages:

- Use the present tense ("Add feature" not "Added feature")
- Use the imperative mood ("Move cursor to..." not "Moves cursor to...")
- Limit the first line to 72 characters or less
- Reference issues and pull requests where appropriate (e.g., "Fixes #123")

Example:
```
Add support for custom source map paths

- Implement --source-map-path flag
- Add unit tests for path resolution
- Update documentation

Fixes #456
```

## License

By contributing to the Bugsnag CLI, you agree that your contributions will be licensed under the [MIT License](LICENSE).

All submissions are subject to review and may be rejected or require changes before being merged.

---

## Questions?

If you have questions about contributing, please:

- Check the [documentation](https://docs.bugsnag.com/build-integrations/bugsnag-cli/)
- Review existing [issues](https://github.com/bugsnag/bugsnag-cli/issues) and [pull requests](https://github.com/bugsnag/bugsnag-cli/pulls)
- Open a discussion in the issues section

For security vulnerabilities, please see [SECURITY.md](SECURITY.md).

Thank you for contributing! 🎉

## Releases

### Preparation

Ensure that:
1. All PRs to be included in the release have been merged into `next`.
2. `CHANGELOG.md` details all changes relevant to end users and that PR links are correct.


#### Performing the release

1. Create the release branch from next. `git checkout -b release/vx.xx.xx`
2. Run `make VERSION=x.xx.x bump` to set the desired version number and date of release.
3. Add and commit the updated files into the release branch. 
    1. `git add main.go CHANGELOG.md install.sh js/package.json`
    2. `git commit -m'bump version to vx.xx.x`
    3. `git tag vx.xx.x`
    4. `git push`
4. Create a PR from the release branch to `main`.
    1. Use the `CHANGELOG.md` entries for the new version as the PR description.
2. Once merged, on GitHub, 'Draft a new release':
    1. Tag version - the one created when bumping the version (`vx.xx.x`)
    2. Target - generally `main` unless the release is a minor/patch for a previous major version for which we have a branch.
    3. Release title - as the Tag version
    4. Description - copy directly from `CHANGLEOG.md`, ensuring that the formatting looks correct in the preview.
    5. Run `make build-all` locally to generate the binaries and then attach them to the release.
3. Publish release
4. Update and push to `next`
    1. `git checkout main`
    2. `git pull`
    3. `git checkout next`
    4. `git merge origin/main`
    5. `git push --force`
5. Run `npm publish` from the `js/` folder to publish the NPM package.
