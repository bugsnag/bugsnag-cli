Feature: Android NDK Integration Test

  Scenario: Upload a single Android NDK sourcemap using all CLI flags
    When I run bugsnag-cli with upload android-ndk --upload-api-root-url=http://localhost:9339 --api-key=1234567890ABCDEF1234567890ABCDEF --overwrite --application-id="2.0" --app-manifest=features/android/fixtures/app/build/intermediates/merged_manifests/release/AndroidManifest.xml --variant=release --version-code=2 --version-name=2.0 features/android/fixtures/app/build/intermediates/merged_native_libs/release/out/lib/arm64-v8a/libbugsnag-ndk.so
    And I wait to receive 1 sourcemaps
    Then the sourcemap is valid for the NDK Build API
    Then the sourcemaps Content-Type header is valid multipart form-data
    And the sourcemap payload field "apiKey" equals "1234567890ABCDEF1234567890ABCDEF"
    And the sourcemap payload field "appId" equals "2.0"
    And the sourcemap payload field "versionCode" equals "2"
    And the sourcemap payload field "versionName" equals "2.0"
    And the sourcemap payload field "overwrite" equals "true"

  Scenario: Upload a single Android NDK sourcemap providing the app-manifest CLI flag
    When I run bugsnag-cli with upload android-ndk --upload-api-root-url=http://localhost:9339 --api-key=1234567890ABCDEF1234567890ABCDEF --overwrite --app-manifest=features/android/fixtures/app/build/intermediates/merged_manifests/release/AndroidManifest.xml features/android/fixtures/app/build/intermediates/merged_native_libs/release/out/lib/arm64-v8a/libbugsnag-ndk.so
    And I wait to receive 1 sourcemaps
    Then the sourcemap is valid for the NDK Build API
    Then the sourcemaps Content-Type header is valid multipart form-data
    And the sourcemap payload field "apiKey" equals "1234567890ABCDEF1234567890ABCDEF"
    And the sourcemap payload field "appId" equals "com.example.bugsnag.android"
    And the sourcemap payload field "versionCode" equals "1"
    And the sourcemap payload field "versionName" equals "1.0"
    And the sourcemap payload field "overwrite" equals "true"

  Scenario: Upload a single Android NDK sourcemap providing the app-manifest CLI flag
    When I run bugsnag-cli with upload android-ndk --upload-api-root-url=http://localhost:9339 --api-key=1234567890ABCDEF1234567890ABCDEF --overwrite features/android/fixtures/app/build/intermediates/merged_native_libs/release/out/lib/arm64-v8a/libbugsnag-ndk.so
    And I wait to receive 1 sourcemaps
    Then the sourcemap is valid for the NDK Build API
    Then the sourcemaps Content-Type header is valid multipart form-data
    And the sourcemap payload field "apiKey" equals "1234567890ABCDEF1234567890ABCDEF"
    And the sourcemap payload field "appId" equals "com.example.bugsnag.android"
    And the sourcemap payload field "versionCode" equals "1"
    And the sourcemap payload field "versionName" equals "1.0"
    And the sourcemap payload field "overwrite" equals "true"

  Scenario: Upload multiple Android NDK sourcemaps providing no flags to the CLI
    When I run bugsnag-cli with upload android-ndk --upload-api-root-url=http://localhost:9339 --api-key=1234567890ABCDEF1234567890ABCDEF --overwrite features/android/fixtures/
    And I wait to receive 16 sourcemaps
    Then the sourcemap is valid for the NDK Build API
    Then the sourcemaps Content-Type header is valid multipart form-data
    And the sourcemap payload field "apiKey" equals "1234567890ABCDEF1234567890ABCDEF"
    And the sourcemap payload field "appId" equals "com.example.bugsnag.android"
    And the sourcemap payload field "versionCode" equals "1"
    And the sourcemap payload field "versionName" equals "1.0"
    And the sourcemap payload field "overwrite" equals "true"
