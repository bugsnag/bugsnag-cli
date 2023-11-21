Feature: React Native 0.72 Cocoa Integration Tests

  Scenario: Upload a single React Native 0.72 Cocoa sourcemap using all CLI flags
    When I run bugsnag-cli with upload react-native-cocoa --api-key=1234567890ABCDEF1234567890ABCDEF --overwrite --bundle-version=1.0-15 --dev --bundle=features/react-native-cocoa/fixtures/rn0_72/ios/build/sourcemaps/main.jsbundle --source-map=features/react-native-cocoa/fixtures/rn0_72/ios/build/sourcemaps/main.jsbundle.map --version-name=1.0 features/react-native-cocoa/fixtures/rn0_72
    And I wait to receive 1 sourcemaps
    Then the sourcemap is valid for the React Native Build API
    Then the sourcemaps Content-Type header is valid multipart form-data
    And the sourcemap payload field "apiKey" equals "1234567890ABCDEF1234567890ABCDEF"
    And the sourcemap payload field "appVersion" equals "1.0"
    And the sourcemap payload field "appVersionCode" equals "1"
    And the sourcemap payload field "codeBundleId" equals "1.0-15"
    And the sourcemap payload field "dev" equals "true"
    And the sourcemap payload field "platform" equals "ios"
    And the sourcemap payload field "overwrite" equals "true"

  Scenario: Upload a single React Native 0.72 Cocoa sourcemap providing the bundle, source-map and plist flags
    When I run bugsnag-cli with upload react-native-cocoa --upload-api-root-url=http://localhost:9339 --api-key=1234567890ABCDEF1234567890ABCDEF --overwrite --bundle=features/react-native-cocoa/fixtures/rn0_72/ios/build/sourcemaps/main.jsbundle --source-map=features/react-native-cocoa/fixtures/rn0_72/ios/build/sourcemaps/main.jsbundle.map --plist=features/react-native-cocoa/fixtures/rn0_72/ios/bugsnag_cli_test/Info.plist features/react-native-cocoa/fixtures/rn0_72
    And I wait to receive 1 sourcemaps
    Then the sourcemap is valid for the React Native Build API
    Then the sourcemaps Content-Type header is valid multipart form-data
    And the sourcemap payload field "apiKey" equals "1234567890ABCDEF1234567890ABCDEF"
    And the sourcemap payload field "appVersion" equals "1.0"
    And the sourcemap payload field "appVersionCode" equals "1"
    And the sourcemap payload field "platform" equals "ios"
    And the sourcemap payload field "overwrite" equals "true"

  Scenario: Upload a single React Native 0.72 Cocoa sourcemap providing the bundle CLI flag
    When I run bugsnag-cli with upload react-native-cocoa --upload-api-root-url=http://localhost:9339 --api-key=1234567890ABCDEF1234567890ABCDEF --overwrite --bundle=features/react-native-cocoa/fixtures/rn0_72/ios/build/sourcemaps/main.jsbundle features/react-native-cocoa/fixtures/rn0_72
    And I wait to receive 1 sourcemaps
    Then the sourcemap is valid for the React Native Build API
    Then the sourcemaps Content-Type header is valid multipart form-data
    And the sourcemap payload field "apiKey" equals "1234567890ABCDEF1234567890ABCDEF"
    And the sourcemap payload field "appVersion" equals "1.0"
    And the sourcemap payload field "appVersionCode" equals "1"
    And the sourcemap payload field "platform" equals "android"
    And the sourcemap payload field "overwrite" equals "true"

  Scenario: Upload a single React Native 0.72 Cocoa sourcemap providing no CLI flags
    When I run bugsnag-cli with upload react-native-cocoa --upload-api-root-url=http://localhost:9339 --api-key=1234567890ABCDEF1234567890ABCDEF --overwrite features/react-native-cocoa/fixtures/rn0_72
    And I wait to receive 1 sourcemaps
    Then the sourcemap is valid for the React Native Build API
    Then the sourcemaps Content-Type header is valid multipart form-data
    And the sourcemap payload field "apiKey" equals "1234567890ABCDEF1234567890ABCDEF"
    And the sourcemap payload field "appVersion" equals "1.0"
    And the sourcemap payload field "appVersionCode" equals "1"
    And the sourcemap payload field "platform" equals "ios"
    And the sourcemap payload field "overwrite" equals "true"

  Scenario: Build and Upload React Native 0.72 Cocoa sourcemaps
    When I make the "features/base-fixtures/rn0_72"
    And I wait for the build to succeed

    When I run bugsnag-cli with upload react-native-cocoa --upload-api-root-url=http://localhost:9339 --api-key=1234567890ABCDEF1234567890ABCDEF --overwrite features/base-fixtures/rn0_72
    And I wait to receive 1 sourcemaps
    Then the sourcemap is valid for the React Native Build API
    Then the sourcemaps Content-Type header is valid multipart form-data
    And the sourcemap payload field "apiKey" equals "1234567890ABCDEF1234567890ABCDEF"
    And the sourcemap payload field "appVersion" equals "1.0"
    And the sourcemap payload field "appVersionCode" equals "1"
    And the sourcemap payload field "platform" equals "android"
    And the sourcemap payload field "overwrite" equals "true"