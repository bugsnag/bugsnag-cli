Feature: Android Proguard Integration Test

  Scenario: Upload an Android Proguard mapping file using all CLI flags
    When I run bugsnag-cli with upload android-proguard --upload-api-root-url=http://localhost:9339 --api-key=1234567890ABCDEF1234567890ABCDEF --overwrite --application-id="com.exampleApp.android" --app-manifest=features/android/fixtures/app/build/intermediates/merged_manifests/release/AndroidManifest.xml --build-uuid=1234567890abcdefghijklmnopqrstuvwxyz --variant=release --version-code=2 --version-name=2.0 features/android/fixtures/app/build/outputs/mapping/release/mapping.txt
    And I wait to receive 1 sourcemaps
    Then the sourcemap is valid for the Proguard Build API
    Then the sourcemaps Content-Type header is valid multipart form-data
    And the sourcemap payload field "apiKey" equals "1234567890ABCDEF1234567890ABCDEF"
    And the sourcemap payload field "appId" equals "com.exampleApp.android"
    And the sourcemap payload field "buildUUID" equals "1234567890abcdefghijklmnopqrstuvwxyz"
    And the sourcemap payload field "versionCode" equals "2"
    And the sourcemap payload field "versionName" equals "2.0"
    And the sourcemap payload field "overwrite" equals "true"

  Scenario: Upload an Android Proguard mapping file providing the app-manifest CLI flag
    When I run bugsnag-cli with upload android-proguard --upload-api-root-url=http://localhost:9339 --api-key=1234567890ABCDEF1234567890ABCDEF --overwrite --app-manifest=features/android/fixtures/app/build/intermediates/merged_manifests/release/AndroidManifest.xml features/android/fixtures/app/build/outputs/mapping/release/mapping.txt
    And I wait to receive 1 sourcemaps
    Then the sourcemap is valid for the Proguard Build API
    Then the sourcemaps Content-Type header is valid multipart form-data
    And the sourcemap payload field "apiKey" equals "1234567890ABCDEF1234567890ABCDEF"
    And the sourcemap payload field "appId" equals "com.example.bugsnag.android"
    And the sourcemap payload field "buildUUID" equals "53e067c2-f338-455d-a4f1-51e2033e89ed"
    And the sourcemap payload field "versionCode" equals "1"
    And the sourcemap payload field "versionName" equals "1.0"
    And the sourcemap payload field "overwrite" equals "true"

  Scenario: Upload an Android Proguard mapping file providing no flags to the CLI
    When I run bugsnag-cli with upload android-proguard --upload-api-root-url=http://localhost:9339 --api-key=1234567890ABCDEF1234567890ABCDEF --overwrite features/android/fixtures/
    And I wait to receive 1 sourcemaps
    Then the sourcemap is valid for the Proguard Build API
    Then the sourcemaps Content-Type header is valid multipart form-data
    And the sourcemap payload field "apiKey" equals "1234567890ABCDEF1234567890ABCDEF"
    And the sourcemap payload field "appId" equals "com.example.bugsnag.android"
    And the sourcemap payload field "buildUUID" equals "53e067c2-f338-455d-a4f1-51e2033e89ed"
    And the sourcemap payload field "versionCode" equals "1"
    And the sourcemap payload field "versionName" equals "1.0"
    And the sourcemap payload field "overwrite" equals "true"

  Scenario: Build and Upload Android Proguard sourcemaps
    When I make the "android-test-fixture"
    And I wait for the build to succeed

    When I run bugsnag-cli with upload android-proguard --upload-api-root-url=http://localhost:9339 --api-key=1234567890ABCDEF1234567890ABCDEF --overwrite features/base-fixtures/android
    And I wait to receive 1 sourcemaps
    Then the sourcemap is valid for the Proguard Build API
    Then the sourcemaps Content-Type header is valid multipart form-data
    And the sourcemap payload field "apiKey" equals "1234567890ABCDEF1234567890ABCDEF"
    And the sourcemap payload field "appId" equals "com.example.picoapp"
    And the sourcemap payload field "versionCode" equals "1"
    And the sourcemap payload field "versionName" equals "1.0"
    And the sourcemap payload field "overwrite" equals "true"
