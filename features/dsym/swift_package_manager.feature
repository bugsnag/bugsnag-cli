Feature: dSYM Uploads for Swift Package Manager Projects Integration Tests

  Scenario: Upload a single dSYM sourcemap using all CLI flags
    When I run bugsnag-cli with upload dsym --upload-api-root-url=http://localhost:9339/dsym --api-key=1234567890ABCDEF1234567890ABCDEF --overwrite --dev --plist=features/base-fixtures/dsym/swift-package-manager/swift-package-manager/Info.plist --scheme=swift-package-manager --version-name=1.0 features/base-fixtures/dsym/swift-package-manager
    And I wait to receive 1 sourcemaps
    Then the sourcemap is valid for the dSYM Build API
    Then the sourcemaps Content-Type header is valid multipart form-data
    And the sourcemap payload field "apiKey" equals "1234567890ABCDEF1234567890ABCDEF"
    And the sourcemap payload field "appVersion" equals "1.0"
    And the sourcemap payload field "dev" equals "true"
    And the sourcemap payload field "platform" equals "ios"
    And the sourcemap payload field "overwrite" equals "true"

  Scenario: Upload a single dSYM sourcemap using only api-key and a path pointing to a SPM project
    When I run bugsnag-cli with upload dsym --upload-api-root-url=http://localhost:9339/dsym --api-key=1234567890ABCDEF1234567890ABCDEF features/base-fixtures/dsym/swift-package-manager
    And I wait to receive 1 sourcemaps
    Then the sourcemap is valid for the dSYM Build API
    Then the sourcemaps Content-Type header is valid multipart form-data
    And the sourcemap payload field "apiKey" equals "1234567890ABCDEF1234567890ABCDEF"
    And the sourcemap payload field "platform" equals "ios"

  Scenario: Upload a single dSYM sourcemap using only api-key and a path pointing to a xcodeproj directory of a SPM project
    When I run bugsnag-cli with upload dsym --upload-api-root-url=http://localhost:9339/dsym --api-key=1234567890ABCDEF1234567890ABCDEF features/base-fixtures/dsym/swift-package-manager/swift-package-manager.xcodeproj
    And I wait to receive 1 sourcemaps
    Then the sourcemap is valid for the dSYM Build API
    Then the sourcemaps Content-Type header is valid multipart form-data
    And the sourcemap payload field "apiKey" equals "1234567890ABCDEF1234567890ABCDEF"
    And the sourcemap payload field "platform" equals "ios"

  Scenario: Upload a single dSYM sourcemap with scheme defined in command
    When I run bugsnag-cli with upload dsym --upload-api-root-url=http://localhost:9339/dsym --scheme=swift-package-manager --api-key=1234567890ABCDEF1234567890ABCDEF features/base-fixtures/dsym/swift-package-manager
    And I wait to receive 1 sourcemaps
    Then the sourcemap is valid for the dSYM Build API
    Then the sourcemaps Content-Type header is valid multipart form-data
    And the sourcemap payload field "apiKey" equals "1234567890ABCDEF1234567890ABCDEF"
    And the sourcemap payload field "platform" equals "ios"

  Scenario: Upload a single dSYM sourcemap using plist to define apiKey and appVersion
    When I run bugsnag-cli with upload dsym --upload-api-root-url=http://localhost:9339/dsym --plist=features/base-fixtures/dsym/swift-package-manager/swift-package-manager/Info.plist features/base-fixtures/dsym/swift-package-manager
    And I wait to receive 1 sourcemaps
    Then the sourcemap is valid for the dSYM Build API
    Then the sourcemaps Content-Type header is valid multipart form-data
    And the sourcemap payload field "apiKey" equals "1234567890ABCDEF1234567890ABCDEF"
    And the sourcemap payload field "appVersion" equals "1.0"
    And the sourcemap payload field "platform" equals "ios"