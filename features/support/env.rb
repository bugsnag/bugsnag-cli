BeforeAll do
  $api_key = '1234567890ABCDEF1234567890ABCDEF'
  ENV['MAZE_RUNNER_PORT'] ||= '9339'
end

def run_output
  @output ||= ''
end
