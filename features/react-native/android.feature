@BuildRNAndroid
Feature: React Native Android Integration Tests
  Scenario: Upload a single React Native Android sourcemap using all escape hatches
    When I run bugsnag-cli with upload react-native-android --upload-api-root-url=http://localhost:9339 --bundle=features/react-native/fixtures/generated/old-arch/$RN_VERSION/android/app/build/generated/assets/createBundleReleaseJsAndAssets/index.android.bundle --dev --source-map=features/react-native/fixtures/generated/old-arch/$RN_VERSION/android/app/build/generated/sourcemaps/react/release/index.android.bundle.map --version-name=2.0 --app-manifest=features/react-native/fixtures/generated/old-arch/$RN_VERSION/android/app/build/intermediates/merged_manifests/release/processReleaseManifest/AndroidManifest.xml --variant=release --version-code=2 --overwrite
    And I wait to receive 1 sourcemaps
    Then the sourcemap is valid for the React Native Build API
    Then the sourcemaps Content-Type header is valid multipart form-data
    And the sourcemap payload field "apiKey" equals "1234567890ABCDEF1234567890ABCDEF"
    And the sourcemap payload field "appVersion" equals "2.0"
    And the sourcemap payload field "appVersionCode" equals "2"
    And the sourcemap payload field "platform" equals "android"
    And the sourcemap payload field "overwrite" equals "true"

  Scenario: Upload a single React Native Android sourcemap, only passing bundle, sourcemap and manifest options
    When I run bugsnag-cli with upload react-native-android --upload-api-root-url=http://localhost:9339 --bundle=features/react-native/fixtures/generated/old-arch/$RN_VERSION/android/app/build/generated/assets/createBundleReleaseJsAndAssets/index.android.bundle --source-map=features/react-native/fixtures/generated/old-arch/$RN_VERSION/android/app/build/generated/sourcemaps/react/release/index.android.bundle.map --app-manifest=features/react-native/fixtures/generated/old-arch/$RN_VERSION/android/app/build/intermediates/merged_manifests/release/processReleaseManifest/AndroidManifest.xml --overwrite features/react-native/fixtures/generated/old-arch/$RN_VERSION
    And I wait to receive 1 sourcemaps
    Then the sourcemap is valid for the React Native Build API
    Then the sourcemaps Content-Type header is valid multipart form-data
    And the sourcemap payload field "apiKey" equals "1234567890ABCDEF1234567890ABCDEF"
    And the sourcemap payload field "appVersion" equals "1.0"
    And the sourcemap payload field "appVersionCode" equals "1"
    And the sourcemap payload field "platform" equals "android"
    And the sourcemap payload field "overwrite" equals "true"

  Scenario: Upload a single React Native Android sourcemap
    When I run bugsnag-cli with upload react-native-android --upload-api-root-url=http://localhost:9339 --overwrite features/react-native/fixtures/generated/old-arch/$RN_VERSION
    And I wait to receive 1 sourcemaps
    Then the sourcemap is valid for the React Native Build API
    Then the sourcemaps Content-Type header is valid multipart form-data
    And the sourcemap payload field "apiKey" equals "1234567890ABCDEF1234567890ABCDEF"
    And the sourcemap payload field "appVersion" equals "1.0"
    And the sourcemap payload field "appVersionCode" equals "1"
    And the sourcemap payload field "platform" equals "android"
    And the sourcemap payload field "overwrite" equals "true"
