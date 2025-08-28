Feature: Unity Android integration tests
  Scenario: Unity Android integration tests
    Given I build the Unity project for Android
    And I wait for the Unity symbols to generate

    When I run bugsnag-cli with upload unity-android --upload-api-root-url=http://localhost:$MAZE_RUNNER_PORT --api-key=1234567890ABCDEF1234567890ABCDEF --overwrite --no-upload-il2cpp-mapping platforms-examples/Unity/
    Then I wait to receive 5 sourcemaps
    Then the sourcemaps are valid for the API
    Then the sourcemaps Content-Type header is valid multipart form-data
    Then the sourcemap payload fields should be:
      | apiKey       | 1234567890ABCDEF1234567890ABCDEF     |
      | appId        | com.bugsnag.example.unity.android    |
      | versionCode  | 1                                    |
      | versionName  | 1.0                                  |
      | overwrite    | true                                 |

  Scenario: Unity Android integration tests passing the aab file
    Given I build the Unity project for Android
    And I wait for the Unity symbols to generate

    When I run bugsnag-cli with upload unity-android --upload-api-root-url=http://localhost:$MAZE_RUNNER_PORT --api-key=1234567890ABCDEF1234567890ABCDEF --overwrite --no-upload-il2cpp-mapping --aab-path platforms-examples/Unity/UnityExample.aab platforms-examples/Unity/
    Then I wait to receive 5 sourcemaps
    Then the sourcemaps are valid for the API
    Then the sourcemaps Content-Type header is valid multipart form-data
    Then the sourcemap payload fields should be:
      | apiKey       | 1234567890ABCDEF1234567890ABCDEF     |
      | appId        | com.bugsnag.example.unity.android    |
      | versionCode  | 1                                    |
      | versionName  | 1.0                                  |
      | overwrite    | true                                 |

  Scenario: Unity Android integration test with Unity Line Mappings
    Given I build the Unity project for Android
    And I wait for the Unity symbols to generate

    When I run bugsnag-cli with upload unity-android --upload-api-root-url=http://localhost:$MAZE_RUNNER_PORT --api-key=1234567890ABCDEF1234567890ABCDEF --overwrite platforms-examples/Unity/
    Then I wait to receive 6 sourcemaps

    And I sort the sourcemaps by path

    Then the sourcemap is valid for the Android Build API
    Then the sourcemaps Content-Type header is valid multipart form-data
    And the sourcemap payload field "appId" equals "com.bugsnag.example.unity.android"
    And the sourcemap payload field "versionCode" equals "1"
    And the sourcemap payload field "versionName" equals "1.0"
    And the sourcemap payload field "overwrite" equals "true"

    And I discard the oldest sourcemaps

    Then the sourcemap is valid for the Android Build API
    Then the sourcemaps Content-Type header is valid multipart form-data
    And the sourcemap payload field "appId" equals "com.bugsnag.example.unity.android"
    And the sourcemap payload field "versionCode" equals "1"
    And the sourcemap payload field "versionName" equals "1.0"
    And the sourcemap payload field "overwrite" equals "true"

    And I discard the oldest sourcemaps

    Then the sourcemap is valid for the Android Build API
    Then the sourcemaps Content-Type header is valid multipart form-data
    And the sourcemap payload field "appId" equals "com.bugsnag.example.unity.android"
    And the sourcemap payload field "versionCode" equals "1"
    And the sourcemap payload field "versionName" equals "1.0"
    And the sourcemap payload field "overwrite" equals "true"

    And I discard the oldest sourcemaps

    Then the sourcemap is valid for the Android Build API
    Then the sourcemaps Content-Type header is valid multipart form-data
    And the sourcemap payload field "appId" equals "com.bugsnag.example.unity.android"
    And the sourcemap payload field "versionCode" equals "1"
    And the sourcemap payload field "versionName" equals "1.0"
    And the sourcemap payload field "overwrite" equals "true"

    And I discard the oldest sourcemaps

    Then the sourcemap is valid for the Android Build API
    Then the sourcemaps Content-Type header is valid multipart form-data
    And the sourcemap payload field "appId" equals "com.bugsnag.example.unity.android"
    And the sourcemap payload field "versionCode" equals "1"
    And the sourcemap payload field "versionName" equals "1.0"
    And the sourcemap payload field "overwrite" equals "true"

    And I discard the oldest sourcemaps

    Then the sourcemap is valid for the Unity Line Mapping API
    Then the sourcemaps Content-Type header is valid multipart form-data
    And the sourcemap payload field "apiKey" equals "1234567890ABCDEF1234567890ABCDEF"
    And the sourcemap payload field "appVersion" equals "1.0"
    And the sourcemap payload field "appVersionCode" equals "1"
    And the sourcemap payload field "soBuildId" is not null
    And the sourcemap payload field "appId" equals "com.bugsnag.example.unity.android"
    And the sourcemap payload field "overwrite" equals "true"

  Scenario: Unity Android integration test with Unity Line Mappings passing version numbers
    Given I build the Unity project for Android
    And I wait for the build to succeed

    When I run bugsnag-cli with upload unity-android --upload-api-root-url=http://localhost:$MAZE_RUNNER_PORT --api-key=1234567890ABCDEF1234567890ABCDEF --overwrite platforms-examples/Unity/ --application-id=com.bugsnag.unity.test --version-code=999.99 --version-name=123.456
    Then I wait to receive 6 sourcemaps

    And I sort the sourcemaps by path

    Then the sourcemap is valid for the Android Build API
    Then the sourcemaps Content-Type header is valid multipart form-data
    And the sourcemap payload field "appId" equals "com.bugsnag.unity.test"
    And the sourcemap payload field "versionCode" equals "999.99"
    And the sourcemap payload field "versionName" equals "123.456"
    And the sourcemap payload field "overwrite" equals "true"

    And I discard the oldest sourcemaps

    Then the sourcemap is valid for the Android Build API
    Then the sourcemaps Content-Type header is valid multipart form-data
    And the sourcemap payload field "appId" equals "com.bugsnag.unity.test"
    And the sourcemap payload field "versionCode" equals "999.99"
    And the sourcemap payload field "versionName" equals "123.456"
    And the sourcemap payload field "overwrite" equals "true"

    And I discard the oldest sourcemaps

    Then the sourcemap is valid for the Android Build API
    Then the sourcemaps Content-Type header is valid multipart form-data
    And the sourcemap payload field "appId" equals "com.bugsnag.unity.test"
    And the sourcemap payload field "versionCode" equals "999.99"
    And the sourcemap payload field "versionName" equals "123.456"
    And the sourcemap payload field "overwrite" equals "true"

    And I discard the oldest sourcemaps

    Then the sourcemap is valid for the Android Build API
    Then the sourcemaps Content-Type header is valid multipart form-data
    And the sourcemap payload field "appId" equals "com.bugsnag.unity.test"
    And the sourcemap payload field "versionCode" equals "999.99"
    And the sourcemap payload field "versionName" equals "123.456"
    And the sourcemap payload field "overwrite" equals "true"

    And I discard the oldest sourcemaps

    Then the sourcemap is valid for the Android Build API
    Then the sourcemaps Content-Type header is valid multipart form-data
    And the sourcemap payload field "appId" equals "com.bugsnag.unity.test"
    And the sourcemap payload field "versionCode" equals "999.99"
    And the sourcemap payload field "versionName" equals "123.456"
    And the sourcemap payload field "overwrite" equals "true"

    And I discard the oldest sourcemaps

    Then the sourcemap is valid for the Unity Line Mapping API
    Then the sourcemaps Content-Type header is valid multipart form-data
    And the sourcemap payload field "apiKey" equals "1234567890ABCDEF1234567890ABCDEF"
    And the sourcemap payload field "appVersion" equals "123.456"
    And the sourcemap payload field "appVersionCode" equals "999.99"
    And the sourcemap payload field "soBuildId" is not null
    And the sourcemap payload field "appId" equals "com.bugsnag.unity.test"
    And the sourcemap payload field "overwrite" equals "true"

