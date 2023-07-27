Feature: Android AAB Integration Test

  Scenario: Upload Android AAB
    When I run bugsnag-cli with upload android-aab --upload-api-root-url=http://localhost:9339 --api-key=1234567890ABCDEF1234567890ABCDEF --overwrite features/fixtures/android/aab/app-release.aab
    And I wait to receive 5 sourcemaps


  Scenario: Upload Android Dexguard AAB
    When I run bugsnag-cli with upload android-aab --upload-api-root-url=http://localhost:9339 --api-key=1234567890ABCDEF1234567890ABCDEF --overwrite features/fixtures/android/aab/app-release-dexguard.aab
    And I wait to receive 5 sourcemaps
