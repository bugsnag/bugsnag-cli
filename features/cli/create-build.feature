Feature: Bugsnag CLI create-build behavior
  Scenario: Starting bugsnag-cli create build on mac without an API Key
    When I run bugsnag-cli with create-build
    Then I should see the API Key error

  Scenario: Starting bugsnag-cli create-build on mac without any options
    When I run bugsnag-cli with create-build --api-key=1234567890ABCDEF1234567890ABCDEF
    Then I should see the missing app version error

  Scenario: Starting bugsnag-cli create-build on mac with app-version
    When I run bugsnag-cli with create-build --build-api-root-url=http://localhost:9339/builds --api-key=1234567890ABCDEF1234567890ABCDEF --version-name=1.2.3
    And I wait to receive 1 builds

    Then the build is valid for the Builds API

    And the builds payload field "apiKey" equals "1234567890ABCDEF1234567890ABCDEF"
    And the builds payload field "appVersion" equals "1.2.3"
    And the builds payload field "sourceControl.repository" equals "git@github.com:bugsnag/bugsnag-cli"
