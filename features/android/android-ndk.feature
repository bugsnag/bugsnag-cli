Feature: Android NDK Integration Test

  Scenario: Upload Single Android NDK Sourcemap
    When I run bugsnag-cli with upload android-ndk --upload-api-root-url=http://localhost:9339 --api-key=1234567890ABCDEF1234567890ABCDEF --overwrite --app-manifest=features/fixtures/android/app/build/intermediates/merged_manifests/release/AndroidManifest.xml --android-ndk-root="/Users/administrator/Library/Android/sdk/ndk/25.1.8937393" features/fixtures/android/app/build/intermediates/merged_native_libs/release/out/lib/arm64-v8a/libbugsnag-ndk.so
    And I wait to receive 1 sourcemaps

    Then the sourcemap is valid for the NDK Build API

    Then the sourcemaps Content-Type header is valid multipart form-data
    And the sourcemap payload field "apiKey" equals "1234567890ABCDEF1234567890ABCDEF"
    And the sourcemap payload field "versionCode" equals "1"
    And the sourcemap payload field "versionName" equals "1.0"
    And the sourcemap payload field "overwrite" equals "true"

  Scenario: Locate and Upload Android NDK
    When I run bugsnag-cli with upload android-ndk --upload-api-root-url=http://localhost:9339 --api-key=1234567890ABCDEF1234567890ABCDEF --overwrite features/fixtures/android/
    And I wait to receive 16 sourcemaps

    Then the sourcemap is valid for the NDK Build API

    Then the sourcemaps Content-Type header is valid multipart form-data
    And the sourcemap payload field "apiKey" equals "1234567890ABCDEF1234567890ABCDEF"
    And the sourcemap payload field "versionCode" equals "1"
    And the sourcemap payload field "versionName" equals "1.0"
    And the sourcemap payload field "overwrite" equals "true"
