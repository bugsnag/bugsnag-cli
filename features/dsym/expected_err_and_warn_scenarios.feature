Feature: dSYM Expected Error and Warning scenario Integration Tests

  Scenario: If --ignore-empty-dsym is set to true, then the log message returned should be [WARN]
    When I run bugsnag-cli with upload dsym --upload-api-root-url=http://localhost:9339 --api-key=1234567890ABCDEF1234567890ABCDEF --project-root=/path/to/project/root --scheme=test --ignore-empty-dsym=true features/dsym/fixtures/ZeroByteDsym
    Then I should see a log level of "[WARN]" when no dSYM files could be found

  Scenario: If --ignore-empty-dsym is not set, then the log message returned should be [ERROR]
    When I run bugsnag-cli with upload dsym --upload-api-root-url=http://localhost:9339 --api-key=1234567890ABCDEF1234567890ABCDEF --project-root=/path/to/project/root --scheme=test features/dsym/fixtures/ZeroByteDsym
    Then I should see a log level of "[ERROR]" when no dSYM files could be found

  Scenario: If --ignore-missing-dwarf is set to true, then the log message returned should be [WARN]
    When I run bugsnag-cli with upload dsym --upload-api-root-url=http://localhost:9339 --api-key=1234567890ABCDEF1234567890ABCDEF --project-root=/path/to/project/root --scheme=test --ignore-missing-dwarf=true features/dsym/fixtures/MissingDWARFdSYM
    Then I should see a log level of "[WARN]" when no dSYM files could be found

  Scenario: If --ignore-missing-dwarf is not set, then the log message returned should be [ERROR]
    When I run bugsnag-cli with upload dsym --upload-api-root-url=http://localhost:9339 --api-key=1234567890ABCDEF1234567890ABCDEF --project-root=/path/to/project/root --scheme=test features/dsym/fixtures/MissingDWARFdSYM
    Then I should see a log level of "[ERROR]" when no dSYM files could be found

  Scenario: If --ignore-missing-dwarf is set to true, then the log message returned should be [WARN]
    When I run bugsnag-cli with upload dsym --upload-api-root-url=http://localhost:9339 --api-key=1234567890ABCDEF1234567890ABCDEF --project-root=/path/to/project/root --scheme=test --ignore-missing-dwarf=true features/dsym/fixtures/MissingDWARFdSYM
    Then I should see a log level of "[WARN]" when no dSYM files could be found

  Scenario: If --ignore-missing-dwarf is not set, then the log message returned should be [ERROR]
    When I run bugsnag-cli with upload dsym --upload-api-root-url=http://localhost:9339 --api-key=1234567890ABCDEF1234567890ABCDEF --project-root=/path/to/project/root --scheme=test features/dsym/fixtures/MissingDWARFdSYM
    Then I should see a log level of "[ERROR]" when no dSYM files could be found

  Scenario: If --fail-on-upload is set to true, then the log message returned should be [ERROR]
    When I run bugsnag-cli with upload dsym --upload-api-root-url=http://localhost:9333 --api-key=1234567890ABCDEF1234567890ABCDEF --project-root=/path/to/project/root --scheme=test --fail-on-upload=true features/dsym/fixtures/dsyms.zip
    Then I should see a log level of "[ERROR]" when no dSYM files could be uploaded

  Scenario: If --fail-on-upload is not set, then the log message returned should be [WARN]
    When I run bugsnag-cli with upload dsym --upload-api-root-url=http://localhost:9333 --api-key=1234567890ABCDEF1234567890ABCDEF --project-root=/path/to/project/root --scheme=test features/dsym/fixtures/dsyms.zip
    Then I should see a log level of "[WARN]" when no dSYM files could be uploaded
