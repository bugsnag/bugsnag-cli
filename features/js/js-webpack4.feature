Feature: Webpack 4 js Integration Tests

  Scenario: Upload a single js sourcemap using all CLI flags
    When I run bugsnag-cli with upload js --upload-api-root-url=http://localhost:9339 --api-key=1234567890ABCDEF1234567890ABCDEF --overwrite --bundle-url=example.com --version-name=2.3.4 --source-map=features/js/fixtures/js-webpack4/dist/main.js.map --bundle=features/js/fixtures/js-webpack4/dist/main.js --project-root=features/js/fixtures/js-webpack4 features/js/fixtures/js-webpack4/dist
    And I wait to receive 1 sourcemaps
    Then the sourcemap is valid for the JS Build API
    Then the sourcemaps Content-Type header is valid multipart form-data
    And the sourcemap payload field "apiKey" equals "1234567890ABCDEF1234567890ABCDEF"
    And the sourcemap payload field "appVersion" equals "2.3.4"
    And the sourcemap payload field "minifiedUrl" equals "example.com"
    And the sourcemap payload field "sourceMap" is valid json
    And the sourcemap payload field "minifiedFile" is not empty
    And the sourcemap payload field "projectRoot" ends with "features/js/fixtures/js-webpack4"
    And the sourcemap payload field "overwrite" equals "true"

  Scenario: Automatically resolves the version number based on the package.json
    When I run bugsnag-cli with upload js --upload-api-root-url=http://localhost:9339 --api-key=1234567890ABCDEF1234567890ABCDEF --overwrite --base-url=example.com --project-root=features/js/fixtures/js-webpack4 features/js/fixtures/js-webpack4/dist
    And I wait to receive 1 sourcemaps
    Then the sourcemap is valid for the JS Build API
    Then the sourcemaps Content-Type header is valid multipart form-data
    And the sourcemap payload field "apiKey" equals "1234567890ABCDEF1234567890ABCDEF"
    And the sourcemap payload field "appVersion" equals "1.2.3"
    And the sourcemap payload field "minifiedUrl" equals "example.com/main.js"
    And the sourcemap payload field "sourceMap" is valid json
    And the sourcemap payload field "minifiedFile" is not empty
    And the sourcemap payload field "projectRoot" ends with "features/js/fixtures/js-webpack4"
    And the sourcemap payload field "overwrite" equals "true"

  Scenario: Resolves the path specified as the map
    When I run bugsnag-cli with upload js --upload-api-root-url=http://localhost:9339 --api-key=1234567890ABCDEF1234567890ABCDEF --overwrite --bundle-url=example.com --project-root=features/js/fixtures/js-webpack4 features/js/fixtures/js-webpack4/dist/main.js.map
    And I wait to receive 1 sourcemaps
    Then the sourcemap is valid for the JS Build API
    Then the sourcemaps Content-Type header is valid multipart form-data
    And the sourcemap payload field "apiKey" equals "1234567890ABCDEF1234567890ABCDEF"
    And the sourcemap payload field "appVersion" equals "1.2.3"
    And the sourcemap payload field "minifiedUrl" equals "example.com"
    And the sourcemap payload field "sourceMap" is valid json
    And the sourcemap payload field "minifiedFile" is not empty
    And the sourcemap payload field "projectRoot" ends with "features/js/fixtures/js-webpack4"
    And the sourcemap payload field "overwrite" equals "true"

  Scenario: Searches in the dist folder automatically
    When I run bugsnag-cli with upload js --upload-api-root-url=http://localhost:9339 --api-key=1234567890ABCDEF1234567890ABCDEF --overwrite --base-url=example.com --project-root=features/js/fixtures/js-webpack4 features/js/fixtures/js-webpack4/
    And I wait to receive 1 sourcemaps
    Then the sourcemap is valid for the JS Build API
    Then the sourcemaps Content-Type header is valid multipart form-data
    And the sourcemap payload field "apiKey" equals "1234567890ABCDEF1234567890ABCDEF"
    And the sourcemap payload field "appVersion" equals "1.2.3"
    And the sourcemap payload field "minifiedUrl" equals "example.com/main.js"
    And the sourcemap payload field "sourceMap" is valid json
    And the sourcemap payload field "minifiedFile" is not empty
    And the sourcemap payload field "projectRoot" ends with "features/js/fixtures/js-webpack4"
    And the sourcemap payload field "overwrite" equals "true"

  Scenario: Uses the working directory as project root
    When I run bugsnag-cli with upload js --upload-api-root-url=http://localhost:9339 --api-key=1234567890ABCDEF1234567890ABCDEF --overwrite --base-url=example.com features/js/fixtures/js-webpack4/
    And I wait to receive 1 sourcemaps
    Then the sourcemap is valid for the JS Build API
    Then the sourcemaps Content-Type header is valid multipart form-data
    And the sourcemap payload field "apiKey" equals "1234567890ABCDEF1234567890ABCDEF"
    And the sourcemap payload field "appVersion" equals "1.2.3"
    And the sourcemap payload field "minifiedUrl" equals "example.com/main.js"
    And the sourcemap payload field "sourceMap" is valid json
    And the sourcemap payload field "minifiedFile" is not empty
    And the sourcemap payload field "projectRoot" ends with "js-webpack4"
    And the sourcemap payload field "overwrite" equals "true"

  Scenario: Base URL correctly appends the path
    When I run bugsnag-cli with upload js --upload-api-root-url=http://localhost:9339 --api-key=1234567890ABCDEF1234567890ABCDEF --overwrite --base-url=example.com --project-root=features/js/fixtures/js-webpack4 features/js/fixtures/js-webpack4/dist
    And I wait to receive 1 sourcemaps
    Then the sourcemap is valid for the JS Build API
    Then the sourcemaps Content-Type header is valid multipart form-data
    And the sourcemap payload field "apiKey" equals "1234567890ABCDEF1234567890ABCDEF"
    And the sourcemap payload field "appVersion" equals "1.2.3"
    And the sourcemap payload field "minifiedUrl" equals "example.com/main.js"
    And the sourcemap payload field "sourceMap" is valid json
    And the sourcemap payload field "minifiedFile" is not empty
    And the sourcemap payload field "projectRoot" ends with "features/js/fixtures/js-webpack4"
    And the sourcemap payload field "overwrite" equals "true"

  Scenario: Build and Upload js webpack4 sourcemaps
    When I make the "features/base-fixtures/js-webpack4"
    And I wait for the build to succeed

    When I run bugsnag-cli with upload js --upload-api-root-url=http://localhost:9339 --api-key=1234567890ABCDEF1234567890ABCDEF --overwrite --bundle-url=example.com --project-root=features/base-fixtures/js-webpack4 features/base-fixtures/js-webpack4/dist/main.js.map
    And I wait to receive 1 sourcemaps
    Then the sourcemap is valid for the JS Build API
    Then the sourcemaps Content-Type header is valid multipart form-data
    And the sourcemap payload field "apiKey" equals "1234567890ABCDEF1234567890ABCDEF"
    And the sourcemap payload field "appVersion" equals "1.2.3"
    And the sourcemap payload field "minifiedUrl" equals "example.com"
    And the sourcemap payload field "sourceMap" is valid json
    And the sourcemap payload field "minifiedFile" is not empty
    And the sourcemap payload field "projectRoot" ends with "features/base-fixtures/js-webpack4"
    And the sourcemap payload field "overwrite" equals "true"