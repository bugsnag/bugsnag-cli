require 'rbconfig'
require 'etc'

os = RbConfig::CONFIG['host_os']
arch = RbConfig::CONFIG['host_cpu']

case
when os.downcase.include?('windows_nt'), ENV['WSL_DISTRO_NAME'] != nil
  os = 'windows'
  binary = 'bugsnag-cli.exe'
when os.downcase.include?('linux')
  os = 'linux'
  binary = 'bugsnag-cli'
when os.downcase.include?('darwin')
  os = 'macos'
  binary = 'bugsnag-cli'
end

When('I run bugsnag-cli') do
  @output = `bin/#{arch}-#{os}-#{binary} 2>&1`
end

When(/^I run bugsnag-cli with (.*)$/) do |flags|
  @output = `bin/#{arch}-#{os}-#{binary} #{flags} 2>&1`
  puts @output
end

Then('I should see a log level of {string} when no dSYM files could be found') do |log_level|
  message = log_level + ' No dSYM files found'
  Maze.check.include(run_output, message)
end

Then('I should see a log level of {string} when no dSYM files could be uploaded') do |log_level|
  message = log_level + ' failed after'
  Maze.check.include(run_output, message)
end

Then('I should see the help banner') do
  Maze.check.include(run_output, "Usage: #{arch}-#{os}-#{binary} <command>")
end

Then('I should see the API Key error') do
  Maze.check.include(run_output, "[FATAL] missing api key, please specify using `--api-key`")
end

Then('I should see the Project Root error') do
  Maze.check.include(run_output, "[FATAL] --project-root is required when uploading dSYMs from a directory that is not an Xcode project or workspace")
end


Then('I should see the missing path error') do
  Maze.check.include(run_output, "error: expected \"<path>\"")
end

Then('I should see the missing app version error') do
  Maze.check.include(run_output, "[FATAL] missing app version, please specify using `--version-name`")
end

Then('I should see the no such file or directory error') do
  Maze.check.include(run_output, "error: <path>: stat /path/to/no/file: no such file or directory")
end

Then('I should see the not an accepted value for the source control provider error') do
  Maze.check.include(run_output, "is not an accepted value for the source control provider. Accepted values are: github, github-enterprise, bitbucket, bitbucket-server, gitlab, gitlab-onpremise")
end

Then('I should see the missing source control provider error') do
  Maze.check.include(run_output, "error: --provider: missing source control provider, please specify using `--provider`. Accepted values are: github, github-enterprise, bitbucket, bitbucket-server, gitlab, gitlab-onpremise")
end

Then('I should see the path ambiguous error') do
  Maze.check.include(run_output, "Path ambiguous: more than one AAB file was found")
end

Then('the sourcemap is valid for the Proguard Build API') do
  steps %(
    Then the sourcemap is valid for the Android Build API
  )
end

Then('the sourcemap is valid for the NDK Build API') do
  steps %(
    Then the sourcemap is valid for the Android Build API
  )
end

Then('the sourcemap is valid for the Dart Build API') do
  steps %(
    And the sourcemap payload field "apiKey" equals "#{$api_key}"
    And the sourcemap payload field "buildId" is not null
  )
end

Then('the sourcemap is valid for the React Native Build API') do
  steps %(
    And the sourcemap payload field "apiKey" equals "#{$api_key}"
    And the sourcemap payload field "appVersion" is not null
  )
end

Then('the sourcemap is valid for the JS Build API') do
  steps %(
    And the sourcemap payload field "apiKey" equals "#{$api_key}"
    And the sourcemap payload field "appVersion" is not null
  )
end

Then('the sourcemap is valid for the dSYM Build API') do
  steps %(
    And the sourcemap payload field "apiKey" equals "#{$api_key}"
  )
end

Then('the sourcemap is valid for the Android Build API') do
  steps %(
    And the sourcemap payload field "apiKey" equals "#{$api_key}"
    And the sourcemap payload field "appId" is not null
  )
end

Then('the build is valid for the Builds API') do
  steps %(
    And the build payload field "apiKey" equals "#{$api_key}"
    And the build payload field "appVersion" is not null
  )
end

Then('the sourcemaps Content-Type header is valid multipart form-data') do
  expected = /^multipart\/form-data; boundary=([^;]+)/
  actual = Maze::Server.sourcemaps.current[:request]['content-type']
  Maze.check.match(expected, actual)
end

Then('the sourcemap payload field "sourceMap" is valid json') do
  require 'json'
  decoded = JSON.parse(Maze::Server.sourcemaps.current[:body]['sourceMap'])
  Maze.check.not_equal(decoded['mappings'].length, 0)
end

Then('the sourcemap payload field "minifiedFile" is not empty') do
  Maze.check.not_equal(Maze::Server.sourcemaps.current[:body]['minifiedFile'].length, 0)
end

Then('{string} should be used as {string}') do |value, field|
  Maze.check.include(run_output, "Using #{value} as #{field} from")
end

Then('I should see the build payload') do
  Maze.check.include(run_output,      "[INFO] (dryrun) Build payload:\n" +
    "{\n" +
    "    \"apiKey\": \"1234567890ABCDEF1234567890ABCDEF\",\n" +
    "    \"appVersionCode\": \"1\",\n" +
    "    \"sourceControl\": {\n")
end

def get_version_number(file_path)
  package_version = nil

  file_content = File.read(file_path)

  file_content.each_line do |line|
    if line =~ /\bpackage_version\s*=\s*(['"])(.*?)\1/
      package_version = $2
      break
    end
  end

  package_version
end

Then(/^the version number should match the version set in main\.go$/) do
  version_number = get_version_number "main.go"
  Maze.check.include(run_output, version_number)
end

And(/^I wait for the build to succeed$/) do
  Maze.check.not_include(run_output, "Error 1")
end

When(/^I make the "([^"]*)"$/) do |arg|
  @output = `make #{arg} 2>&1`
end

Then(/^I should only see the fatal log level messages$/) do
  Maze.check.include(run_output, "[FATAL]")
  Maze.check.not_include(run_output, "[ERROR]")
  Maze.check.not_include(run_output, "[WARN]")
  Maze.check.not_include(run_output, "[INFO]")
  Maze.check.not_include(run_output, "[DEBUG]")
end

Before('@installation') do
  @base_dir = Dir.pwd
  @output = `npm pack`
  @bugsnag_cli_package_path = "#{@base_dir}/#{@output}"
end

When('I install the bugsnag-cli via {string} in a new directory') do |package_manager|
  @fixture_dir = "#{@base_dir}/features/cli/fixtures/#{package_manager}"
  Dir.mkdir(@fixture_dir)
  Dir.chdir(@fixture_dir)

  case package_manager
  when 'npm'
    @init_output = `npm init -y`
    @install_output = `npm install #{@bugsnag_cli_package_path}`
  when 'yarn'
    @init_output = `yarn init -y`
    @install_output = `yarn add #{@bugsnag_cli_package_path}`
  when 'pnpm'
    @init_output = `pnpm init -y`
    @install_output = `pnpm add #{@bugsnag_cli_package_path}`
  end

  Dir.chdir(@base_dir)
end

Then('the {string} directory should contain {string}') do |directory, package|
  Maze.check.include(`ls #{@fixture_dir}/#{directory}`, package)
end

Given('I build the Unity Android example project') do
  @base_dir = Dir.pwd
  @fixture_dir = "#{@base_dir}/platforms-examples/Unity"
  Dir.chdir(@fixture_dir)
  @output = `./build_android.sh aab`
  Dir.chdir(@base_dir)
end

And('I wait for the Unity symbols to generate') do
  Maze.check.include(`ls #{@fixture_dir}`, 'UnityExample-1.0-v1-IL2CPP.symbols.zip')
end

Given(/^I set the NDK path to the Unity bundled version$/) do
#  Set the environment variable to the path of the NDK bundled with Unity
  ENV['ANDROID_NDK_ROOT'] = "/Applications/Unity/Hub/Editor/#{ENV['UNITY_VERSION']}/PlaybackEngines/AndroidPlayer/NDK"
end
