Feature: React Native 0.72 Cocoa Integration Tests

  Scenario: Upload a single React Native 0.72 Cocoa sourcemap using all CLI flags
    When I run bugsnag-cli with upload upload react-native-cocoa --upload-api-root-url=http://localhost:9339 --api-key=1234567890ABCDEF1234567890ABCDEF --bundle-version=1.0-15 --scheme=bugsnag_cli_test --dev --bundle=features/react-native-cocoa/fixtures/rn0_72/ios/build/sourcemaps/main.jsbundle --source-map=features/react-native-cocoa/fixtures/rn0_72/ios/build/sourcemaps/main.jsbundle.map --plist=features/react-native-cocoa/fixtures/rn0_72/ios/bugsnag_cli_test/Info.plist --version-name=1.0 features/react-native-cocoa/fixtures/rn0_72
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
