Feature: React Native Integration Tests
  Scenario: Upload a single React Native sourcemap
    When I run bugsnag-cli with upload react-native-sourcemaps --upload-api-root-url=http://localhost:$MAZE_RUNNER_PORT --api-key=1234567890ABCDEF1234567890ABCDEF --version-name=1.0 --version-code=123.456 --source-map=features/react-native/fixtures/vega/index.bundle.map --bundle=features/react-native/fixtures/vega/index.bundle --platform=vega
    And I wait to receive 1 sourcemaps
    Then the sourcemaps are valid for the API
    Then the sourcemaps Content-Type header is valid multipart form-data
    Then the sourcemap payload fields should be:
      | apiKey            | 1234567890ABCDEF1234567890ABCDEF  |
      | appVersion        | 1.0                               |
      | appVersionCode    | 123.456                           |
      | platform          | vega                              |

  Scenario: Upload a single React Native sourcemap providing both version-code and bundle-version
    When I run bugsnag-cli with upload react-native-sourcemaps --upload-api-root-url=http://localhost:$MAZE_RUNNER_PORT --api-key=1234567890ABCDEF1234567890ABCDEF --version-name=1.0 --version-code=123.456  --bundle-version=123.456 --source-map=features/react-native/fixtures/vega/index.bundle.map --bundle=features/react-native/fixtures/vega/index.bundle --platform=vega
    Then the error should contain "--version-code and --bundle-version can't be used together"

  Scenario: Upload a single React Native sourcemap providing code-bundle-id
    When I run bugsnag-cli with upload react-native-sourcemaps --upload-api-root-url=http://localhost:$MAZE_RUNNER_PORT --api-key=1234567890ABCDEF1234567890ABCDEF --code-bundle-id=123.456 --source-map=features/react-native/fixtures/vega/index.bundle.map --bundle=features/react-native/fixtures/vega/index.bundle --platform=vega
    And I wait to receive 1 sourcemaps
    Then the sourcemaps are valid for the API
    Then the sourcemaps Content-Type header is valid multipart form-data
    Then the sourcemap payload fields should be:
      | apiKey            | 1234567890ABCDEF1234567890ABCDEF  |
      | codeBundleId      | 123.456                           |
      | platform          | vega                              |
