Feature: Android AAB Integration Test

  Scenario: Upload Android AAB
    When I run bugsnag-cli with upload android-aab --upload-api-root-url=http://localhost:9339 --api-key=1234567890ABCDEF1234567890ABCDEF --overwrite features/fixtures/android/aab/app-release.aab
    And I wait to receive 5 sourcemaps

    Then the sourcemap is valid for the Android Build API

    Then the sourcemap is valid for the Android Build API
    Then the sourcemaps Content-Type header is valid multipart form-data
    And the sourcemap payload field "apiKey" equals "1234567890ABCDEF1234567890ABCDEF"
    And the sourcemap payload field "versionCode" equals "1"
    And the sourcemap payload field "versionName" equals "1.0"
    And the sourcemap payload field "overwrite" equals "true"

  Scenario: Upload Android Dexguard AAB
    When I run bugsnag-cli with upload android-aab --upload-api-root-url=http://localhost:9339 --api-key=1234567890ABCDEF1234567890ABCDEF --overwrite features/fixtures/android/aab/app-release-dexguard.aab
    And I wait to receive 5 sourcemaps

    Then the sourcemap is valid for the Android Build API
    Then the sourcemaps Content-Type header is valid multipart form-data
    And the sourcemap payload field "apiKey" equals "1234567890ABCDEF1234567890ABCDEF"
    And the sourcemap payload field "versionCode" equals "1"
    And the sourcemap payload field "versionName" equals "3.0"
    And the sourcemap payload field "overwrite" equals "true"

  Scenario: Uploading AAB with no overrides
    When I run bugsnag-cli with upload android-aab features/fixtures/min-app-release.aab
    Then "decafbaddecafbaddecafbaddecafbad" should be used as "API key"
    And "f3112c3dbdd73ae5dee677e407af196f101e97f5" should be used as "build ID"
