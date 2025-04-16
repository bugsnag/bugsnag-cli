Feature: React Native iOS Integration Tests
  @BuildRNiOS
  Scenario: Upload a single React Native iOS sourcemap
    When I run bugsnag-cli with upload react-native-ios --upload-api-root-url=http://localhost:$MAZE_RUNNER_PORT --api-key=1234567890ABCDEF1234567890ABCDEF --overwrite features/react-native/fixtures/generated/old-arch/$RN_VERSION
    And I wait to receive 1 sourcemaps
    Then the sourcemap is valid for the React Native Build API
    Then the sourcemaps Content-Type header is valid multipart form-data
    And the sourcemap payload field "apiKey" equals "1234567890ABCDEF1234567890ABCDEF"
    And the sourcemap payload field "appVersion" equals "1.0"
    And the sourcemap payload field "appBundleVersion" equals "1"
    And the sourcemap payload field "platform" equals "ios"
    And the sourcemap payload field "overwrite" equals "true"

  @BuildExportRNiOS
  Scenario: Upload a single React Native iOS sourcemap using escape hatches
    When I run bugsnag-cli with upload react-native-ios --upload-api-root-url=http://localhost:$MAZE_RUNNER_PORT --api-key=1234567890ABCDEF1234567890ABCDEF --overwrite --bundle=features/react-native/fixtures/generated/old-arch/$RN_VERSION/reactnative.xcarchive/Products/Applications/reactnative.app/main.jsbundle --source-map=features/react-native/fixtures/generated/old-arch/$RN_VERSION/ios/build/sourcemaps/main.jsbundle.map --plist=features/react-native/fixtures/generated/old-arch/$RN_VERSION/reactnative.xcarchive/Products/Applications/reactnative.app/Info.plist --xcode-project=features/react-native/fixtures/generated/old-arch/$RN_VERSION/ios/reactnative.xcodeproj --scheme=reactnative  --bundle-version=2 --version-name=2.0
    Then the sourcemap is valid for the React Native Build API
    Then the sourcemaps Content-Type header is valid multipart form-data
    And the sourcemap payload field "apiKey" equals "1234567890ABCDEF1234567890ABCDEF"
    And the sourcemap payload field "appVersion" equals "2.0"
    And the sourcemap payload field "appBundleVersion" equals "2"
    And the sourcemap payload field "platform" equals "ios"
    And the sourcemap payload field "overwrite" equals "true"

  Scenario: Upload a single React Native iOS sourcemap using escape hatches
    When I run bugsnag-cli with upload react-native-ios --upload-api-root-url=http://localhost:$MAZE_RUNNER_PORT --api-key=1234567890ABCDEF1234567890ABCDEF --overwrite --bundle=features/react-native/fixtures/generated/old-arch/$RN_VERSION/reactnative.xcarchive/Products/Applications/reactnative.app/main.jsbundle --source-map=features/react-native/fixtures/generated/old-arch/$RN_VERSION/ios/build/sourcemaps/main.jsbundle.map --plist=features/react-native/fixtures/generated/old-arch/$RN_VERSION/reactnative.xcarchive/Products/Applications/reactnative.app/Info.plist --xcode-project=features/react-native/fixtures/generated/old-arch/$RN_VERSION/ios/reactnative.xcodeproj --scheme=reactnative
    Then the sourcemap is valid for the React Native Build API
    Then the sourcemaps Content-Type header is valid multipart form-data
    And the sourcemap payload field "apiKey" equals "1234567890ABCDEF1234567890ABCDEF"
    And the sourcemap payload field "appVersion" equals "1.0"
    And the sourcemap payload field "appBundleVersion" equals "1"
    And the sourcemap payload field "platform" equals "ios"
    And the sourcemap payload field "overwrite" equals "true"
