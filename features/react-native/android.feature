# The following are set as part of the BuildRNAndroid step:
#  - APP_MANIFEST_PATH
#  - BUNDLE_PATH
#  - SOURCE_MAP_PATH

@BuildRNAndroid
Feature: React Native Android Integration Tests
  Scenario: Upload a single React Native Android sourcemap using all escape hatches
    When I run bugsnag-cli with upload react-native-android --upload-api-root-url=http://localhost:$MAZE_RUNNER_PORT --bundle=$BUNDLE_PATH --source-map=$SOURCE_MAP_PATH --app-manifest=$APP_MANIFEST_PATH --version-name=2.0 --variant=release --version-code=2 --overwrite
    And I wait to receive 1 sourcemaps
    Then the sourcemaps are valid for the API
    Then the sourcemaps Content-Type header is valid multipart form-data
    And the sourcemap payload field "apiKey" equals "1234567890ABCDEF1234567890ABCDEF"
    And the sourcemap payload field "appVersion" equals "2.0"
    And the sourcemap payload field "appVersionCode" equals "2"
    And the sourcemap payload field "platform" equals "android"
    And the sourcemap payload field "overwrite" equals "true"

  Scenario: Upload a single React Native Android sourcemap, only passing bundle, sourcemap and manifest options
    When I run bugsnag-cli with upload react-native-android --upload-api-root-url=http://localhost:$MAZE_RUNNER_PORT --bundle=$BUNDLE_PATH --source-map=$SOURCE_MAP_PATH --app-manifest=$APP_MANIFEST_PATH --overwrite features/react-native/fixtures/generated/old-arch/$RN_VERSION
    And I wait to receive 1 sourcemaps
    Then the sourcemaps are valid for the API
    Then the sourcemaps Content-Type header is valid multipart form-data
    And the sourcemap payload field "apiKey" equals "1234567890ABCDEF1234567890ABCDEF"
    And the sourcemap payload field "appVersion" equals "1.0"
    And the sourcemap payload field "appVersionCode" equals "1"
    And the sourcemap payload field "platform" equals "android"
    And the sourcemap payload field "overwrite" equals "true"

  Scenario: Upload a single React Native Android sourcemap
    When I run bugsnag-cli with upload react-native-android --upload-api-root-url=http://localhost:$MAZE_RUNNER_PORT --overwrite features/react-native/fixtures/generated/old-arch/$RN_VERSION
    And I wait to receive 1 sourcemaps
    Then the sourcemaps are valid for the API
    Then the sourcemaps Content-Type header is valid multipart form-data
    And the sourcemap payload field "apiKey" equals "1234567890ABCDEF1234567890ABCDEF"
    And the sourcemap payload field "appVersion" equals "1.0"
    And the sourcemap payload field "appVersionCode" equals "1"
    And the sourcemap payload field "platform" equals "android"
    And the sourcemap payload field "overwrite" equals "true"
