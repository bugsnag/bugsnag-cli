Feature: dSYM Expected Error and Warning scenario Integration Tests

  Scenario: If --ignore-missing-dwarf is set to true, then the log message returned should be [WARN]
    When I run bugsnag-cli with upload dsym --upload-api-root-url=http://localhost:9339 --api-key=1234567890ABCDEF1234567890ABCDEF --project-root=/path/to/project/root --scheme=test --ignore-missing-dwarf=true --dsym-path=features/dsym/fixtures/empty-path features/dsym/fixtures/empty-path
    Then I should see a log level of "[WARN]" when no dSYM files could be found

  Scenario: If --ignore-missing-dwarf is not set, then the log message returned should be [ERROR]
    When I run bugsnag-cli with upload dsym --upload-api-root-url=http://localhost:9339 --api-key=1234567890ABCDEF1234567890ABCDEF --project-root=/path/to/project/root --scheme=test --dsym-path=features/dsym/fixtures/empty-path features/dsym/fixtures/empty-path
    Then I should see a log level of "[ERROR]" when no dSYM files could be found
    