Feature: Android Proguard Integration Test

  Scenario: Upload an Android Proguard mapping file using all CLI flags
    When I run bugsnag-cli upload "android-proguard" with the following arguments:
      | --upload-api-root-url                                                   | http://localhost:$MAZE_RUNNER_PORT                                                             |
      | --api-key                                                               | 1234567890ABCDEF1234567890ABCDEF                                                               |
      | --application-id                                                        | com.exampleApp.android                                                                         |
      | --app-manifest                                                          | features/android/fixtures/app/build/intermediates/merged_manifests/release/AndroidManifest.xml |
      | --build-uuid                                                            | 1234567890abcdefghijklmnopqrstuvwxyz                                                           |
      | --variant                                                               | release                                                                                        |
      | --version-code                                                          | 2                                                                                              |
      | --version-name                                                          | 2.0                                                                                            |
      | features/android/fixtures/app/build/outputs/mapping/release/mapping.txt |                                                                                                |
    And I wait to receive 1 sourcemaps
    Then the sourcemaps are valid for the API
    Then the sourcemaps Content-Type header is valid multipart form-data
    Then the sourcemap payload fields should be:
      | apiKey       | 1234567890ABCDEF1234567890ABCDEF      |
      | appId        | com.exampleApp.android                |
      | buildUUID    | 1234567890abcdefghijklmnopqrstuvwxyz  |
      | versionCode  | 2                                     |
      | versionName  | 2.0                                   |

  Scenario: Upload an Android Proguard mapping file providing the app-manifest CLI flag
    When I run bugsnag-cli upload "android-proguard" with the following arguments:
      | --upload-api-root-url                                                   | http://localhost:$MAZE_RUNNER_PORT                                                             |
      | --api-key                                                               | 1234567890ABCDEF1234567890ABCDEF                                                               |
      | --app-manifest                                                          | features/android/fixtures/app/build/intermediates/merged_manifests/release/AndroidManifest.xml |
      | features/android/fixtures/app/build/outputs/mapping/release/mapping.txt |                                                                                                |
    And I wait to receive 1 sourcemaps
    Then the sourcemaps are valid for the API
    Then the sourcemaps Content-Type header is valid multipart form-data
    Then the sourcemap payload fields should be:
      | apiKey       | 1234567890ABCDEF1234567890ABCDEF      |
      | appId        | com.example.bugsnag.android           |
      | buildUUID    | 53e067c2-f338-455d-a4f1-51e2033e89ed  |
      | versionCode  | 1                                     |
      | versionName  | 1.0                                   |

  Scenario: Upload an Android Proguard mapping file providing no flags to the CLI
    When I run bugsnag-cli upload "android-proguard" with the following arguments:
      | --upload-api-root-url      | http://localhost:$MAZE_RUNNER_PORT|
      | --api-key                  | 1234567890ABCDEF1234567890ABCDEF  |
      | features/android/fixtures/ |                                   |
    And I wait to receive 1 sourcemaps
    Then the sourcemaps are valid for the API
    Then the sourcemaps Content-Type header is valid multipart form-data
    Then the sourcemap payload fields should be:
      | apiKey       | 1234567890ABCDEF1234567890ABCDEF      |
      | appId        | com.example.bugsnag.android           |
      | buildUUID    | 53e067c2-f338-455d-a4f1-51e2033e89ed  |
      | versionCode  | 1                                     |
      | versionName  | 1.0                                   |


  Scenario: Build and Upload Android Proguard sourcemaps
    When I make the "features/base-fixtures/android"
    And I wait for the build to succeed

    When I run bugsnag-cli upload "android-proguard" with the following arguments:
      | --upload-api-root-url           | http://localhost:$MAZE_RUNNER_PORT|
      | --api-key                       | 1234567890ABCDEF1234567890ABCDEF  |
      | features/base-fixtures/android/ |                                   |
    And I wait to receive 1 sourcemaps
    Then the sourcemaps are valid for the API
    Then the sourcemaps Content-Type header is valid multipart form-data
    Then the sourcemap payload fields should be:
      | apiKey       | 1234567890ABCDEF1234567890ABCDEF      |
      | appId        | com.example.picoapp                   |
      | versionCode  | 1                                     |
      | versionName  | 1.0                                   |
