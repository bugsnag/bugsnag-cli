Feature: Dsym Integration Tests

  Scenario: Upload a single Dsym sourcemap using some flags
    When I run bugsnag-cli with upload dsym --upload-api-root-url=http://localhost:9339 --api-key=1234567890ABCDEF1234567890ABCDEF --dev --scheme="SingleSchemeExample" --overwrite --version-name=1.0 features/base-fixtures/dsym/SingleSchemeExample
    And I wait to receive 1 sourcemaps
    Then the sourcemap is valid for the Dsym Build API
    Then the sourcemaps Content-Type header is valid multipart form-data
    And the sourcemap payload field "scheme" equals "SingleSchemeExample"
    And the sourcemap payload field "apiKey" equals "1234567890ABCDEF1234567890ABCDEF"
    And the sourcemap payload field "appVersion" equals "1.0"
    And the sourcemap payload field "dev" equals "true"
    And the sourcemap payload field "platform" equals "ios"
    And the sourcemap payload field "overwrite" equals "true"