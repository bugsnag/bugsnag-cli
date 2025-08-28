Feature: Dart Integration Tests

  Scenario: Upload a single Android Dart sourcemap using all CLI flags
    When I run bugsnag-cli upload "dart" with the following arguments:
      | --upload-api-root-url                                           | http://localhost:$MAZE_RUNNER_PORT |
      | --api-key                                                       | 1234567890ABCDEF1234567890ABCDEF   |
      | --version-name                                                  | 2.0                                |
      | --version-code                                                  | 1.0                                |
      | features/dart/fixtures/app-debug-info/app.android-arm64.symbols |                                    |
    And I wait to receive 1 sourcemaps
    Then the sourcemaps are valid for the API
    Then the sourcemaps Content-Type header is valid multipart form-data
    Then the sourcemap payload fields should be:
      | apiKey          | 1234567890ABCDEF1234567890ABCDEF     |
      | appVersion      | 2.0                                  |
      | appVersionCode  | 1.0                                  |
      | buildId         | 07cc131ca803c124e93268ce19322737     |
      | overwrite       | true                                 |

  Scenario: Upload a single iOS Dart sourcemap using all CLI flags
    When I run bugsnag-cli upload "dart" with the following arguments:
      | --upload-api-root-url                                           | http://localhost:$MAZE_RUNNER_PORT                                                |
      | --api-key                                                       | 1234567890ABCDEF1234567890ABCDEF                                                  |
      | --version-name                                                  | 2.0                                                                               |
      | --version-code                                                  | 1.0                                                                               |
      | --ios-app-path                                                  | features/dart/fixtures/build/ios/iphoneos/Runner.app/Frameworks/App.framework/App |
      | features/dart/fixtures/app-debug-info/app.android-arm64.symbols |                                                                                   |
    And I wait to receive 1 sourcemaps
    Then the sourcemaps are valid for the API
    Then the sourcemaps Content-Type header is valid multipart form-data
    Then the sourcemap payload fields should be:
      | apiKey            | 1234567890ABCDEF1234567890ABCDEF      |
      | appVersion        | 2.0                                   |
      | appBundleVersion  | 1.0                                   |
      | buildId           | E30C1BE5-DEB6-373C-98B4-52D827B7FF0D  |
      | overwrite         | true                                  |


  Scenario: Upload multiple Dart sourcemaps providing no flags to the CLI
    When I run bugsnag-cli upload "dart" with the following arguments:
      | --upload-api-root-url                 | http://localhost:$MAZE_RUNNER_PORT |
      | --api-key                             | 1234567890ABCDEF1234567890ABCDEF   |
      | features/dart/fixtures/app-debug-info |                                    |
    And I wait to receive 4 sourcemaps
    Then the sourcemaps are valid for the API
    Then the sourcemaps Content-Type header is valid multipart form-data

    And the sourcemap payload field "apiKey" equals "1234567890ABCDEF1234567890ABCDEF"
    And the sourcemap payload field "buildId" equals "1dacae8b2a346d3dc6271a742f5d5210"
    And the sourcemap payload field "overwrite" equals "true"

    And I discard the oldest sourcemap

    And the sourcemap payload field "apiKey" equals "1234567890ABCDEF1234567890ABCDEF"
    And the sourcemap payload field "buildId" equals "07cc131ca803c124e93268ce19322737"
    And the sourcemap payload field "overwrite" equals "true"

    And I discard the oldest sourcemap

    And the sourcemap payload field "apiKey" equals "1234567890ABCDEF1234567890ABCDEF"
    And the sourcemap payload field "buildId" equals "c6e70c11ee73f14347202175430ab226"
    And the sourcemap payload field "overwrite" equals "true"

    And I discard the oldest sourcemap

    And the sourcemap payload field "apiKey" equals "1234567890ABCDEF1234567890ABCDEF"
    And the sourcemap payload field "buildId" equals "E30C1BE5-DEB6-373C-98B4-52D827B7FF0D"
    And the sourcemap payload field "overwrite" equals "true"
