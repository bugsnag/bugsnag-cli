Feature: dSYM Upload Integration Tests

  Scenario: Upload a single dSYM sourcemap using path containing dSYM
    When I run bugsnag-cli with upload dsym --upload-api-root-url=http://localhost:9339 --api-key=1234567890ABCDEF1234567890ABCDEF --project-root=/my/project/root/ features/dsym/fixtures/single-dsym
    And I wait to receive 1 sourcemaps
    Then the sourcemap is valid for the dSYM Build API
    Then the sourcemaps Content-Type header is valid multipart form-data
    And the sourcemap payload field "apiKey" equals "1234567890ABCDEF1234567890ABCDEF"
    And the sourcemap payload field "platform" equals "ios"

  Scenario: Upload multiple dSYM sourcemaps using path containing dSYM
    When I run bugsnag-cli with upload dsym --upload-api-root-url=http://localhost:9339 --api-key=1234567890ABCDEF1234567890ABCDEF --project-root=/my/project/root/ features/dsym/fixtures/dsyms
    And I wait to receive 2 sourcemaps
    Then the sourcemap is valid for the dSYM Build API
    Then the sourcemaps Content-Type header is valid multipart form-data
    And the sourcemap payload field "apiKey" equals "1234567890ABCDEF1234567890ABCDEF"
    And the sourcemap payload field "platform" equals "ios"