Feature: Dart Integration Tests

  Scenario: Locate and Upload Dart Source Maps
    When I run bugsnag-cli with upload dart --upload-api-root-url=http://localhost:9339 --api-key=1234567890ABCDEF1234567890ABCDEF --overwrite features/fixtures/dart/app-debug-info
    And I wait to receive 4 sourcemaps

    Then the sourcemap is valid for the Dart Build API

    Then the sourcemaps Content-Type header is valid multipart form-data
    And the sourcemap payload field "apiKey" equals "1234567890ABCDEF1234567890ABCDEF"
    And the sourcemap payload field "overwrite" equals "true"
