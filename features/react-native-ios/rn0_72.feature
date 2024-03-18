Feature: React Native 0.72 iOS Integration Tests

  Scenario: Upload a single React Native 0.72 iOS sourcemap using all CLI flags
    When I run bugsnag-cli with upload react-native-ios --upload-api-root-url=http://localhost:9339 --api-key=1234567890ABCDEF1234567890ABCDEF --overwrite --dev --bundle=features/react-native-ios/fixtures/rn0_72/ios/build/sourcemaps/main.jsbundle --source-map=features/react-native-ios/fixtures/rn0_72/ios/build/sourcemaps/main.jsbundle.map --plist=features/react-native-ios/fixtures/rn0_72/ios/build/sourcemaps/Info.plist --xcode-project=features/react-native-ios/fixtures/rn0_72/ios/bugsnag_cli_test.xcworkspace --scheme=bugsnag_cli_test --code-bundle-id=1.0-15 --version-name=1.0 --bundle-version=1
    And I wait to receive 1 sourcemaps
    Then the sourcemap is valid for the React Native Build API
    Then the sourcemaps Content-Type header is valid multipart form-data
    And the sourcemap payload field "apiKey" equals "1234567890ABCDEF1234567890ABCDEF"
    And the sourcemap payload field "appVersion" equals "1.0"
    And the sourcemap payload field "appBundleVersion" equals "1"
    And the sourcemap payload field "codeBundleId" equals "1.0-15"
    And the sourcemap payload field "dev" equals "true"
    And the sourcemap payload field "platform" equals "ios"
    And the sourcemap payload field "overwrite" equals "true"

  Scenario: Upload a single React Native 0.72 iOS sourcemap providing the bundle, source-map and plist flags
    When I run bugsnag-cli with upload react-native-ios --upload-api-root-url=http://localhost:9339 --api-key=1234567890ABCDEF1234567890ABCDEF --overwrite --dev --bundle=features/react-native-ios/fixtures/rn0_72/ios/build/sourcemaps/main.jsbundle --source-map=features/react-native-ios/fixtures/rn0_72/ios/build/sourcemaps/main.jsbundle.map --plist=features/react-native-ios/fixtures/rn0_72/ios/build/sourcemaps/Info.plist
    And I wait to receive 1 sourcemaps
    Then the sourcemap is valid for the React Native Build API
    Then the sourcemaps Content-Type header is valid multipart form-data
    And the sourcemap payload field "apiKey" equals "1234567890ABCDEF1234567890ABCDEF"
    And the sourcemap payload field "appVersion" equals "1.0"
    And the sourcemap payload field "appBundleVersion" equals "1"
    And the sourcemap payload field "platform" equals "ios"
    And the sourcemap payload field "overwrite" equals "true"

  Scenario: Upload a single React Native 0.72 iOS sourcemap without a source-map flag
    When I run bugsnag-cli with upload react-native-ios --upload-api-root-url=http://localhost:9339 --api-key=1234567890ABCDEF1234567890ABCDEF --overwrite --dev --bundle=features/react-native-ios/fixtures/rn0_72/ios/build/sourcemaps/main.jsbundle --plist=features/react-native-ios/fixtures/rn0_72/ios/build/sourcemaps/Info.plist --xcode-project=features/react-native-ios/fixtures/rn0_72/ios/bugsnag_cli_test.xcworkspace --scheme=bugsnag_cli_test --code-bundle-id=1.0-15 --version-name=1.0 --bundle-version=1 features/react-native-ios/fixtures/rn0_72
    And I wait to receive 1 sourcemaps
    Then the sourcemap is valid for the React Native Build API
    Then the sourcemaps Content-Type header is valid multipart form-data
    And the sourcemap payload field "apiKey" equals "1234567890ABCDEF1234567890ABCDEF"
    And the sourcemap payload field "appVersion" equals "1.0"
    And the sourcemap payload field "appBundleVersion" equals "1"
    And the sourcemap payload field "platform" equals "ios"
    And the sourcemap payload field "overwrite" equals "true"

  Scenario: Build and Upload React Native 0.72 iOS sourcemaps
    When I make the "features/base-fixtures/rn0_72/ios"
    And I wait for the build to succeed

    When I run bugsnag-cli with upload react-native-ios --upload-api-root-url=http://localhost:9339 --api-key=1234567890ABCDEF1234567890ABCDEF --overwrite --dev --scheme=rn0_72 features/base-fixtures/rn0_72
    And I wait to receive 1 sourcemaps
    Then the sourcemap is valid for the React Native Build API
    Then the sourcemaps Content-Type header is valid multipart form-data
    And the sourcemap payload field "apiKey" equals "1234567890ABCDEF1234567890ABCDEF"
    And the sourcemap payload field "appVersion" equals "1.0"
    And the sourcemap payload field "appBundleVersion" equals "1"
    And the sourcemap payload field "platform" equals "ios"
    And the sourcemap payload field "overwrite" equals "true"

  Scenario: Build, Archive and Upload React Native 0.72 iOS sourcemaps
    When I make the "features/base-fixtures/rn0_72/ios/archive"
    And I wait for the build to succeed

    When I run bugsnag-cli with upload react-native-ios --upload-api-root-url=http://localhost:9339 --api-key=1234567890ABCDEF1234567890ABCDEF --overwrite --dev --scheme=rn0_72 features/base-fixtures/rn0_72
    And I wait to receive 1 sourcemaps
    Then the sourcemap is valid for the React Native Build API
    Then the sourcemaps Content-Type header is valid multipart form-data
    And the sourcemap payload field "apiKey" equals "1234567890ABCDEF1234567890ABCDEF"
    And the sourcemap payload field "appVersion" equals "1.0"
    And the sourcemap payload field "appBundleVersion" equals "1"
    And the sourcemap payload field "platform" equals "ios"
    And the sourcemap payload field "overwrite" equals "true"
