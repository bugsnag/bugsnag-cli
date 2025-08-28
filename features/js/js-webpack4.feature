Feature: Webpack 4 js Integration Tests

  Scenario: Upload a single js sourcemap using all CLI flags
    When I run bugsnag-cli upload "js" with the following arguments:
      | --upload-api-root-url                 | http://localhost:$MAZE_RUNNER_PORT                |
      | --api-key                             | 1234567890ABCDEF1234567890ABCDEF                  |
      | --bundle-url                          | example.com                                       |
      | --version-name                        | 2.3.4                                             |
      | --source-map                          | features/js/fixtures/js-webpack4/dist/main.js.map |
      | --bundle                              | features/js/fixtures/js-webpack4/dist/main.js     |
      | --project-root                        | features/js/fixtures/js-webpack4                  |
      | features/js/fixtures/js-webpack4/dist |                                                   |
    And I wait to receive 1 sourcemaps
    Then the sourcemaps are valid for the API
    Then the sourcemaps Content-Type header is valid multipart form-data
    And the sourcemap payload field "apiKey" equals "1234567890ABCDEF1234567890ABCDEF"
    And the sourcemap payload field "appVersion" equals "2.3.4"
    And the sourcemap payload field "minifiedUrl" equals "example.com"
    And the sourcemap payload field "sourceMap" is valid json
    And the sourcemap payload field "minifiedFile" is not empty
    And the sourcemap payload field "projectRoot" ends with "features/js/fixtures/js-webpack4"

  Scenario: Automatically resolves the version number based on the package.json
    When I run bugsnag-cli upload "js" with the following arguments:
      | --upload-api-root-url                 | http://localhost:$MAZE_RUNNER_PORT |
      | --api-key                             | 1234567890ABCDEF1234567890ABCDEF   |
      | --base-url                            | example.com                        |
      | --project-root                        | features/js/fixtures/js-webpack4   |
      | features/js/fixtures/js-webpack4/dist |                                    |
    And I wait to receive 1 sourcemaps
    Then the sourcemaps are valid for the API
    Then the sourcemaps Content-Type header is valid multipart form-data
    And the sourcemap payload field "apiKey" equals "1234567890ABCDEF1234567890ABCDEF"
    And the sourcemap payload field "appVersion" equals "1.2.3"
    And the sourcemap payload field "minifiedUrl" equals "example.com/main.js"
    And the sourcemap payload field "sourceMap" is valid json
    And the sourcemap payload field "minifiedFile" is not empty
    And the sourcemap payload field "projectRoot" ends with "features/js/fixtures/js-webpack4"

  Scenario: Resolves the path specified as the map
    When I run bugsnag-cli upload "js" with the following arguments:
      | --upload-api-root-url                             | http://localhost:$MAZE_RUNNER_PORT |
      | --api-key                                         | 1234567890ABCDEF1234567890ABCDEF   |
      | --bundle-url                                      | example.com                        |
      | --project-root                                    | features/js/fixtures/js-webpack4   |
      | features/js/fixtures/js-webpack4/dist/main.js.map |                                    |
    And I wait to receive 1 sourcemaps
    Then the sourcemaps are valid for the API
    Then the sourcemaps Content-Type header is valid multipart form-data
    And the sourcemap payload field "apiKey" equals "1234567890ABCDEF1234567890ABCDEF"
    And the sourcemap payload field "appVersion" equals "1.2.3"
    And the sourcemap payload field "minifiedUrl" equals "example.com"
    And the sourcemap payload field "sourceMap" is valid json
    And the sourcemap payload field "minifiedFile" is not empty
    And the sourcemap payload field "projectRoot" ends with "features/js/fixtures/js-webpack4"

  Scenario: Searches in the dist folder automatically
    When I run bugsnag-cli upload "js" with the following arguments:
      | --upload-api-root-url             | http://localhost:$MAZE_RUNNER_PORT |
      | --api-key                         | 1234567890ABCDEF1234567890ABCDEF   |
      | --base-url                        | example.com                        |
      | --project-root                    | features/js/fixtures/js-webpack4   |
      | features/js/fixtures/js-webpack4/ |                                    |
    And I wait to receive 1 sourcemaps
    Then the sourcemaps are valid for the API
    Then the sourcemaps Content-Type header is valid multipart form-data
    And the sourcemap payload field "apiKey" equals "1234567890ABCDEF1234567890ABCDEF"
    And the sourcemap payload field "appVersion" equals "1.2.3"
    And the sourcemap payload field "minifiedUrl" equals "example.com/dist/main.js"
    And the sourcemap payload field "sourceMap" is valid json
    And the sourcemap payload field "minifiedFile" is not empty
    And the sourcemap payload field "projectRoot" ends with "features/js/fixtures/js-webpack4"

  Scenario: Uses the working directory as project root
    When I run bugsnag-cli upload "js" with the following arguments:
      | --upload-api-root-url                   | http://localhost:$MAZE_RUNNER_PORT |
      | --api-key                               | 1234567890ABCDEF1234567890ABCDEF   |
      | --base-url                              | example.com                        |
      | features/js/fixtures/js-webpack4/dist/  |                                    |
    And I wait to receive 1 sourcemaps
    Then the sourcemaps are valid for the API
    Then the sourcemaps Content-Type header is valid multipart form-data
    And the sourcemap payload field "apiKey" equals "1234567890ABCDEF1234567890ABCDEF"
    And the sourcemap payload field "appVersion" equals "1.2.3"
    And the sourcemap payload field "minifiedUrl" equals "example.com/main.js"
    And the sourcemap payload field "sourceMap" is valid json
    And the sourcemap payload field "minifiedFile" is not empty
    And the sourcemap payload field "projectRoot" ends with "js-webpack4"

  Scenario: Base URL correctly appends the path
    When I run bugsnag-cli upload "js" with the following arguments:
      | --upload-api-root-url                   | http://localhost:$MAZE_RUNNER_PORT |
      | --api-key                               | 1234567890ABCDEF1234567890ABCDEF   |
      | --base-url                              | example.com                        |
      | features/js/fixtures/js-webpack4/dist/  |                                    |
    And I wait to receive 1 sourcemaps
    Then the sourcemaps are valid for the API
    Then the sourcemaps Content-Type header is valid multipart form-data
    And the sourcemap payload field "apiKey" equals "1234567890ABCDEF1234567890ABCDEF"
    And the sourcemap payload field "appVersion" equals "1.2.3"
    And the sourcemap payload field "minifiedUrl" equals "example.com/main.js"
    And the sourcemap payload field "sourceMap" is valid json
    And the sourcemap payload field "minifiedFile" is not empty
    And the sourcemap payload field "projectRoot" ends with "features/js/fixtures/js-webpack4"

  Scenario: Build and Upload js webpack4 sourcemaps
    When I make the "features/base-fixtures/js-webpack4"
    And I wait for the build to succeed

    When I run bugsnag-cli upload "js" with the following arguments:
      | --upload-api-root-url                     | http://localhost:$MAZE_RUNNER_PORT |
      | --api-key                                 | 1234567890ABCDEF1234567890ABCDEF   |
      | --base-url                                | example.com                        |
      | features/base-fixtures/js-webpack4/dist/  |                                    |
    And I wait to receive 1 sourcemaps
    Then the sourcemaps are valid for the API
    Then the sourcemaps Content-Type header is valid multipart form-data
    And the sourcemap payload field "apiKey" equals "1234567890ABCDEF1234567890ABCDEF"
    And the sourcemap payload field "appVersion" equals "1.2.3"
    And the sourcemap payload field "minifiedUrl" equals "example.com/main.js"
    And the sourcemap payload field "minifiedFile" is not empty
    And the sourcemap payload field "projectRoot" ends with "features/base-fixtures/js-webpack4"
