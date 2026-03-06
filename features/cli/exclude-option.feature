Feature: Exclude option tests

  Scenario: Exclude files with wildcard extension pattern
    When I run bugsnag-cli with upload js --upload-api-root-url=http://localhost:$MAZE_RUNNER_PORT --api-key=1234567890ABCDEF1234567890ABCDEF --overwrite --base-url=example.com --exclude=*.map --project-root=features/js/fixtures/js-multiple-maps features/js/fixtures/js-multiple-maps/dist/
    And I wait for 2 seconds
    Then I should receive no sourcemaps

  Scenario: Exclude files matching specific filename pattern
    When I run bugsnag-cli with upload js --upload-api-root-url=http://localhost:$MAZE_RUNNER_PORT --api-key=1234567890ABCDEF1234567890ABCDEF --overwrite --base-url=example.com --exclude=main.js.map --project-root=features/js/fixtures/js-multiple-maps features/js/fixtures/js-multiple-maps/dist/
    And I wait to receive 1 sourcemap
    Then the sourcemaps are valid for the API
    And the sourcemap payload field "minifiedUrl" equals "example.com/other.js"

  Scenario: Exclude files with multiple patterns
    When I run bugsnag-cli with upload js --upload-api-root-url=http://localhost:$MAZE_RUNNER_PORT --api-key=1234567890ABCDEF1234567890ABCDEF --overwrite --base-url=example.com --exclude=main.js.map --exclude=other.js.map --project-root=features/js/fixtures/js-multiple-maps features/js/fixtures/js-multiple-maps/dist/
    And I wait for 2 seconds
    Then I should receive no sourcemaps

  Scenario: Exclude with path pattern
    When I run bugsnag-cli with upload js --upload-api-root-url=http://localhost:$MAZE_RUNNER_PORT --api-key=1234567890ABCDEF1234567890ABCDEF --overwrite --base-url=example.com --exclude=/dist/ --project-root=features/js/fixtures/js-multiple-maps features/js/fixtures/js-multiple-maps/dist/
    And I wait for 2 seconds
    Then I should receive no sourcemaps

  Scenario: Upload succeeds when exclude pattern doesn't match
    When I run bugsnag-cli with upload js --upload-api-root-url=http://localhost:$MAZE_RUNNER_PORT --api-key=1234567890ABCDEF1234567890ABCDEF --overwrite --base-url=example.com --exclude=*.log --project-root=features/js/fixtures/js-multiple-maps features/js/fixtures/js-multiple-maps/dist/
    And I wait to receive 2 sourcemaps
    Then the sourcemaps are valid for the API
