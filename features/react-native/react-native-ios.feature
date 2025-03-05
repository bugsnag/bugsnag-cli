@BuildRNiOS
Feature: React Native iOS Integration Tests
#  TODO: Reinstate tests to cover all escape hatches to ensure proper functionality and edge case handling.
  Scenario: Upload a single React Native iOS sourcemap
    When I run bugsnag-cli with upload react-native-ios --upload-api-root-url=http://localhost:9339 --api-key=1234567890ABCDEF1234567890ABCDEF --overwrite features/react-native/fixtures/generated/old-arch/$RN_VERSION
    And I wait to receive 1 sourcemaps
    Then the sourcemap is valid for the React Native Build API
    Then the sourcemaps Content-Type header is valid multipart form-data
    And the sourcemap payload field "apiKey" equals "1234567890ABCDEF1234567890ABCDEF"
    And the sourcemap payload field "appVersion" equals "1.0"
    And the sourcemap payload field "appBundleVersion" equals "1"
    And the sourcemap payload field "platform" equals "ios"
    And the sourcemap payload field "overwrite" equals "true"

