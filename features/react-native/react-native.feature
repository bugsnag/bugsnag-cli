@BuildRN
Feature: React Native Integration Tests
  Scenario: Upload source maps for React Native
    When I run bugsnag-cli with upload react-native --upload-api-root-url=http://localhost:9339 --api-key=1234567890ABCDEF1234567890ABCDEF --overwrite features/react-native/fixtures/generated/old-arch/$RN_VERSION
    Given the following React Native versions and expected source maps:
      | React Native Version | Expected Source Maps |
      | 0.70               | 216                  |
      | 0.73               | 240                  |
      | 0.75               | 216                  |
    Then I wait to receive the correct number of sourcemaps
    Then the sourcemap is valid for the React Native Build API
    Then the sourcemaps Content-Type header is valid multipart form-data
    And the sourcemap payload field "apiKey" equals "1234567890ABCDEF1234567890ABCDEF"
    And the sourcemap payload field "appVersion" equals "1.0"
    And the sourcemap payload field "appVersionCode" equals "1"
    And the sourcemap payload field "platform" equals "android"
    And the sourcemap payload field "overwrite" equals "true"
    And I discard the oldest sourcemap

    Then the sourcemap is valid for the React Native Build API
    Then the sourcemaps Content-Type header is valid multipart form-data
    And the sourcemap payload field "apiKey" equals "1234567890ABCDEF1234567890ABCDEF"
    And the sourcemap payload field "appVersion" equals "1.0"
    And the sourcemap payload field "appBundleVersion" equals "1"
    And the sourcemap payload field "platform" equals "ios"
    And the sourcemap payload field "overwrite" equals "true"
    And I discard the oldest sourcemap

    Then the sourcemap is valid for the Proguard Build API
    Then the sourcemaps Content-Type header is valid multipart form-data
    And the sourcemap payload field "apiKey" equals "1234567890ABCDEF1234567890ABCDEF"
    And the sourcemap payload field "appId" equals "com.reactnative"
    And the sourcemap payload field "versionCode" equals "1"
    And the sourcemap payload field "versionName" equals "1.0"
    And the sourcemap payload field "overwrite" equals "true"
    And I discard the oldest sourcemap

    Then the sourcemap is valid for the NDK Build API
    Then the sourcemaps Content-Type header is valid multipart form-data
    And the sourcemap payload field "apiKey" equals "1234567890ABCDEF1234567890ABCDEF"
    And the sourcemap payload field "appId" equals "com.reactnative"
    And the sourcemap payload field "versionCode" equals "1"
    And the sourcemap payload field "versionName" equals "1.0"
    And the sourcemap payload field "overwrite" equals "true"

    Then the sourcemap is valid for the dSYM Build API
    Then the sourcemaps Content-Type header is valid multipart form-data
    And the sourcemap payload field "apiKey" equals "1234567890ABCDEF1234567890ABCDEF"
    And I discard the oldest sourcemap