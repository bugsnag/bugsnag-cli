Feature: Upload Xcode Archives
  Scenario: Uploading an .xcarchive containing commas and special characters
    When I run bugsnag-cli with upload xcode-archive --upload-api-root-url=http://localhost:9339 --api-key=1234567890ABCDEF1234567890ABCDEF "features/xcode/fixtures/bugsnag-example 14-05-2021,,, 11.27éøœåñü#.xcarchive"
    And I wait to receive 2 sourcemaps
    Then the sourcemap is valid for the dSYM Build API
    Then the sourcemaps Content-Type header is valid multipart form-data
    And the sourcemap payload field "apiKey" equals "1234567890ABCDEF1234567890ABCDEF"

  Scenario: Archive and upload an .xcarchive
    When I make the "features/base-fixtures/dsym/archive"
    Then I wait for the build to succeed


    When I run bugsnag-cli with upload xcode-archive --upload-api-root-url=http://localhost:9339 --api-key=1234567890ABCDEF1234567890ABCDEF features/base-fixtures/dsym/
    And I wait to receive 1 sourcemaps
    Then the sourcemap is valid for the dSYM Build API
    Then the sourcemaps Content-Type header is valid multipart form-data
    And the sourcemap payload field "apiKey" equals "1234567890ABCDEF1234567890ABCDEF"
