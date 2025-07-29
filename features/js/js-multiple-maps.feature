Feature: Js integration tests multiple source maps

  Scenario: Searches in the dist folder automatically
    When I run bugsnag-cli with upload js --upload-api-root-url=http://localhost:$MAZE_RUNNER_PORT --api-key=1234567890ABCDEF1234567890ABCDEF --overwrite --base-url=example.com --project-root=features/js/fixtures/js-multiple-maps features/js/fixtures/js-multiple-maps/dist/
    And I wait to receive 2 sourcemaps
    Then the sourcemaps are valid for the API
    Then the sourcemaps Content-Type header is valid multipart form-data
    And the sourcemap payload field "apiKey" equals "1234567890ABCDEF1234567890ABCDEF"
    And the sourcemap payload field "appVersion" equals "1.2.3"
    And the sourcemap payload field "minifiedUrl" equals "example.com/main.js"
    And the sourcemap payload field "sourceMap" is valid json
    And the sourcemap payload field "minifiedFile" is not empty
    And the sourcemap payload field "projectRoot" ends with "features/js/fixtures/js-multiple-maps"
    And the sourcemap payload field "overwrite" equals "true"
