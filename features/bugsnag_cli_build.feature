Feature: Bugsnag CLI create-build behavior
  Scenario: Starting bugsnag-cli create build on mac without an API Key
    When I run bugsnag-cli with create-build
    Then I should see the API Key error

  Scenario: Starting bugsnag-cli create-build on mac without any options
    When I run bugsnag-cli with create-build --api-key=1234567890ABCDEF1234567890ABCDEF
    Then I should see the missing app version error

  Scenario: Starting bugsnag-cli create-build on mac with app-version
    When I run bugsnag-cli with create-build --api-key=1234567890ABCDEF1234567890ABCDEF --app-version=1.2.3
    Then the payload should match local information