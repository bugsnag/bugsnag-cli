Feature: React Native 0.72 Integration Tests

  Scenario: Build and Upload React Native 0.72 sourcemaps
    When I make the "features/base-fixtures/rn0_72/android"
    And I wait for the build to succeed

    When I run bugsnag-cli with upload react-native --upload-api-root-url=http://localhost:9339 --api-key=1234567890ABCDEF1234567890ABCDEF --overwrite features/base-fixtures/rn0_72
    And I wait to receive 1 sourcemaps
    Then the sourcemap is valid for the React Native Build API
    Then the sourcemaps Content-Type header is valid multipart form-data
    And the sourcemap payload field "apiKey" equals "1234567890ABCDEF1234567890ABCDEF"
    And the sourcemap payload field "appVersion" equals "1.0"
    And the sourcemap payload field "appVersionCode" equals "1"
    And the sourcemap payload field "platform" equals "android"
    And the sourcemap payload field "overwrite" equals "true"
