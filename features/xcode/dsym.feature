Feature: dSYM Integration Tests
  @CleanAndBuildDsym
  Scenario: Upload a dsym from an Xcode Build
    When I run bugsnag-cli with upload dsym --upload-api-root-url=http://localhost:$MAZE_RUNNER_PORT --api-key=1234567890ABCDEF1234567890ABCDEF features/base-fixtures/dsym/
    And I wait to receive 1 sourcemaps
    Then the sourcemap is valid for the dSYM Build API
    Then the sourcemaps Content-Type header is valid multipart form-data
    And the sourcemap payload field "apiKey" equals "1234567890ABCDEF1234567890ABCDEF"

  @CleanAndArchiveDsym
  Scenario: Upload a dsym from an Xcode Archive
    When I run bugsnag-cli with upload dsym --upload-api-root-url=http://localhost:$MAZE_RUNNER_PORT --api-key=1234567890ABCDEF1234567890ABCDEF features/base-fixtures/dsym/
    And I wait to receive 1 sourcemaps
    Then the sourcemap is valid for the dSYM Build API
    Then the sourcemaps Content-Type header is valid multipart form-data
    And the sourcemap payload field "apiKey" equals "1234567890ABCDEF1234567890ABCDEF"

  @PathNormalization @AutoDetect @BugFix1
  Scenario: Upload dSYM with auto-detected project-root (omit flag - defaults to current directory)
    When I run bugsnag-cli with upload dsym --upload-api-root-url=http://localhost:$MAZE_RUNNER_PORT --api-key=1234567890ABCDEF1234567890ABCDEF features/base-fixtures/dsym/
    And I wait to receive 1 sourcemaps
    Then the sourcemap is valid for the dSYM Build API
    Then the sourcemaps Content-Type header is valid multipart form-data
    And the sourcemap payload field "apiKey" equals "1234567890ABCDEF1234567890ABCDEF"
    And the sourcemap payload field "projectRoot" is an absolute path

  @PathNormalization @AbsolutePathWithSlash @BugFix2
  Scenario: Upload dSYM with absolute project-root (WITH leading slash)
    When I run bugsnag-cli with upload dsym --upload-api-root-url=http://localhost:$MAZE_RUNNER_PORT --api-key=1234567890ABCDEF1234567890ABCDEF --project-root=/Users/user_name/project/sub-project features/base-fixtures/dsym/
    And I wait to receive 1 sourcemaps
    Then the sourcemap is valid for the dSYM Build API
    Then the sourcemaps Content-Type header is valid multipart form-data
    And the sourcemap payload field "apiKey" equals "1234567890ABCDEF1234567890ABCDEF"
    And the sourcemap payload field "projectRoot" equals "/Users/user_name/project/sub-project"

  @PathNormalization @RelativePathWithoutSlash @BugFix2
  Scenario: Upload dSYM with relative project-root (WITHOUT leading slash - normalized to absolute)
    When I run bugsnag-cli with upload dsym --upload-api-root-url=http://localhost:$MAZE_RUNNER_PORT --api-key=1234567890ABCDEF1234567890ABCDEF --project-root=Users/user_name/project/sub-project features/base-fixtures/dsym/
    And I wait to receive 1 sourcemaps
    Then the sourcemap is valid for the dSYM Build API
    Then the sourcemaps Content-Type header is valid multipart form-data
    And the sourcemap payload field "apiKey" equals "1234567890ABCDEF1234567890ABCDEF"
    And the sourcemap payload field "projectRoot" equals "/Users/user_name/project/sub-project"

  @PathNormalization @RelativePathWithDot @BugFix2
  Scenario: Upload dSYM with relative project-root starting with dot (normalized to absolute)
    When I run bugsnag-cli with upload dsym --upload-api-root-url=http://localhost:$MAZE_RUNNER_PORT --api-key=1234567890ABCDEF1234567890ABCDEF --project-root=./features/base-fixtures features/base-fixtures/dsym/
    And I wait to receive 1 sourcemaps
    Then the sourcemap is valid for the dSYM Build API
    Then the sourcemaps Content-Type header is valid multipart form-data
    And the sourcemap payload field "apiKey" equals "1234567890ABCDEF1234567890ABCDEF"
    And the sourcemap payload field "projectRoot" is an absolute path

  @PathNormalization @TrailingSlash @BugFix3
  Scenario: Upload dSYM with project-root ending in a single trailing slash (stripped before upload)
    When I run bugsnag-cli with upload dsym --upload-api-root-url=http://localhost:$MAZE_RUNNER_PORT --api-key=1234567890ABCDEF1234567890ABCDEF --project-root=/Users/user_name/project/sub-project/ features/base-fixtures/dsym/
    And I wait to receive 1 sourcemaps
    Then the sourcemap is valid for the dSYM Build API
    Then the sourcemaps Content-Type header is valid multipart form-data
    And the sourcemap payload field "apiKey" equals "1234567890ABCDEF1234567890ABCDEF"
    And the sourcemap payload field "projectRoot" equals "/Users/user_name/project/sub-project"

  @PathNormalization @DoubleTrailingSlash @BugFix3
  Scenario: Upload dSYM with project-root ending in double trailing slashes (stripped before upload)
    When I run bugsnag-cli with upload dsym --upload-api-root-url=http://localhost:$MAZE_RUNNER_PORT --api-key=1234567890ABCDEF1234567890ABCDEF --project-root=/Users/user_name/project/sub-project// features/base-fixtures/dsym/
    And I wait to receive 1 sourcemaps
    Then the sourcemap is valid for the dSYM Build API
    Then the sourcemaps Content-Type header is valid multipart form-data
    And the sourcemap payload field "apiKey" equals "1234567890ABCDEF1234567890ABCDEF"
    And the sourcemap payload field "projectRoot" equals "/Users/user_name/project/sub-project"

  @PathNormalization @ParentRelative @QAVerified
  Scenario: Upload dSYM with parent-relative project-root using .. (resolved to absolute)
    When I run bugsnag-cli with upload dsym --upload-api-root-url=http://localhost:$MAZE_RUNNER_PORT --api-key=1234567890ABCDEF1234567890ABCDEF --project-root=../.. features/base-fixtures/dsym/
    And I wait to receive 1 sourcemaps
    Then the sourcemap is valid for the dSYM Build API
    Then the sourcemaps Content-Type header is valid multipart form-data
    And the sourcemap payload field "apiKey" equals "1234567890ABCDEF1234567890ABCDEF"
    And the sourcemap payload field "projectRoot" is an absolute path

  @PathNormalization @ParentRelativeWithSlash @QAVerified
  Scenario: Upload dSYM with parent-relative project-root using ../../../ (resolved to absolute)
    When I run bugsnag-cli with upload dsym --upload-api-root-url=http://localhost:$MAZE_RUNNER_PORT --api-key=1234567890ABCDEF1234567890ABCDEF --project-root=../../../ features/base-fixtures/dsym/
    And I wait to receive 1 sourcemaps
    Then the sourcemap is valid for the dSYM Build API
    Then the sourcemaps Content-Type header is valid multipart form-data
    And the sourcemap payload field "apiKey" equals "1234567890ABCDEF1234567890ABCDEF"
    And the sourcemap payload field "projectRoot" is an absolute path

  @PathNormalization @DotSlash @QAVerified
  Scenario: Upload dSYM with project-root set to ./ (resolved to current working directory)
    When I run bugsnag-cli with upload dsym --upload-api-root-url=http://localhost:$MAZE_RUNNER_PORT --api-key=1234567890ABCDEF1234567890ABCDEF --project-root=./ features/base-fixtures/dsym/
    And I wait to receive 1 sourcemaps
    Then the sourcemap is valid for the dSYM Build API
    Then the sourcemaps Content-Type header is valid multipart form-data
    And the sourcemap payload field "apiKey" equals "1234567890ABCDEF1234567890ABCDEF"
    And the sourcemap payload field "projectRoot" is an absolute path

  @PathNormalization @DotOnly @QAVerified
  Scenario: Upload dSYM with project-root set to . (resolved to current working directory)
    When I run bugsnag-cli with upload dsym --upload-api-root-url=http://localhost:$MAZE_RUNNER_PORT --api-key=1234567890ABCDEF1234567890ABCDEF --project-root=. features/base-fixtures/dsym/
    And I wait to receive 1 sourcemaps
    Then the sourcemap is valid for the dSYM Build API
    Then the sourcemaps Content-Type header is valid multipart form-data
    And the sourcemap payload field "apiKey" equals "1234567890ABCDEF1234567890ABCDEF"
    And the sourcemap payload field "projectRoot" is an absolute path

  @PathNormalization @EmptyProjectRoot @QAVerified
  Scenario: Upload dSYM with project-root set to empty string (falls back to current working directory)
    When I run bugsnag-cli with upload dsym --upload-api-root-url=http://localhost:$MAZE_RUNNER_PORT --api-key=1234567890ABCDEF1234567890ABCDEF --project-root= features/base-fixtures/dsym/
    And I wait to receive 1 sourcemaps
    Then the sourcemap is valid for the dSYM Build API
    Then the sourcemaps Content-Type header is valid multipart form-data
    And the sourcemap payload field "apiKey" equals "1234567890ABCDEF1234567890ABCDEF"
    And the sourcemap payload field "projectRoot" is an absolute path
