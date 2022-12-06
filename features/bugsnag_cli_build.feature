Feature: Bugsnag CLI Create Build behavior
  Scenario: Starting bugsnag-cli create build on mac without an API Key
    When I run bugsnag-cli with create-build
    Then I should see the API Key error

  Scenario: Starting bugsnag-cli upload dart on mac without a path
    When I run bugsnag-cli with create-build --api-key=1234567890ABCDEF1234567890ABCDEF
    Then I should see the missing app version error
