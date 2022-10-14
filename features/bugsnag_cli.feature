Feature: Basic CLI behavior

  Scenario: Starting bugsnag-cli on mac without any flags
    When I run bugsnag-cli on mac
    Then I should see the help banner

  Scenario: Starting bugsnag-cli upload all on mac without an API Key
    When I run bugsnag-cli upload all on mac without an API key
    Then I should see the API Key error

  Scenario: Starting bugsnag-cli upload all on mac without a path
    When I run bugsnag-cli upload all on mac without a path
    Then I should see the missing path error

  Scenario: Starting bugsnag-cli upload all with an invalid path
    When I run bugsnag-cli upload all on mac with an invalid path
    Then I should see the no such file or directory error
    