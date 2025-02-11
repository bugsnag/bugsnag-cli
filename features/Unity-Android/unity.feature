Feature: Unity Android integration tests
  Scenario: Unity Android integration tests
    Given I build the Unity Android example project
    And I wait for the Unity symbols to generate

    When I run bugsnag-cli with upload unity-android --upload-api-root-url=http://localhost:9339 --api-key=1234567890ABCDEF1234567890ABCDEF --overwrite platforms-examples/Unity/
    Then I wait to receive 5 sourcemaps
    Then the sourcemap is valid for the Android Build API
    Then the sourcemaps Content-Type header is valid multipart form-data
    And the sourcemap payload field "apiKey" equals "1234567890ABCDEF1234567890ABCDEF"
    And the sourcemap payload field "appId" equals "com.bugsnag.example.unity.android"
    And the sourcemap payload field "versionCode" equals "1"
    And the sourcemap payload field "versionName" equals "1.0"
    And the sourcemap payload field "overwrite" equals "true"

  Scenario: Unity Android integration tests using the bundled NDK
    Given I build the Unity Android example project
    And I wait for the Unity symbols to generate

    Given I set the NDK path to the Unity bundled version
    When I run bugsnag-cli with upload unity-android --upload-api-root-url=http://localhost:9339 --api-key=1234567890ABCDEF1234567890ABCDEF --overwrite platforms-examples/Unity/
    Then I wait to receive 5 sourcemaps
    Then the sourcemap is valid for the Android Build API
    Then the sourcemaps Content-Type header is valid multipart form-data
    And the sourcemap payload field "apiKey" equals "1234567890ABCDEF1234567890ABCDEF"
    And the sourcemap payload field "appId" equals "com.bugsnag.example.unity.android"
    And the sourcemap payload field "versionCode" equals "1"
    And the sourcemap payload field "versionName" equals "1.0"
    And the sourcemap payload field "overwrite" equals "true"


  Scenario: Unity Android integration tests passing the aab file
    Given I build the Unity Android example project
    And I wait for the Unity symbols to generate

    When I run bugsnag-cli with upload unity-android --upload-api-root-url=http://localhost:9339 --api-key=1234567890ABCDEF1234567890ABCDEF --overwrite  --aab-path platforms-examples/Unity/UnityExample.aab platforms-examples/Unity/
    Then I wait to receive 5 sourcemaps
    Then the sourcemap is valid for the Android Build API
    Then the sourcemaps Content-Type header is valid multipart form-data
    And the sourcemap payload field "apiKey" equals "1234567890ABCDEF1234567890ABCDEF"
    And the sourcemap payload field "appId" equals "com.bugsnag.example.unity.android"
    And the sourcemap payload field "versionCode" equals "1"
    And the sourcemap payload field "versionName" equals "1.0"
    And the sourcemap payload field "overwrite" equals "true"
