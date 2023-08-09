Feature: React Native 0.72 Android Integration Tests

  Scenario: Upload a single React Native 0.72 Android sourcemap using all CLI flags
    When I run bugsnag-cli with upload react-native-android --upload-api-root-url=http://localhost:9339 --api-key=1234567890ABCDEF1234567890ABCDEF --overwrite --app-manifest=features/react-native-android/fixtures/rn0_72/android/app/build/intermediates/merged_manifests/release/AndroidManifest.xml --bundle=features/react-native-android/fixtures/rn0_72/android/app/build/generated/sourcemaps/react/release/index.android.bundle.map --code-bundle-id=1.0-15 --dev --source-map=features/react-native-android/fixtures/rn0_72/android/app/build/generated/sourcemaps/react/release/index.android.bundle.map --variant=release --version-name=1.0 --version-code=1
    And I wait to receive 1 sourcemaps
    Then the sourcemap is valid for the React Native Build API
    Then the sourcemaps Content-Type header is valid multipart form-data
    And the sourcemap payload field "apiKey" equals "1234567890ABCDEF1234567890ABCDEF"
    And the sourcemap payload field "appVersion" equals "1.0"
    And the sourcemap payload field "appVersionCode" equals "1"
    And the sourcemap payload field "codeBundleId" equals "1.0-15"
    And the sourcemap payload field "dev" equals "true"
    And the sourcemap payload field "platform" equals "android"
    And the sourcemap payload field "overwrite" equals "true"

  Scenario: Upload a single React Native 0.72 Android sourcemap providing the app-manifest, bundle and source-map CLI flag
    When I run bugsnag-cli with upload react-native-android --upload-api-root-url=http://localhost:9339 --api-key=1234567890ABCDEF1234567890ABCDEF --overwrite --app-manifest=features/react-native-android/fixtures/rn0_72/android/app/build/intermediates/merged_manifests/release/AndroidManifest.xml --bundle=features/react-native-android/fixtures/rn0_72/android/app/build/generated/sourcemaps/react/release/index.android.bundle.map --source-map=features/react-native-android/fixtures/rn0_72/android/app/build/generated/sourcemaps/react/release/index.android.bundle.map
    And I wait to receive 1 sourcemaps
    Then the sourcemap is valid for the React Native Build API
    Then the sourcemaps Content-Type header is valid multipart form-data
    And the sourcemap payload field "apiKey" equals "1234567890ABCDEF1234567890ABCDEF"
    And the sourcemap payload field "appVersion" equals "1.0"
    And the sourcemap payload field "appVersionCode" equals "1"
    And the sourcemap payload field "platform" equals "android"
    And the sourcemap payload field "overwrite" equals "true"

  Scenario: Upload a single React Native 0.72 Android sourcemap providing the app-manifest CLI flag
    When I run bugsnag-cli with upload react-native-android --upload-api-root-url=http://localhost:9339 --api-key=1234567890ABCDEF1234567890ABCDEF --overwrite --app-manifest=features/react-native-android/fixtures/rn0_72/android/app/build/intermediates/merged_manifests/release/AndroidManifest.xml features/react-native-android/fixtures/rn0_72
    And I wait to receive 1 sourcemaps
    Then the sourcemap is valid for the React Native Build API
    Then the sourcemaps Content-Type header is valid multipart form-data
    And the sourcemap payload field "apiKey" equals "1234567890ABCDEF1234567890ABCDEF"
    And the sourcemap payload field "appVersion" equals "1.0"
    And the sourcemap payload field "appVersionCode" equals "1"
    And the sourcemap payload field "platform" equals "android"
    And the sourcemap payload field "overwrite" equals "true"

  Scenario: Upload a single React Native 0.72 Android sourcemap providing no CLI flag
    When I run bugsnag-cli with upload react-native-android --upload-api-root-url=http://localhost:9339 --api-key=1234567890ABCDEF1234567890ABCDEF --overwrite features/react-native-android/fixtures/rn0_72
    And I wait to receive 1 sourcemaps
    Then the sourcemap is valid for the React Native Build API
    Then the sourcemaps Content-Type header is valid multipart form-data
    And the sourcemap payload field "apiKey" equals "1234567890ABCDEF1234567890ABCDEF"
    And the sourcemap payload field "appVersion" equals "1.0"
    And the sourcemap payload field "appVersionCode" equals "1"
    And the sourcemap payload field "platform" equals "android"
    And the sourcemap payload field "overwrite" equals "true"

  Scenario: Build and Upload React Native 0.72 Android sourcemaps
    When I make the "rn0-72-test-fixture"
    And I wait for the build to succeed

    When I run bugsnag-cli with upload react-native-android --upload-api-root-url=http://localhost:9339 --api-key=1234567890ABCDEF1234567890ABCDEF --overwrite features/base-fixtures/rn0_72
    And I wait to receive 1 sourcemaps
    Then the sourcemap is valid for the React Native Build API
    Then the sourcemaps Content-Type header is valid multipart form-data
    And the sourcemap payload field "apiKey" equals "1234567890ABCDEF1234567890ABCDEF"
    And the sourcemap payload field "appVersion" equals "1.0"
    And the sourcemap payload field "appVersionCode" equals "1"
    And the sourcemap payload field "platform" equals "android"
    And the sourcemap payload field "overwrite" equals "true"
