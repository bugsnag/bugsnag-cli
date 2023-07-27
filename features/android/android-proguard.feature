Feature: Android Proguard Integration Test

  Scenario: Upload Android Proguard
    When I run bugsnag-cli with upload android-proguard --upload-api-root-url=http://localhost:9339 --api-key=1234567890ABCDEF1234567890ABCDEF --overwrite --app-manifest=features/fixtures/android/app/build/intermediates/merged_manifests/release/AndroidManifest.xml features/fixtures/android/app/build/outputs/mapping/release/mapping.txt
    And I wait to receive 1 sourcemaps

    Then the sourcemap is valid for the Proguard Build API

    Then the sourcemaps Content-Type header is valid multipart form-data
    And the sourcemap payload field "apiKey" equals "1234567890ABCDEF1234567890ABCDEF"
    And the sourcemap payload field "versionCode" equals "1"
    And the sourcemap payload field "versionName" equals "1.0"
    And the sourcemap payload field "buildUUID" equals "53e067c2-f338-455d-a4f1-51e2033e89ed"
    And the sourcemap payload field "overwrite" equals "true"

  Scenario: Locate and Upload Android Proguard
    When I run bugsnag-cli with upload android-proguard --upload-api-root-url=http://localhost:9339 --api-key=1234567890ABCDEF1234567890ABCDEF --overwrite features/fixtures/android/
    And I wait to receive 1 sourcemaps

    Then the sourcemap is valid for the Proguard Build API

    Then the sourcemaps Content-Type header is valid multipart form-data
    And the sourcemap payload field "apiKey" equals "1234567890ABCDEF1234567890ABCDEF"
    And the sourcemap payload field "versionCode" equals "1"
    And the sourcemap payload field "versionName" equals "1.0"
    And the sourcemap payload field "buildUUID" equals "53e067c2-f338-455d-a4f1-51e2033e89ed"
    And the sourcemap payload field "overwrite" equals "true"
