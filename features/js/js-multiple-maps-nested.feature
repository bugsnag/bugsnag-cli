Feature: Js integration tests multiple nested source maps

  @BuildNestedJS
  Scenario: Searches in the dist folder automatically
    When I run bugsnag-cli upload "js" with the following arguments:
      | --upload-api-root-url          | http://localhost:$MAZE_RUNNER_PORT    |
      | --api-key                      | 1234567890ABCDEF1234567890ABCDEF      |
      | --base-url                     | example.com                           |
      | features/base-fixtures/js/out/ |                                       |
    And I wait to receive 4 sourcemaps
    Then the sourcemaps are valid for the API
    Then the sourcemaps Content-Type header is valid multipart form-data
    And the sourcemap payload field "apiKey" equals "1234567890ABCDEF1234567890ABCDEF"
    And the sourcemap payload field "appVersion" equals "3.0.0"
    And the sourcemap payload field "minifiedUrl" equals "example.com/dir1/file1.js"
    And the sourcemap payload field "sourceMap" is valid json
    And the sourcemap payload field "minifiedFile" is not empty
    And the sourcemap payload field "projectRoot" ends with "base-fixtures/js"

    And I discard the oldest sourcemap

    Then the sourcemaps are valid for the API
    Then the sourcemaps Content-Type header is valid multipart form-data
    And the sourcemap payload field "apiKey" equals "1234567890ABCDEF1234567890ABCDEF"
    And the sourcemap payload field "appVersion" equals "3.0.0"
    And the sourcemap payload field "minifiedUrl" equals "example.com/dir2/dir22/file3.js"
    And the sourcemap payload field "sourceMap" is valid json
    And the sourcemap payload field "minifiedFile" is not empty
    And the sourcemap payload field "projectRoot" ends with "base-fixtures/js"

    And I discard the oldest sourcemap

    Then the sourcemaps are valid for the API
    Then the sourcemaps Content-Type header is valid multipart form-data
    And the sourcemap payload field "apiKey" equals "1234567890ABCDEF1234567890ABCDEF"
    And the sourcemap payload field "appVersion" equals "3.0.0"
    And the sourcemap payload field "minifiedUrl" equals "example.com/dir2/file2.js"
    And the sourcemap payload field "sourceMap" is valid json
    And the sourcemap payload field "minifiedFile" is not empty
    And the sourcemap payload field "projectRoot" ends with "base-fixtures/js"

    And I discard the oldest sourcemap

    Then the sourcemaps are valid for the API
    Then the sourcemaps Content-Type header is valid multipart form-data
    And the sourcemap payload field "apiKey" equals "1234567890ABCDEF1234567890ABCDEF"
    And the sourcemap payload field "appVersion" equals "3.0.0"
    And the sourcemap payload field "minifiedUrl" equals "example.com/index.js"
    And the sourcemap payload field "sourceMap" is valid json
    And the sourcemap payload field "minifiedFile" is not empty
    And the sourcemap payload field "projectRoot" ends with "base-fixtures/js"
