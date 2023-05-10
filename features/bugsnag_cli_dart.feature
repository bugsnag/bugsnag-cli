Feature: Bugsnag CLI Dart behavior
  Scenario: Starting bugsnag-cli upload dart on mac without an API Key
    When I run bugsnag-cli with upload dart
    Then I should see the API Key error

  Scenario: Starting bugsnag-cli upload dart on mac without a path
    When I run bugsnag-cli with upload dart --api-key=1234567890ABCDEF1234567890ABCDEF
    Then I should see the missing path error

  Scenario: Starting bugsnag-cli upload dart with an invalid path
    When I run bugsnag-cli with upload dart --api-key=1234567890ABCDEF1234567890ABCDEF /path/to/no/file
    Then I should see the no such file or directory error
    