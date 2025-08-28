Feature: React Native iOS Integration Tests
  @BuildRNiOS
  Scenario: Upload a single React Native iOS sourcemap
    When I run bugsnag-cli upload "react-native-ios" with the following arguments:
      | --upload-api-root-url                                         | http://localhost:$MAZE_RUNNER_PORT |
      | --api-key                                                     | 1234567890ABCDEF1234567890ABCDEF   |
      | features/react-native/fixtures/generated/old-arch/$RN_VERSION |                                    |
    And I wait to receive 1 sourcemaps
    Then the sourcemaps are valid for the API
    Then the sourcemaps Content-Type header is valid multipart form-data
    Then the sourcemap payload fields should be:
      | apiKey           | 1234567890ABCDEF1234567890ABCDEF |
      | appVersion       | 1.0                              |
      | appBundleVersion | 1                                |
      | platform         | ios                              |

  @BuildExportRNiOS
  Scenario: Upload a single React Native iOS sourcemap using escape hatches
    When I run bugsnag-cli upload "react-native-ios" with the following arguments:
      | --upload-api-root-url                                         | http://localhost:$MAZE_RUNNER_PORT                                                                                                      |
      | --api-key                                                     | 1234567890ABCDEF1234567890ABCDEF                                                                                                        |
      | --bundle                                                      | features/react-native/fixtures/generated/old-arch/$RN_VERSION/reactnative.xcarchive/Products/Applications/reactnative.app/main.jsbundle |
      | --source-map                                                  | features/react-native/fixtures/generated/old-arch/$RN_VERSION/ios/build/sourcemaps/main.jsbundle.map                                    |
      | --plist                                                       | features/react-native/fixtures/generated/old-arch/$RN_VERSION/reactnative.xcarchive/Products/Applications/reactnative.app/Info.plist    |
      | --xcode-project                                               | features/react-native/fixtures/generated/old-arch/$RN_VERSION/ios/reactnative.xcodeproj                                                 |
      | --scheme                                                      | reactnative                                                                                                                             |
      | --bundle-version                                              | 2                                                                                                                                       |
      | --version-name                                                | 2.0                                                                                                                                     |
      | features/react-native/fixtures/generated/old-arch/$RN_VERSION |                                                                                                                                         |
    And I wait to receive 1 sourcemaps
    Then the sourcemaps are valid for the API
    Then the sourcemaps Content-Type header is valid multipart form-data
    Then the sourcemap payload fields should be:
      | apiKey           | 1234567890ABCDEF1234567890ABCDEF |
      | appVersion       | 2.0                              |
      | appBundleVersion | 2                                |
      | platform         | ios                              |

  Scenario: Upload a single React Native iOS sourcemap using escape hatches
    When I run bugsnag-cli upload "react-native-ios" with the following arguments:
      | --upload-api-root-url                                         | http://localhost:$MAZE_RUNNER_PORT                                                                                                      |
      | --api-key                                                     | 1234567890ABCDEF1234567890ABCDEF                                                                                                        |
      | --bundle                                                      | features/react-native/fixtures/generated/old-arch/$RN_VERSION/reactnative.xcarchive/Products/Applications/reactnative.app/main.jsbundle |
      | --source-map                                                  | features/react-native/fixtures/generated/old-arch/$RN_VERSION/ios/build/sourcemaps/main.jsbundle.map                                    |
      | --plist                                                       | features/react-native/fixtures/generated/old-arch/$RN_VERSION/reactnative.xcarchive/Products/Applications/reactnative.app/Info.plist    |
      | --xcode-project                                               | features/react-native/fixtures/generated/old-arch/$RN_VERSION/ios/reactnative.xcodeproj                                                 |
      | --scheme                                                      | reactnative                                                                                                                             |
      | --bundle-version                                              | 2                                                                                                                                       |
      | --version-name                                                | 2.0                                                                                                                                     |
      | features/react-native/fixtures/generated/old-arch/$RN_VERSION |                                                                                                                                         |
    And I wait to receive 1 sourcemaps
    Then the sourcemaps are valid for the API
    Then the sourcemaps Content-Type header is valid multipart form-data
    Then the sourcemap payload fields should be:
      | apiKey           | 1234567890ABCDEF1234567890ABCDEF |
      | appVersion       | 1.0                              |
      | appBundleVersion | 1                                |
      | platform         | ios                              |
