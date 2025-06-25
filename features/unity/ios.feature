Feature: Unity iOS integration tests
  Scenario: Unity iOS integration test
    Given I build the Unity project for iOS
    And I wait for the build to succeed

    When I run bugsnag-cli with upload unity-ios --upload-api-root-url=http://localhost:$MAZE_RUNNER_PORT --api-key=1234567890ABCDEF1234567890ABCDEF --overwrite --no-upload-il2cpp-mapping platforms-examples/Unity/UnityExample/
    Then I wait to receive 2 sourcemaps
    Then the sourcemap is valid for the dSYM Build API
    Then the sourcemaps Content-Type header is valid multipart form-data
    And the sourcemap payload field "apiKey" equals "1234567890ABCDEF1234567890ABCDEF"

  Scenario: Unity iOS integration test with Unity Line Mappings
    Given I build the Unity project for iOS
    And I wait for the build to succeed

    When I run bugsnag-cli with upload unity-ios --upload-api-root-url=http://localhost:$MAZE_RUNNER_PORT --api-key=1234567890ABCDEF1234567890ABCDEF --overwrite platforms-examples/Unity/UnityExample/
    Then I wait to receive 3 sourcemaps

    Then the sourcemap is valid for the dSYM Build API
    Then the sourcemaps Content-Type header is valid multipart form-data
    And the sourcemap payload field "apiKey" equals "1234567890ABCDEF1234567890ABCDEF"

    And I discard the oldest sourcemaps

    Then the sourcemap is valid for the dSYM Build API
    Then the sourcemaps Content-Type header is valid multipart form-data
    And the sourcemap payload field "apiKey" equals "1234567890ABCDEF1234567890ABCDEF"

    And I discard the oldest sourcemaps

    Then the sourcemap is valid for the Unity Line Mapping API
    Then the sourcemaps Content-Type header is valid multipart form-data
    And the sourcemap payload field "appVersion" equals "1.0"
    And the sourcemap payload field "appBundleVersion" equals "1.0"
    And the sourcemap payload field "dsymUUID" is not null
    And the sourcemap payload field "appId" equals "com.apple.xcode.dsym.com.unity3d.framework"
    And the sourcemap payload field "overwrite" equals "true"


  Scenario: Unity iOS integration test with Unity Line Mappings passing version numbers
    Given I build the Unity project for iOS
    And I wait for the build to succeed

    When I run bugsnag-cli with upload unity-ios --upload-api-root-url=http://localhost:$MAZE_RUNNER_PORT --api-key=1234567890ABCDEF1234567890ABCDEF --overwrite platforms-examples/Unity/UnityExample/ --application-id=com.bugsnag.unity.test --bundle-version=999.99 --version-name=123.456
    Then I wait to receive 3 sourcemaps

    Then the sourcemap is valid for the dSYM Build API
    Then the sourcemaps Content-Type header is valid multipart form-data
    And the sourcemap payload field "apiKey" equals "1234567890ABCDEF1234567890ABCDEF"

    And I discard the oldest sourcemaps

    Then the sourcemap is valid for the dSYM Build API
    Then the sourcemaps Content-Type header is valid multipart form-data
    And the sourcemap payload field "apiKey" equals "1234567890ABCDEF1234567890ABCDEF"

    And I discard the oldest sourcemaps

    Then the sourcemap is valid for the Unity Line Mapping API
    Then the sourcemaps Content-Type header is valid multipart form-data
    And the sourcemap payload field "appVersion" equals "123.456"
    And the sourcemap payload field "appBundleVersion" equals "999.99"
    And the sourcemap payload field "dsymUUID" is not null
    And the sourcemap payload field "appId" equals "com.bugsnag.unity.test"
    And the sourcemap payload field "overwrite" equals "true"
