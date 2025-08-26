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
