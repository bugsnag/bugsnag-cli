Feature: Node.js Integration Tests
  @CleanAndBuildNodeJs
  Scenario: Upload Node.js sourcemaps providing a directory
    When I run bugsnag-cli with upload js --upload-api-root-url=http://localhost:$MAZE_RUNNER_PORT --api-key=1234567890ABCDEF1234567890ABCDEF --base-url=example.com --version-name=2.3.4 features/node/fixtures/dist
    And I wait to receive 2 sourcemaps
    Then the sourcemaps are valid for the API
    Then the sourcemaps Content-Type header is valid multipart form-data
    And the sourcemap payload field "apiKey" equals "1234567890ABCDEF1234567890ABCDEF"
    And the sourcemap payload field "appVersion" equals "2.3.4"
    And the sourcemap payload field "minifiedUrl" equals "example.com/index.js"
    And the sourcemap payload field "sourceMap" is valid json
    And the sourcemap payload field "minifiedFile" is not empty
    And the sourcemap payload field "projectRoot" ends with "features/node/fixtures"

    And I discard the oldest sourcemaps

    And the sourcemap payload field "apiKey" equals "1234567890ABCDEF1234567890ABCDEF"
    And the sourcemap payload field "appVersion" equals "2.3.4"
    And the sourcemap payload field "minifiedUrl" equals "example.com/utils.js"
    And the sourcemap payload field "sourceMap" is valid json
    And the sourcemap payload field "minifiedFile" is not empty
    And the sourcemap payload field "projectRoot" ends with "features/node/fixtures"

  @CleanAndBuildNodeJs
  Scenario: Upload a single Node.js sourcemap using all CLI flags
    When I run bugsnag-cli with upload js --upload-api-root-url=http://localhost:$MAZE_RUNNER_PORT --api-key=1234567890ABCDEF1234567890ABCDEF --bundle-url=example.com --version-name=2.3.4 --source-map=features/node/fixtures/dist/index.js.map --bundle=features/node/fixtures/dist/index.js --project-root=features/node/fixtures features/node/fixtures/dist
    And I wait to receive 1 sourcemaps
    Then the sourcemaps are valid for the API
    Then the sourcemaps Content-Type header is valid multipart form-data
    And the sourcemap payload field "apiKey" equals "1234567890ABCDEF1234567890ABCDEF"
    And the sourcemap payload field "appVersion" equals "2.3.4"
    And the sourcemap payload field "minifiedUrl" equals "example.com"
    And the sourcemap payload field "sourceMap" is valid json
    And the sourcemap payload field "minifiedFile" is not empty
    And the sourcemap payload field "projectRoot" ends with "features/node/fixtures"
