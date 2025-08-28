Feature: Breakpad Integration Tests

    Scenario: Upload a single breakpad .sym using all CLI flags
      When I run bugsnag-cli upload "breakpad" with the following arguments:
        | --upload-api-root-url                           | http://localhost:$MAZE_RUNNER_PORT               |
        | --api-key                                       | 1234567890ABCDEF1234567890ABCDEF                 |
        | --project-root                                  | /features/breakpad/fixtures/breakpad-symbols.sym |
        | --debug-file                                    | /features/breakpad/fixtures/breakpad-symbols.sym |
        | --code-file                                     | /features/breakpad/fixtures/breakpad-symbols.sym |
        | --version-name                                  | 2                                                |
        | --cpu-arch                                      | x86_64                                           |
        | --os-name                                       | Linux                                            |
        | --debug-identifier                              | 1234567890ABCDEF1234567890ABCDEF                 |
        | --product-name                                  | test-product                                     |
        | features/breakpad/fixtures/breakpad-symbols.sym |                                                  |
        And I wait to receive 1 sourcemaps
        Then the sourcemaps are valid for the API
        Then the sourcemaps Content-Type header is valid multipart form-data
        And the sourcemap "api_key" query parameter equals "1234567890ABCDEF1234567890ABCDEF"
        And the sourcemap "project_root" query parameter equals "/features/breakpad/fixtures/breakpad-symbols.sym"
        Then the sourcemap payload fields should be:
        | version           | 2                                                 |
        | os                | Linux                                             |
        | cpu               | x86_64                                            |
        | debug_file        | /features/breakpad/fixtures/breakpad-symbols.sym  |
        | code_file         | /features/breakpad/fixtures/breakpad-symbols.sym  |
        | debug_identifier  | 1234567890ABCDEF1234567890ABCDEF                  |
        | product           | test-product                                      |
