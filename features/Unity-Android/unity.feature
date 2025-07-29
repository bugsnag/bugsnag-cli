Feature: Unity Android integration tests
  Scenario: Unity Android integration tests
    Given I build the Unity Android example project
    And I wait for the Unity symbols to generate

    When I run bugsnag-cli with upload unity-android --upload-api-root-url=http://localhost:$MAZE_RUNNER_PORT --api-key=1234567890ABCDEF1234567890ABCDEF --overwrite platforms-examples/Unity/
    Then I wait to receive 5 sourcemaps
    Then the sourcemaps are valid for the API
    Then the sourcemaps Content-Type header is valid multipart form-data
    Then the sourcemap payload fields should be:
      | apiKey       | 1234567890ABCDEF1234567890ABCDEF     |
      | appId        | com.bugsnag.example.unity.android    |
      | versionCode  | 1                                    |
      | versionName  | 1.0                                  |
      | overwrite    | true                                 |

  Scenario: Unity Android integration tests using the bundled NDK
    Given I build the Unity Android example project
    And I wait for the Unity symbols to generate

    Given I set the NDK path to the Unity bundled version
    When I run bugsnag-cli with upload unity-android --upload-api-root-url=http://localhost:$MAZE_RUNNER_PORT --api-key=1234567890ABCDEF1234567890ABCDEF --overwrite platforms-examples/Unity/
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
    Given I build the Unity Android example project
    And I wait for the Unity symbols to generate

    When I run bugsnag-cli with upload unity-android --upload-api-root-url=http://localhost:$MAZE_RUNNER_PORT --api-key=1234567890ABCDEF1234567890ABCDEF --overwrite  --aab-path platforms-examples/Unity/UnityExample.aab platforms-examples/Unity/
    Then I wait to receive 5 sourcemaps
    Then the sourcemaps are valid for the API
    Then the sourcemaps Content-Type header is valid multipart form-data
    Then the sourcemap payload fields should be:
      | apiKey       | 1234567890ABCDEF1234567890ABCDEF     |
      | appId        | com.bugsnag.example.unity.android    |
      | versionCode  | 1                                    |
      | versionName  | 1.0                                  |
      | overwrite    | true                                 |
