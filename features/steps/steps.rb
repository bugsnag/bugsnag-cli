require 'rbconfig'
require 'etc'
require 'digest/md5'

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
    And the requests are different
  )
end

Then('the sourcemap is valid for the Breakpad Build API') do
  steps %(
    And the sourcemap "api_key" query parameter equals "#{$api_key}"
    And the sourcemap "project_root" query parameter is not null
    And the requests are different
  )
end

Then('the sourcemap is valid for the React Native Build API') do
  steps %(
    And the sourcemap payload field "apiKey" equals "#{$api_key}"
    And the sourcemap payload field "appVersion" is not null
    And the requests are different
  )
end

Then('the sourcemap is valid for the JS Build API') do
  steps %(
    And the sourcemap payload field "apiKey" equals "#{$api_key}"
    And the sourcemap payload field "appVersion" is not null
    And the requests are different
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
    And the requests are different
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
  # Change directory to the js directory
  @js_dir = "#{base_dir}/js"
  Dir.chdir(@js_dir)
  @output = `npm i && npm pack`
  @bugsnag_cli_package_path = "#{@js_dir}/#{@output}"
  Dir.chdir(base_dir)
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
def clean_and_build(scheme, project_path, build_target)
  puts "🚀 Starting dSYM cleaning and build process..."

  # Find the Xcode archive path dynamically
  custom_archives_path = `defaults read com.apple.dt.Xcode IDECustomDistributionArchivesLocation`.strip
  archives_path = custom_archives_path.empty? ? File.expand_path("~/Library/Developer/Xcode/Archives/") : custom_archives_path
  today = Date.today.strftime('%d-%m-%Y')
  archives_path = File.join(archives_path, Date.today.strftime('%Y-%m-%d'))

  # Delete archives for the given scheme created today
  archives_deleted = false
  Dir.glob(File.join(archives_path, "#{scheme} #{today}*")) do |archive|
    if archive.include?(scheme)
      puts "🗑️ Removing archive: #{archive}"
      FileUtils.rm_rf(archive)
      archives_deleted = true
    end
  end
  puts archives_deleted ? "✅ Archives deleted successfully." : "ℹ️ No matching archives found."

  # Clear Xcode build directories matching wildcard pattern
  build_paths = Dir.glob(File.expand_path("~/Library/Developer/Xcode/DerivedData/#{scheme}-*"))

  if build_paths.any?
    build_paths.each do |path|
      puts "🧹 Clearing Xcode build directory: #{path}"
      FileUtils.rm_rf(path)
    end
    puts "✅ Xcode build directories cleared."
  else
    puts "ℹ️ No matching build directories found."
  end

  # Build the project
  puts "🏗️ Building project: #{project_path}"
  output = `make #{build_target}`

  if $?.success?
    puts "✅ Build completed successfully."
  else
    puts "❌ Build failed. Output:\n#{output}"
    raise "Build process failed."
  end
end

Before('@CleanAndBuildDsym') do
  scheme = 'dSYM-Example'
  project_path = 'features/base-fixtures/dsym'
  build_target = project_path

  clean_and_build(scheme, project_path, build_target)
end

Before('@CleanAndArchiveDsym') do
  scheme = 'dSYM-Example'
  project_path = 'features/base-fixtures/dsym'
  build_target = "#{project_path}/archive"

  clean_and_build(scheme, project_path, build_target)
end

# React Native
Before('@BuildRNAndroid') do
  unless defined?($setup_android) && $setup_android
    puts "🚀 Setting up React Native Android app and generating sourcemap..."

    generate_command = 'node features/react-native/scripts/generate.js'
    output = `#{generate_command}`

    if $?.success?
      puts "✅ Android setup completed successfully."
    else
      puts "❌ Android setup failed. Output:\n#{output}"
      raise "Failed to set up React Native Android."
    end

    base_path = "features/react-native/fixtures/generated/old-arch/#{ENV['RN_VERSION']}/android/app/build"
    sourcemap_path = "#{base_path}/generated/sourcemaps/react/release"

    Maze.check.include(`ls #{sourcemap_path}`, 'index.android.bundle.map')

    puts "📍 Sourcemap verified successfully."

    # Set environment variables
    ENV['APP_MANIFEST_PATH'] = "#{base_path}/intermediates/merged_manifests/release/AndroidManifest.xml"
    ENV['BUNDLE_PATH'] = "#{base_path}/generated/assets/createBundleReleaseJsAndAssets/index.android.bundle"
    ENV['SOURCE_MAP_PATH'] = "#{sourcemap_path}/index.android.bundle.map"

    case ENV['RN_VERSION'].to_f
    when 0.70
      ENV['BUNDLE_PATH'] = "#{base_path}/generated/assets/react/release/index.android.bundle"
    when 0.71
      ENV['BUNDLE_PATH'] = "#{base_path}/ASSETS/createBundleReleaseJsAndAssets/index.android.bundle"
    when 0.75
      ENV['APP_MANIFEST_PATH'] = "#{base_path}/intermediates/merged_manifests/release/processReleaseManifest/AndroidManifest.xml"
    end

    puts "🔧 Environment variables set:"
    puts "  - APP_MANIFEST_PATH: #{ENV['APP_MANIFEST_PATH']}"
    puts "  - BUNDLE_PATH: #{ENV['BUNDLE_PATH']}"
    puts "  - SOURCE_MAP_PATH: #{ENV['SOURCE_MAP_PATH']}"

    $setup_android = true
  end
end

Before('@BuildRNiOS') do
  unless defined?($setup_ios) && $setup_ios
    puts "🚀 Setting up React Native iOS app and generating sourcemap..."

    generate_command = 'node features/react-native/scripts/generate.js'
    output = `#{generate_command}`

    if $?.success?
      puts "✅ React Native iOS setup completed successfully."
    else
      puts "❌ Setup failed. Output:\n#{output}"
      raise "Failed to set up React Native iOS."
    end

    sourcemap_path = "features/react-native/fixtures/generated/old-arch/#{ENV['RN_VERSION']}/ios/build/sourcemaps"

    Maze.check.include(`ls #{sourcemap_path}`, 'main.jsbundle.map')

    puts "📍 Sourcemap verified successfully."
    $setup_ios = true
  end
end

Before('@BuildExportRNiOS') do
  unless defined?($export_ios) && $export_ios
    puts "🚀 Setting up React Native iOS app, generating sourcemaps, and exporting the archive..."

    export_command = 'EXPORT_ARCHIVE=true node features/react-native/scripts/generate.js'
    output = `#{export_command}`

    if $?.success?
      puts "✅ Archive export completed successfully."
    else
      puts "❌ Archive export failed. Output:\n#{output}"
      raise "Failed to export archive."
    end

    sourcemap_path = "features/react-native/fixtures/generated/old-arch/#{ENV['RN_VERSION']}/ios/build/sourcemaps"

    Maze.check.include(`ls #{sourcemap_path}`, 'main.jsbundle.map')

    puts "📍 Sourcemap verified successfully."
    $export_ios = true
  end
end

Then('I should see the {string} in the output') do |log_message|
  Maze.check.include(run_output, log_message)
end

Before('@BuildNestedJS') do
  unless defined?($nested_js) && $nested_js
    puts "🚀 Building Nested JS fixture and generating sourcemaps..."

    @fixture_dir= "#{base_dir}/features/base-fixtures/js"

    # Change to the Nested JS fixture directory
    Dir.chdir(@fixture_dir)

    # Run NPM install
    npm_install = `npm install`
    if $?.success?
      puts "✅ NPM install completed successfully."
    else
      puts "❌ NPM install failed. Output:\n#{npm_install}"
      raise "Failed to install NPM dependencies."
    end

    # Run the build script
    build_command = 'npm run build'
    output = `#{build_command}`

    if $?.success?
      puts "✅ Build completed successfully."
    else
      puts "❌ Build failed. Output:\n#{output}"
      raise "Failed to build Nested JS."
    end

    Maze.check.include(`ls #{@fixture_dir}`, 'out')

    puts "📍 Sourcemap verified successfully."

    # Change back to the base directory
    Dir.chdir(base_dir)
    $nested_js = true
  end
end

And(/^the builds payload field "([^"]*)" hash equals \{"([^"]*)"=>"([^"]*)", "([^"]*)"=>"([^"]*)"\}$/) do |arg1, arg2, arg3, arg4, arg5|
  builds = Maze::Server.builds.current[:body]
  expected_hash = { arg2 => arg3, arg4 => arg5 }

  Maze.check.equal(builds[arg1], expected_hash, "Expected builds payload field '#{arg1}' to equal #{expected_hash}, but got #{builds[arg1]}")
end

Then('the requests are different') do
  requests = Maze::Server.sourcemaps.remaining
  last_md5 = nil
  requests.each do |request|
    puts "last md5: #{last_md5}"
    body = request[:body]
    body = body.to_s
    current_md5 = Digest::MD5.hexdigest(body)
    Maze.check.not_equal(last_md5, current_md5) unless last_md5.nil?
    last_md5 = current_md5
  end
end
