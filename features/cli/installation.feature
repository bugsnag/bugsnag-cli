Feature: CLI Installation
  Scenario: Install the bugsnag-cli via NPM
    Given I package the bugsnag-cli
    And I init a new nodejs project
    Then I install the bugsnag-cli via npm in the new project
    Then the node_modules bin directory should contain "bugsnag-cli"
    Then I clean up the project

  Scenario: Install the bugsnag-cli via YARN
    Given I package the bugsnag-cli
    And I init a new nodejs project
    Then I install the bugsnag-cli via yarn in the new project
    Then the node_modules bin directory should contain "bugsnag-cli"
    Then I clean up the project
