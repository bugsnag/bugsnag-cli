Feature: Bugsnag CLI Android AAB behavior
  Scenario: Uploading AAB with no overrides
    When I run bugsnag-cli with upload android-aab features/fixtures/min-app-release.aab
    Then "decafbaddecafbaddecafbaddecafbad" should be used as "API key"
    And "f3112c3dbdd73ae5dee677e407af196f101e97f5" should be used as "build ID"
