Feature: dSYM Expected Error and Warning scenario Integration Tests

  Scenario: If --ignore-empty-dsym is set to true, then the log message returned should be [WARN]
    When I run bugsnag-cli with upload xcode-build --upload-api-root-url=http://localhost:9339 --api-key=1234567890ABCDEF1234567890ABCDEF --project-root=/path/to/project/root --scheme=test --ignore-empty-dsym=true features/xcode/fixtures/ZeroByteDsym
    Then I should see a log level of "[FATAL]" when no dSYM files could be found

  Scenario: If --ignore-empty-dsym is not set, then the log message returned should be [ERROR]
    When I run bugsnag-cli with upload xcode-build --upload-api-root-url=http://localhost:9339 --api-key=1234567890ABCDEF1234567890ABCDEF --project-root=/path/to/project/root --scheme=test features/xcode/fixtures/ZeroByteDsym
    Then I should see a log level of "[FATAL]" when no dSYM files could be found

  Scenario: If --ignore-missing-dwarf is set to true, then the log message returned should be [WARN]
    When I run bugsnag-cli with upload xcode-build --upload-api-root-url=http://localhost:9339 --api-key=1234567890ABCDEF1234567890ABCDEF --project-root=/path/to/project/root --scheme=test --ignore-missing-dwarf=true features/xcode/fixtures/MissingDWARFdSYM
    Then I should see a log level of "[FATAL]" when no dSYM files could be found

  Scenario: If --ignore-missing-dwarf is not set, then the log message returned should be [ERROR]
    When I run bugsnag-cli with upload xcode-build --upload-api-root-url=http://localhost:9339 --api-key=1234567890ABCDEF1234567890ABCDEF --project-root=/path/to/project/root --scheme=test features/xcode/fixtures/MissingDWARFdSYM
    Then I should see a log level of "[FATAL]" when no dSYM files could be found

  Scenario: If --ignore-missing-dwarf is set to true, then the log message returned should be [WARN]
    When I run bugsnag-cli with upload xcode-build --upload-api-root-url=http://localhost:9339 --api-key=1234567890ABCDEF1234567890ABCDEF --project-root=/path/to/project/root --scheme=test --ignore-missing-dwarf=true features/xcode/fixtures/MissingDWARFdSYM
    Then I should see a log level of "[FATAL]" when no dSYM files could be found

  Scenario: If --ignore-missing-dwarf is not set, then the log message returned should be [ERROR]
    When I run bugsnag-cli with upload xcode-build --upload-api-root-url=http://localhost:9339 --api-key=1234567890ABCDEF1234567890ABCDEF --project-root=/path/to/project/root --scheme=test features/xcode/fixtures/MissingDWARFdSYM
    Then I should see a log level of "[FATAL]" when no dSYM files could be found
