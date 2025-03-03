Feature: Breakpad Integration Tests

    Scenario: Upload a single breakpad .sym using all CLI flags
        When I run bugsnag-cli with upload breakpad features/breakpad/fixtures/breakpad-symbols.sym --upload-api-root-url=http://localhost:9339 --version-code=2 --version-name=2.0 --os-name="Linux" --cpu-info="x86_64" --product="test-product" --api-key=1234567890ABCDEF1234567890ABCDEF --overwrite
        And I wait to receive 1 sourcemaps
        Then the sourcemaps Content-Type header is valid multipart form-data
        And the sourcemap payload field "api_Key" equals "1234567890ABCDEF1234567890ABCDEF"
        And the sourcemap payload field "OsName" equals "Linux, x86_64"
        And the sourcemap payload field "overwrite" equals "true"
        And the sourcemap payload field "debug_identifier" equals "CB77944DB22F6H3S00000000000000000"
        And the sourcemap payload field "cpuInfo" equals "x86_64"
        And the sourcemap payload field "versionCode" equals "2"
        And the sourcemap payload field "versionName" equals "2.0"
        And the sourcemap payload field "product" equals "test-product"
