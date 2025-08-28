Feature: Android NDK Integration Test

  Scenario: Upload a single Android NDK sourcemap using all CLI flags
    When I run bugsnag-cli upload "android-ndk" with the following arguments:
      | --upload-api-root-url                                                                                            | http://localhost:$MAZE_RUNNER_PORT                                                             |
      | --api-key                                                                                                        | 1234567890ABCDEF1234567890ABCDEF                                                               |
      | --application-id                                                                                                 | 2.0                                                                                            |
      | --app-manifest                                                                                                   | features/android/fixtures/app/build/intermediates/merged_manifests/release/AndroidManifest.xml |
      | --variant                                                                                                        | release                                                                                         |
      | --version-code                                                                                                   | 2                                                                                               |
      | --version-name                                                                                                   | 2.0                                                                                            |
      | features/android/fixtures/app/build/intermediates/merged_native_libs/release/out/lib/arm64-v8a/libbugsnag-ndk.so |                                                                                                |
    And I wait to receive 1 sourcemaps
    Then the sourcemaps are valid for the API
    Then the sourcemaps Content-Type header is valid multipart form-data
    Then the sourcemap payload fields should be:
      | apiKey       | 1234567890ABCDEF1234567890ABCDEF     |
      | appId        | 2.0                                  |
      | versionCode  | 2                                    |
      | versionName  | 2.0                                  |

  Scenario: Upload a single Android NDK sourcemap providing the app-manifest CLI flag
    When I run bugsnag-cli upload "android-ndk" with the following arguments:
      | --upload-api-root-url                                                                                            | http://localhost:$MAZE_RUNNER_PORT                                                             |
      | --api-key                                                                                                        | 1234567890ABCDEF1234567890ABCDEF                                                               |
      | --app-manifest                                                                                                   | features/android/fixtures/app/build/intermediates/merged_manifests/release/AndroidManifest.xml |
      | features/android/fixtures/app/build/intermediates/merged_native_libs/release/out/lib/arm64-v8a/libbugsnag-ndk.so |                                                                                                |
    And I wait to receive 1 sourcemaps
    Then the sourcemaps are valid for the API
    Then the sourcemaps Content-Type header is valid multipart form-data
    Then the sourcemap payload fields should be:
      | apiKey       | 1234567890ABCDEF1234567890ABCDEF     |
      | appId        | com.example.bugsnag.android          |
      | versionCode  | 1                                    |
      | versionName  | 1.0                                  |

  Scenario: Upload a single Android NDK sourcemap
    When I run bugsnag-cli upload "android-ndk" with the following arguments:
      | --upload-api-root-url                                                                                            | http://localhost:$MAZE_RUNNER_PORT                                                             |
      | --api-key                                                                                                        | 1234567890ABCDEF1234567890ABCDEF                                                               |
      | features/android/fixtures/app/build/intermediates/merged_native_libs/release/out/lib/arm64-v8a/libbugsnag-ndk.so |                                                                                                |
    And I wait to receive 1 sourcemaps
    Then the sourcemaps are valid for the API
    Then the sourcemaps Content-Type header is valid multipart form-data
    Then the sourcemap payload fields should be:
      | apiKey       | 1234567890ABCDEF1234567890ABCDEF     |
      | appId        | com.example.bugsnag.android          |
      | versionCode  | 1                                    |
      | versionName  | 1.0                                  |

  Scenario: Upload multiple Android NDK sourcemaps when command is run from within app directory
    When I run bugsnag-cli upload "android-ndk" with the following arguments:
      | --upload-api-root-url          | http://localhost:$MAZE_RUNNER_PORT |
      | --api-key                      | 1234567890ABCDEF1234567890ABCDEF   |
      | features/android/fixtures/app/ |                                    |
    And I wait to receive 16 sourcemaps
    Then the sourcemaps are valid for the API
    Then the sourcemaps Content-Type header is valid multipart form-data
    Then the sourcemap payload fields should be:
      | apiKey       | 1234567890ABCDEF1234567890ABCDEF     |
      | appId        | com.example.bugsnag.android          |
      | versionCode  | 1                                    |
      | versionName  | 1.0                                  |

  Scenario: Upload multiple Android NDK sourcemaps when command is run from within x86 directory
    When I run bugsnag-cli upload "android-ndk" with the following arguments:
      | --upload-api-root-url                                                                    | http://localhost:$MAZE_RUNNER_PORT |
      | --api-key                                                                                | 1234567890ABCDEF1234567890ABCDEF   |
      | features/android/fixtures/app/build/intermediates/merged_native_libs/release/out/lib/x86 |                                    |
    And I wait to receive 4 sourcemaps
    Then the sourcemaps are valid for the API
    Then the sourcemaps Content-Type header is valid multipart form-data
    Then the sourcemap payload fields should be:
      | apiKey       | 1234567890ABCDEF1234567890ABCDEF     |
      | appId        | com.example.bugsnag.android          |
      | versionCode  | 1                                    |
      | versionName  | 1.0                                  |

  Scenario: Upload multiple Android NDK sourcemaps providing no flags to the CLI
    When I run bugsnag-cli upload "android-ndk" with the following arguments:
      | --upload-api-root-url      | http://localhost:$MAZE_RUNNER_PORT |
      | --api-key                  | 1234567890ABCDEF1234567890ABCDEF   |
      | features/android/fixtures/ |                                    |
    And I wait to receive 16 sourcemaps
    Then the sourcemaps are valid for the API
    Then the sourcemaps Content-Type header is valid multipart form-data
    Then the sourcemap payload fields should be:
      | apiKey       | 1234567890ABCDEF1234567890ABCDEF     |
      | appId        | com.example.bugsnag.android          |
      | versionCode  | 1                                    |
      | versionName  | 1.0                                  |

  Scenario: Build and Upload Android NDK sourcemaps
    When I make the "features/base-fixtures/android"
    And I wait for the build to succeed

    When I run bugsnag-cli upload "android-ndk" with the following arguments:
      | --upload-api-root-url      | http://localhost:$MAZE_RUNNER_PORT |
      | --api-key                  | 1234567890ABCDEF1234567890ABCDEF   |
      | features/base-fixtures/android/ |                                    |
    And I wait to receive 4 sourcemaps
    Then the sourcemaps are valid for the API
    Then the sourcemaps Content-Type header is valid multipart form-data
    Then the sourcemap payload fields should be:
      | apiKey       | 1234567890ABCDEF1234567890ABCDEF     |
      | appId        | com.example.picoapp                  |
      | versionCode  | 1                                    |
      | versionName  | 1.0                                  |
