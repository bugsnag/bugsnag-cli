require 'rbconfig'

os = RbConfig::CONFIG['host_os']
arch = RbConfig::CONFIG['host_cpu']

case
when os.downcase.include?('linux')
  os = 'linux'
when os.downcase.include?('darwin')
  os = 'macos'
end

When(/^I run bugsnag-cli on mac$/) do
  @output = `bin/#{arch}-#{os}/bugsnag-cli 2>&1`
end

When(/^I run bugsnag-cli with (.*)$/) do |flags|
  @output = `bin/#{arch}-#{os}/bugsnag-cli #{flags}`
end

Then(/^I should see the help banner$/) do
  run_output.include?("Usage: bugsnag-cli <command>")
end

Then(/^I should see the API Key error$/) do
  run_output.include?("[ERROR] no API key provided")
end

Then(/^I should see the missing path error$/) do
  run_output.include?("bugsnag-cli-arm64-darwin: error: expected \"<path>\"")
end

Then(/^I should see the no such file or directory error$/) do
  run_output.include?("bugsnag-cli-arm64-darwin: error: <path>: stat /path/to/no/file: no such file or directory")
end
