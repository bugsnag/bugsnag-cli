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

base_dir = Dir.pwd

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
  @output = `npm pack`
  @bugsnag_cli_package_path = "#{base_dir}/#{@output}"
end

When('I install the bugsnag-cli via {string} in a new directory') do |package_manager|
  @fixture_dir = "#{base_dir}/features/cli/fixtures/#{package_manager}"
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

  Dir.chdir(base_dir)
end

Then('the {string} directory should contain {string}') do |directory, package|
  Maze.check.include(`ls #{@fixture_dir}/#{directory}`, package)
end

Given('I build the Unity Android example project') do
  @fixture_dir = "#{base_dir}/platforms-examples/Unity"
  Dir.chdir(@fixture_dir)
  @output = `./build_android.sh aab`
  Dir.chdir(base_dir)
end

And('I wait for the Unity symbols to generate') do
  Maze.check.include(`ls #{@fixture_dir}`, 'UnityExample-1.0-v1-IL2CPP.symbols.zip')
end

Given(/^I set the NDK path to the Unity bundled version$/) do
#  Set the environment variable to the path of the NDK bundled with Unity
  ENV['ANDROID_NDK_ROOT'] = "/Applications/Unity/Hub/Editor/#{ENV['UNITY_VERSION']}/PlaybackEngines/AndroidPlayer/NDK"
end

# dSYM
Before('@CleanAndBuildDsym') do
  scheme = 'dSYM-Example'
  project_path = 'features/base-fixtures/dsym'

  # Find the Xcode archive path dynamically
  custom_archives_path = `defaults read com.apple.dt.Xcode IDECustomDistributionArchivesLocation`.strip
  archives_path = custom_archives_path.empty? ? File.expand_path("~/Library/Developer/Xcode/Archives/") : custom_archives_path
  today = Date.today.strftime('%Y-%m-%d')

  # Delete archives for the given scheme created today
  Dir.glob(File.join(archives_path, "#{today}*")) do |archive|
    if archive.include?(scheme)
      puts "Removing archive: #{archive}"
      FileUtils.rm_rf(archive)
    end
  end

  # Clear Xcode build directories matching wildcard pattern
  build_paths = Dir.glob(File.expand_path("~/Library/Developer/Xcode/DerivedData/dSYM-Example-*"))

  if build_paths.any?
    build_paths.each do |path|
      puts "Clearing Xcode build directory: #{path}"
      FileUtils.rm_rf(path)
    end
  else
    puts "No matching build directories found."
  end

  # Build the project
  puts "Building project: #{project_path}"
  @output = `make features/base-fixtures/dsym`
end

Before('@CleanAndArchiveDsym') do
  scheme = 'dSYM-Example'
  project_path = 'features/base-fixtures/dsym'

  # Find the Xcode archive path dynamically
  custom_archives_path = `defaults read com.apple.dt.Xcode IDECustomDistributionArchivesLocation`.strip
  archives_path = custom_archives_path.empty? ? File.expand_path("~/Library/Developer/Xcode/Archives/") : custom_archives_path
  today = Date.today.strftime('%Y-%m-%d')

  # Delete archives for the given scheme created today
  Dir.glob(File.join(archives_path, "#{today}*")) do |archive|
    if archive.include?(scheme)
      puts "Removing archive: #{archive}"
      FileUtils.rm_rf(archive)
    end
  end

  # Clear Xcode build directories matching wildcard pattern
  build_paths = Dir.glob(File.expand_path("~/Library/Developer/Xcode/DerivedData/dSYM-Example-*"))

  if build_paths.any?
    build_paths.each do |path|
      puts "Clearing Xcode build directory: #{path}"
      FileUtils.rm_rf(path)
    end
  else
    puts "No matching build directories found."
  end

  # Build the project
  puts "Building project: #{project_path}"
  @output = `make features/base-fixtures/dsym/archive`
end

# React Native
Before('@BuildRNAndroid') do
  unless defined?($setup_android) && $setup_android
    puts "Setting up React Native Android app and sourcemap..."
    @output = `node features/react-native/scripts/generate.js`
    Maze.check.include(`ls features/react-native/fixtures/generated/old-arch/#{ENV['RN_VERSION']}/android/app/build/generated/sourcemaps/react/release`, 'index.android.bundle.map')

    ENV['APP_MANIFEST_PATH'] = "features/react-native/fixtures/generated/old-arch/#{ENV['RN_VERSION']}/android/app/build/intermediates/merged_manifests/release/AndroidManifest.xml"
    ENV['BUNDLE_PATH'] = "features/react-native/fixtures/generated/old-arch/#{ENV['RN_VERSION']}/android/app/build/generated/assets/createBundleReleaseJsAndAssets/index.android.bundle"
    ENV['SOURCE_MAP_PATH'] = "features/react-native/fixtures/generated/old-arch/#{ENV['RN_VERSION']}/android/app/build/generated/sourcemaps/react/release/index.android.bundle.map"

    if ENV['RN_VERSION'].to_f == 0.70
      ENV['BUNDLE_PATH'] = "features/react-native/fixtures/generated/old-arch/#{ENV['RN_VERSION']}/android/app/build/generated/assets/react/release/index.android.bundle"
    end

    if ENV['RN_VERSION'].to_f == 0.71
      ENV['BUNDLE_PATH'] = "features/react-native/fixtures/generated/old-arch/#{ENV['RN_VERSION']}/android/app/build/ASSETS/createBundleReleaseJsAndAssets/index.android.bundle"
    end

    if ENV['RN_VERSION'].to_f == 0.75
      ENV['APP_MANIFEST_PATH'] = "features/react-native/fixtures/generated/old-arch/#{ENV['RN_VERSION']}/android/app/build/intermediates/merged_manifests/release/processReleaseManifest/AndroidManifest.xml"
    end

    $setup_android = true
  end
end

Before('@BuildRNiOS') do
  unless defined?($setup_ios) && $setup_ios
    puts "Setting up React Native iOS app and sourcemap..."
    @output = `node features/react-native/scripts/generate.js`
    Maze.check.include(`ls features/react-native/fixtures/generated/old-arch/#{ENV['RN_VERSION']}/ios/build/sourcemaps`, 'main.jsbundle.map')
    $setup_ios = true
  end
end

Before('@BuildExportRNiOS') do
  unless defined?($export_ios) && $export_ios
    puts "Setting up React Native iOS app and sourcemap and exporting the archive..."
    @output = `EXPORT_ARCHIVE=true node features/react-native/scripts/generate.js`
    Maze.check.include(`ls features/react-native/fixtures/generated/old-arch/#{ENV['RN_VERSION']}/ios/build/sourcemaps`, 'main.jsbundle.map')
    $export_ios = true
  end
end
