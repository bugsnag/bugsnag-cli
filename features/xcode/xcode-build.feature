Feature: Xcode Build Integration Tests

  Scenario: Upload a single dSYM sourcemap using path containing one dSYM
    When I run bugsnag-cli with upload xcode-build --upload-api-root-url=http://localhost:$MAZE_RUNNER_PORT --api-key=1234567890ABCDEF1234567890ABCDEF --project-root=/my/project/root/ features/xcode/fixtures/single-dsym
    And I wait to receive 1 sourcemaps
    Then the sourcemap is valid for the dSYM Build API
    Then the sourcemaps Content-Type header is valid multipart form-data
    And the sourcemap payload field "apiKey" equals "1234567890ABCDEF1234567890ABCDEF"

  Scenario: Upload multiple dSYM sourcemaps using path containing multiple dSYMs
    When I run bugsnag-cli with upload xcode-build --upload-api-root-url=http://localhost:$MAZE_RUNNER_PORT --api-key=1234567890ABCDEF1234567890ABCDEF --project-root=/my/project/root/ features/xcode/fixtures/dsyms
    And I wait to receive 2 sourcemaps
    Then the sourcemap is valid for the dSYM Build API
    Then the sourcemaps Content-Type header is valid multipart form-data
    And the sourcemap payload field "apiKey" equals "1234567890ABCDEF1234567890ABCDEF"

  Scenario: Upload a single dSYM sourcemap using zip file containing one dSYM
    When I run bugsnag-cli with upload xcode-build --upload-api-root-url=http://localhost:$MAZE_RUNNER_PORT --api-key=1234567890ABCDEF1234567890ABCDEF --project-root=/my/project/root/ features/xcode/fixtures/single-dsym.zip
    And I wait to receive 1 sourcemaps
    Then the sourcemap is valid for the dSYM Build API
    Then the sourcemaps Content-Type header is valid multipart form-data
    And the sourcemap payload field "apiKey" equals "1234567890ABCDEF1234567890ABCDEF"

  Scenario: Upload multiple dSYM sourcemaps using zip file containing multiple dSYMs
    When I run bugsnag-cli with upload xcode-build --upload-api-root-url=http://localhost:$MAZE_RUNNER_PORT --api-key=1234567890ABCDEF1234567890ABCDEF --project-root=/my/project/root/ features/xcode/fixtures/dsyms.zip
    And I wait to receive 2 sourcemaps
    Then the sourcemap is valid for the dSYM Build API
    Then the sourcemaps Content-Type header is valid multipart form-data
    And the sourcemap payload field "apiKey" equals "1234567890ABCDEF1234567890ABCDEF"

  Scenario: Uploading a zip file containing directory of dSYM files that was compressed with macOS Archive Utility
    When I run bugsnag-cli with upload xcode-build --upload-api-root-url=http://localhost:$MAZE_RUNNER_PORT --api-key=1234567890ABCDEF1234567890ABCDEF --project-root=/my/project/root/ features/xcode/fixtures/macos-compressed-dsyms.zip
    And I wait to receive 2 sourcemaps
    Then the sourcemap is valid for the dSYM Build API
    Then the sourcemaps Content-Type header is valid multipart form-data
    And the sourcemap payload field "apiKey" equals "1234567890ABCDEF1234567890ABCDEF"

  Scenario: Upload symbols from an AppCenter zip
    When I run bugsnag-cli with upload xcode-build --upload-api-root-url=http://localhost:$MAZE_RUNNER_PORT --api-key=1234567890ABCDEF1234567890ABCDEF --project-root=/my/project/root/ "features/xcode/fixtures/app-center/symbols.zip"
    And I wait to receive 2 sourcemaps
    Then the sourcemap is valid for the dSYM Build API
    Then the sourcemaps Content-Type header is valid multipart form-data
    And the sourcemap payload field "apiKey" equals "1234567890ABCDEF1234567890ABCDEF"

  Scenario: Build and Upload dSYM
    When I make the "features/base-fixtures/dsym"
    Then I wait for the build to succeed

    When I run bugsnag-cli with upload xcode-build --upload-api-root-url=http://localhost:$MAZE_RUNNER_PORT --api-key=1234567890ABCDEF1234567890ABCDEF --project-root=/my/project/root/ features/base-fixtures/dsym/
    And I wait to receive 1 sourcemaps
    Then the sourcemap is valid for the dSYM Build API
    Then the sourcemaps Content-Type header is valid multipart form-data
    And the sourcemap payload field "apiKey" equals "1234567890ABCDEF1234567890ABCDEF"