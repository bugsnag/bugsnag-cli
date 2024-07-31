@installation
Feature: CLI Installation

  Scenario: Install the bugsnag-cli via NPM
    When I install the bugsnag-cli via 'npm' in a new directory
    Then the 'node_modules/.bin' directory should contain "bugsnag-cli"

  Scenario: Install the bugsnag-cli via YARN
    When I install the bugsnag-cli via 'yarn' in a new directory
    Then the 'node_modules/.bin' directory should contain "bugsnag-cli"
