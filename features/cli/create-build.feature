Feature: Bugsnag CLI create-build behavior
  Scenario: Starting bugsnag-cli create build on mac without an API Key
    When I run bugsnag-cli with create-build
    Then I should see the API Key error

  Scenario: Starting bugsnag-cli create-build on mac without any options
    When I run bugsnag-cli with create-build --api-key=1234567890ABCDEF1234567890ABCDEF
    Then I should see the missing app version error

  Scenario: Starting bugsnag-cli create-build on mac with app-version
    When I run bugsnag-cli with create-build --build-api-root-url=http://localhost:$MAZE_RUNNER_PORT/builds --api-key=1234567890ABCDEF1234567890ABCDEF --version-name=1.2.3
    And I wait to receive 1 builds
    Then the build is valid for the Builds API
    And the builds payload field "apiKey" equals "1234567890ABCDEF1234567890ABCDEF"
    And the builds payload field "appVersion" equals "1.2.3"
    And the builds payload field "sourceControl.repository" equals "https://github.com/bugsnag/bugsnag-cli.git"

  Scenario: Starting bugsnag-cli create-build and passing an AAB file
    When I run bugsnag-cli with create-build --build-api-root-url=http://localhost:$MAZE_RUNNER_PORT/builds --api-key=1234567890ABCDEF1234567890ABCDEF --android-aab=features/android/fixtures/aab/app-release.aab
    And I wait to receive 1 builds
    Then the build is valid for the Builds API
    And the builds payload field "apiKey" equals "1234567890ABCDEF1234567890ABCDEF"
    And the builds payload field "appVersion" equals "1.0"
    And the builds payload field "sourceControl.repository" equals "https://github.com/bugsnag/bugsnag-cli.git"

  Scenario: Starting bugsnag-cli create-build and passing an Android manifest file
    When I run bugsnag-cli with create-build --build-api-root-url=http://localhost:$MAZE_RUNNER_PORT/builds --app-manifest=features/android/fixtures/app/build/intermediates/merged_manifests/release/AndroidManifest.xml
    And I wait to receive 1 builds
    Then the build is valid for the Builds API
    And the builds payload field "apiKey" equals "1234567890ABCDEF1234567890ABCDEF"
    And the builds payload field "appVersion" equals "1.0"
    And the builds payload field "sourceControl.repository" equals "https://github.com/bugsnag/bugsnag-cli.git"

  Scenario: Starting bugsnag-cli create-build with invalid source control provider
    When I run bugsnag-cli with create-build --build-api-root-url=http://localhost:$MAZE_RUNNER_PORT/builds --api-key=1234567890ABCDEF1234567890ABCDEF --version-name=1.2.3 --provider=test
    Then I should see the not an accepted value for the source control provider error

  Scenario: Starting bugsnag-cli create-build with no source control provider
    When I run bugsnag-cli with create-build --build-api-root-url=http://localhost:$MAZE_RUNNER_PORT/builds --api-key=1234567890ABCDEF1234567890ABCDEF --version-name=1.2.3 --provider=
    Then I should see the missing source control provider error

  Scenario: Starting bugsnag-cli create-build and passing an Android manifest file with dry-run and verbose
    When I run bugsnag-cli with create-build --build-api-root-url=http://localhost:$MAZE_RUNNER_PORT/builds --app-manifest=features/android/fixtures/app/build/intermediates/merged_manifests/release/AndroidManifest.xml --dry-run --verbose
    Then I should see the build payload
    And I wait to receive 0 builds
