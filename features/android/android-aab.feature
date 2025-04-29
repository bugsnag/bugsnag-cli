Feature: Android AAB Integration Test

  Scenario: Uploading Android AAB file
    When I run bugsnag-cli with upload android-aab --upload-api-root-url=http://localhost:$MAZE_RUNNER_PORT --api-key=1234567890ABCDEF1234567890ABCDEF --overwrite features/android/fixtures/aab/app-release.aab
    And I wait to receive 1 sourcemaps
    Then the sourcemap is valid for the Android Build API
    Then the sourcemaps Content-Type header is valid multipart form-data
    And the sourcemap payload field "apiKey" equals "1234567890ABCDEF1234567890ABCDEF"
    And the sourcemap payload field "versionCode" equals "1"
    And the sourcemap payload field "versionName" equals "1.0"
    And the sourcemap payload field "overwrite" equals "true"

  Scenario: Uploading Android AAB file with --verbose
    When I run bugsnag-cli with upload android-aab --upload-api-root-url=http://localhost:$MAZE_RUNNER_PORT --api-key=1234567890ABCDEF1234567890ABCDEF --overwrite features/android/fixtures/aab/app-release.aab --verbose
    And I wait to receive 1 sourcemaps
    Then the sourcemap is valid for the Android Build API
    And "f3112c3dbdd73ae5dee677e407af196f101e97f5" should be used as "build ID"
    Then the sourcemaps Content-Type header is valid multipart form-data
    And the sourcemap payload field "apiKey" equals "1234567890ABCDEF1234567890ABCDEF"
    And the sourcemap payload field "versionCode" equals "1"
    And the sourcemap payload field "versionName" equals "1.0"
    And the sourcemap payload field "overwrite" equals "true"

  Scenario: Uploading Android AAB file with Dexguard
    When I run bugsnag-cli with upload android-aab --upload-api-root-url=http://localhost:$MAZE_RUNNER_PORT --api-key=1234567890ABCDEF1234567890ABCDEF --overwrite features/android/fixtures/aab/app-release-dexguard.aab
    And I wait to receive 5 sourcemaps
    Then the sourcemap is valid for the Android Build API
    Then the sourcemaps Content-Type header is valid multipart form-data
    And the sourcemap payload field "apiKey" equals "1234567890ABCDEF1234567890ABCDEF"
    And the sourcemap payload field "versionCode" equals "1"
    And the sourcemap payload field "versionName" equals "3.0"
    And the sourcemap payload field "overwrite" equals "true"

  Scenario: Uploading Android AAB file with Dexguard with --verbose
    When I run bugsnag-cli with upload android-aab --upload-api-root-url=http://localhost:$MAZE_RUNNER_PORT --api-key=1234567890ABCDEF1234567890ABCDEF --overwrite features/android/fixtures/aab/app-release-dexguard.aab --verbose
    And I wait to receive 5 sourcemaps
    Then the sourcemap is valid for the Android Build API
    And "fb0d77a7-5df2-4f47-a823-b011f89a2b70" should be used as "build ID"
    Then the sourcemaps Content-Type header is valid multipart form-data
    And the sourcemap payload field "apiKey" equals "1234567890ABCDEF1234567890ABCDEF"
    And the sourcemap payload field "versionCode" equals "1"
    And the sourcemap payload field "versionName" equals "3.0"
    And the sourcemap payload field "overwrite" equals "true"

  Scenario: Build and Upload Android AAB file
    When I make the "features/base-fixtures/android"
    And I wait for the build to succeed

    When I run bugsnag-cli with upload android-aab --upload-api-root-url=http://localhost:$MAZE_RUNNER_PORT --api-key=1234567890ABCDEF1234567890ABCDEF --overwrite features/base-fixtures/android/app/build/outputs/bundle/release/app-release.aab
    And I wait to receive 5 sourcemaps
    Then the sourcemap is valid for the Android Build API
    Then the sourcemaps Content-Type header is valid multipart form-data
    And the sourcemap payload field "apiKey" equals "1234567890ABCDEF1234567890ABCDEF"
    And the sourcemap payload field "versionCode" equals "1"
    And the sourcemap payload field "versionName" equals "1.0"
    And the sourcemap payload field "overwrite" equals "true"

  Scenario: Build and Upload Android AAB file with --verbose
    When I make the "features/base-fixtures/android"
    And I wait for the build to succeed

    When I run bugsnag-cli with upload android-aab --upload-api-root-url=http://localhost:$MAZE_RUNNER_PORT --api-key=1234567890ABCDEF1234567890ABCDEF --overwrite features/base-fixtures/android/app/build/outputs/bundle/release/app-release.aab --verbose
    And I wait to receive 5 sourcemaps
    Then the sourcemap is valid for the Android Build API
    And "f88f420ede59cd6695cea71aa0c7345eccd594cb" should be used as "build ID"
    Then the sourcemaps Content-Type header is valid multipart form-data
    And the sourcemap payload field "apiKey" equals "1234567890ABCDEF1234567890ABCDEF"
    And the sourcemap payload field "versionCode" equals "1"
    And the sourcemap payload field "versionName" equals "1.0"
    And the sourcemap payload field "overwrite" equals "true"

  Scenario: Uploading Android AAB file when command is run from project root
    When I run bugsnag-cli with upload android-aab --upload-api-root-url=http://localhost:$MAZE_RUNNER_PORT --api-key=1234567890ABCDEF1234567890ABCDEF --overwrite features/android/fixtures/ --verbose
    And I wait to receive 1 sourcemaps
    Then the sourcemap is valid for the Android Build API
    And "f3112c3dbdd73ae5dee677e407af196f101e97f5" should be used as "build ID"
    Then the sourcemaps Content-Type header is valid multipart form-data
    And the sourcemap payload field "apiKey" equals "1234567890ABCDEF1234567890ABCDEF"
    And the sourcemap payload field "versionCode" equals "1"
    And the sourcemap payload field "versionName" equals "1.0"
    And the sourcemap payload field "overwrite" equals "true"

  Scenario: Uploading Android AAB file when command is run from within app directory
    When I run bugsnag-cli with upload android-aab --upload-api-root-url=http://localhost:$MAZE_RUNNER_PORT --api-key=1234567890ABCDEF1234567890ABCDEF --overwrite features/android/fixtures/app/ --verbose
    And I wait to receive 1 sourcemaps
    Then the sourcemap is valid for the Android Build API
    And "f3112c3dbdd73ae5dee677e407af196f101e97f5" should be used as "build ID"
    Then the sourcemaps Content-Type header is valid multipart form-data
    And the sourcemap payload field "apiKey" equals "1234567890ABCDEF1234567890ABCDEF"
    And the sourcemap payload field "versionCode" equals "1"
    And the sourcemap payload field "versionName" equals "1.0"
    And the sourcemap payload field "overwrite" equals "true"

  Scenario: Uploading Android AAB file when more than a single AAB is found
    When I run bugsnag-cli with upload android-aab --upload-api-root-url=http://localhost:$MAZE_RUNNER_PORT --api-key=1234567890ABCDEF1234567890ABCDEF --overwrite features/android/fixtures/aab/
    Then I should see the path ambiguous error
