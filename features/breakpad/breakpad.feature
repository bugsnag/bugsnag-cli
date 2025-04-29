Feature: Breakpad Integration Tests

    Scenario: Upload a single breakpad .sym using all CLI flags
        When I run bugsnag-cli with upload breakpad features/breakpad/fixtures/breakpad-symbols.sym --upload-api-root-url=http://localhost:$MAZE_RUNNER_PORT --project-root="/features/breakpad/fixtures/breakpad-symbols.sym" --debug-file="/features/breakpad/fixtures/breakpad-symbols.sym" --code-file="/features/breakpad/fixtures/breakpad-symbols.sym" --version-name="2" --cpu-arch="x86_64" --os-name="Linux" --debug-identifier="1234567890ABCDEF1234567890ABCDEF" --product-name="test-product" --api-key=1234567890ABCDEF1234567890ABCDEF --overwrite
        And I wait to receive 1 sourcemaps
        Then the sourcemap is valid for the Breakpad Build API
        Then the sourcemaps Content-Type header is valid multipart form-data
        And the sourcemap "api_key" query parameter equals "1234567890ABCDEF1234567890ABCDEF"
        And the sourcemap "project_root" query parameter equals "/features/breakpad/fixtures/breakpad-symbols.sym"
        And the sourcemap "overwrite" query parameter equals "true"
        And the sourcemap payload field "version" equals "2"
        And the sourcemap payload field "os" equals "Linux"
        And the sourcemap payload field "cpu" equals "x86_64"
        And the sourcemap payload field "debug_file" equals "/features/breakpad/fixtures/breakpad-symbols.sym"
        And the sourcemap payload field "code_file" equals "/features/breakpad/fixtures/breakpad-symbols.sym"
        And the sourcemap payload field "debug_identifier" equals "1234567890ABCDEF1234567890ABCDEF"
        And the sourcemap payload field "product" equals "test-product"
