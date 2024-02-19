Feature: dSYM Upload Integration Tests

  Scenario: Upload a single dSYM sourcemap using path containing one dSYM
    When I run bugsnag-cli with upload dsym --upload-api-root-url=http://localhost:9339 --api-key=1234567890ABCDEF1234567890ABCDEF --project-root=/my/project/root/ features/dsym/fixtures/single-dsym
    And I wait to receive 1 sourcemaps
    Then the sourcemap is valid for the dSYM Build API
    Then the sourcemaps Content-Type header is valid multipart form-data
    And the sourcemap payload field "apiKey" equals "1234567890ABCDEF1234567890ABCDEF"
    And the sourcemap payload field "platform" equals "ios"

  Scenario: Upload multiple dSYM sourcemaps using path containing multiple dSYMs
    When I run bugsnag-cli with upload dsym --upload-api-root-url=http://localhost:9339 --api-key=1234567890ABCDEF1234567890ABCDEF --project-root=/my/project/root/ features/dsym/fixtures/dsyms
    And I wait to receive 2 sourcemaps
    Then the sourcemap is valid for the dSYM Build API
    Then the sourcemaps Content-Type header is valid multipart form-data
    And the sourcemap payload field "apiKey" equals "1234567890ABCDEF1234567890ABCDEF"
    And the sourcemap payload field "platform" equals "ios"

  Scenario: Upload a single dSYM sourcemap using zip file containing one dSYM
    When I run bugsnag-cli with upload dsym --upload-api-root-url=http://localhost:9339 --api-key=1234567890ABCDEF1234567890ABCDEF --project-root=/my/project/root/ features/dsym/fixtures/single-dsym.zip
    And I wait to receive 1 sourcemaps
    Then the sourcemap is valid for the dSYM Build API
    Then the sourcemaps Content-Type header is valid multipart form-data
    And the sourcemap payload field "apiKey" equals "1234567890ABCDEF1234567890ABCDEF"
    And the sourcemap payload field "platform" equals "ios"

  Scenario: Upload multiple dSYM sourcemaps using zip file containing multiple dSYMs
    When I run bugsnag-cli with upload dsym --upload-api-root-url=http://localhost:9339 --api-key=1234567890ABCDEF1234567890ABCDEF --project-root=/my/project/root/ features/dsym/fixtures/dsyms.zip
    And I wait to receive 2 sourcemaps
    Then the sourcemap is valid for the dSYM Build API
    Then the sourcemaps Content-Type header is valid multipart form-data
    And the sourcemap payload field "apiKey" equals "1234567890ABCDEF1234567890ABCDEF"
    And the sourcemap payload field "platform" equals "ios"

  Scenario: Uploading a zip file containing directory of dSYM files that was compressed with macOS Archive Utility
    When I run bugsnag-cli with upload dsym --upload-api-root-url=http://localhost:9339 --api-key=1234567890ABCDEF1234567890ABCDEF --project-root=/my/project/root/ features/dsym/fixtures/macos-compressed-dsyms.zip
    And I wait to receive 2 sourcemaps
    Then the sourcemap is valid for the dSYM Build API
    Then the sourcemaps Content-Type header is valid multipart form-data
    And the sourcemap payload field "apiKey" equals "1234567890ABCDEF1234567890ABCDEF"
    And the sourcemap payload field "platform" equals "ios"

  Scenario: Upload symbols from an AppCenter zip
    When I run bugsnag-cli with upload dsym --upload-api-root-url=http://localhost:9339 --api-key=1234567890ABCDEF1234567890ABCDEF --project-root=/my/project/root/ "features/dsym/fixtures/app-center/symbols.zip"
    And I wait to receive 2 sourcemaps
    Then the sourcemap is valid for the dSYM Build API
    Then the sourcemaps Content-Type header is valid multipart form-data
    And the sourcemap payload field "apiKey" equals "1234567890ABCDEF1234567890ABCDEF"
    And the sourcemap payload field "platform" equals "ios"

  Scenario: Uploading an .xcarchive containing commas and special characters
    When I run bugsnag-cli with upload dsym --upload-api-root-url=http://localhost:9339 --api-key=1234567890ABCDEF1234567890ABCDEF --project-root=/my/project/root/ "features/dsym/fixtures/bugsnag-example 14-05-2021,,, 11.27éøœåñü#.xcarchive"
    And I wait to receive 2 sourcemaps
    Then the sourcemap is valid for the dSYM Build API
    Then the sourcemaps Content-Type header is valid multipart form-data
    And the sourcemap payload field "apiKey" equals "1234567890ABCDEF1234567890ABCDEF"
    And the sourcemap payload field "platform" equals "ios"

  Scenario: Attempt to upload a single dSYM sourcemap using path containing one dSYM but --project-root is not defined in the command
    When I run bugsnag-cli with upload dsym --upload-api-root-url=http://localhost:9339 --api-key=1234567890ABCDEF1234567890ABCDEF features/dsym/fixtures/single-dsym
    And I wait to receive 1 sourcemaps
    Then I should see the Project Root error