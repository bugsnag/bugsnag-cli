Feature: Dart Integration Tests

  Scenario: Upload a single Android Dart sourcemap using all CLI flags
    When I run bugsnag-cli with upload dart --upload-api-root-url=http://localhost:$MAZE_RUNNER_PORT --api-key=1234567890ABCDEF1234567890ABCDEF --overwrite --version-name=2.0 --version-code=1.0 features/dart/fixtures/app-debug-info/app.android-arm64.symbols
    And I wait to receive 1 sourcemaps
    Then the sourcemap is valid for the Dart Build API
    Then the sourcemaps Content-Type header is valid multipart form-data
    And the sourcemap payload field "apiKey" equals "1234567890ABCDEF1234567890ABCDEF"
    And the sourcemap payload field "appVersion" equals "2.0"
    And the sourcemap payload field "appVersionCode" equals "1.0"
    And the sourcemap payload field "buildId" equals "07cc131ca803c124e93268ce19322737"
    And the sourcemap payload field "overwrite" equals "true"

  Scenario: Upload a single iOS Dart sourcemap using all CLI flags
    When I run bugsnag-cli with upload dart --upload-api-root-url=http://localhost:$MAZE_RUNNER_PORT --api-key=1234567890ABCDEF1234567890ABCDEF --overwrite --ios-app-path=features/dart/fixtures/build/ios/iphoneos/Runner.app/Frameworks/App.framework/App --version-name=2.0 --bundle-version=1.0 features/dart/fixtures/app-debug-info/app.ios-arm64.symbols
    And I wait to receive 1 sourcemaps
    Then the sourcemap is valid for the Dart Build API
    Then the sourcemaps Content-Type header is valid multipart form-data
    And the sourcemap payload field "apiKey" equals "1234567890ABCDEF1234567890ABCDEF"
    And the sourcemap payload field "appVersion" equals "2.0"
    And the sourcemap payload field "appBundleVersion" equals "1.0"
    And the sourcemap payload field "buildId" equals "E30C1BE5-DEB6-373C-98B4-52D827B7FF0D"
    And the sourcemap payload field "overwrite" equals "true"

  Scenario: Upload multiple Dart sourcemaps providing no flags to the CLI
    When I run bugsnag-cli with upload dart --upload-api-root-url=http://localhost:$MAZE_RUNNER_PORT --api-key=1234567890ABCDEF1234567890ABCDEF --overwrite features/dart/fixtures/app-debug-info
    And I wait to receive 4 sourcemaps
    Then the sourcemap is valid for the Dart Build API
    Then the sourcemaps Content-Type header is valid multipart form-data
    And the sourcemap payload field "apiKey" equals "1234567890ABCDEF1234567890ABCDEF"
    And the sourcemap payload field "buildId" equals "1dacae8b2a346d3dc6271a742f5d5210"
    And the sourcemap payload field "overwrite" equals "true"
