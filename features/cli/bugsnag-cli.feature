Feature: Basic CLI behavior

  Scenario: Starting bugsnag-cli on mac without any flags
    When I run bugsnag-cli
    Then I should see the help banner

  Scenario: Starting bugsnag-cli upload all on mac without an API Key
    When I run bugsnag-cli with upload all
    Then I should see the missing path error

  Scenario: Starting bugsnag-cli upload all on mac without a path
    When I run bugsnag-cli with upload all features/dart/fixtures/app-debug-info/app.android-arm.symbols
    Then I should see the API Key error

  Scenario: Starting bugsnag-cli upload all with an invalid path
    When I run bugsnag-cli with upload all --api-key=1234567890ABCDEF1234567890ABCDEF /path/to/no/file
    Then I should see the no such file or directory error

  Scenario: Starting bugsnag-cli and checking the version
    When I run bugsnag-cli with --version
    Then the version number should match the version set in main.go
