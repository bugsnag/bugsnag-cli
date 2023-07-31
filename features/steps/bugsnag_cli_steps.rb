require 'rbconfig'
require 'etc'

os = RbConfig::CONFIG['host_os']
arch = RbConfig::CONFIG['host_cpu']
user = Etc.getlogin
commit_hash = `git rev-parse HEAD`

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

When(/^I run bugsnag-cli on mac$/) do
  @output = `bin/#{arch}-#{os}-#{binary} 2>&1`
end

When(/^I run bugsnag-cli with (.*)$/) do |flags|
  @output = `bin/#{arch}-#{os}-#{binary} #{flags}`
end

Then(/^I should see the help banner$/) do
  run_output.include?("Usage: bugsnag-cli <command>")
end

Then(/^I should see the API Key error$/) do
  run_output.include?("[ERROR] no API key provided")
end

Then(/^I should see the missing path error$/) do
  run_output.include?("error: expected \"<path>\"")
  end

Then(/^I should see the missing app version error$/) do
  run_output.include?("[ERROR] Missing app version, please provide this via the command line options")
end

Then(/^I should see the no such file or directory error$/) do
  run_output.include?("error: <path>: stat /path/to/no/file: no such file or directory")
end

Then(/^the payload should match local information$/) do
  run_output.include?("\"appVersion\": \"1.2.3\"")
  run_output.include?("\"apiKey\": \"1234567890ABCDEF1234567890ABCDEF\"")
  run_output.include?("\"builderName\": \"#{user}\"")
  run_output.include?("\"revision\": \"#{commit_hash}\"")
  run_output.include?("\"repository\": \"git@github.com:bugsnag/bugsnag-cli\"")
end

Then('{string} should be used as {string}') do |value, field|
    Maze.check.include(run_output, "Using #{value} as #{field} from")
end
