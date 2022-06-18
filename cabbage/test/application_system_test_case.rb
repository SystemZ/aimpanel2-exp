require "test_helper"

class ApplicationSystemTestCase < ActionDispatch::SystemTestCase
  driven_by :selenium, using: :firefox, screen_size: [1400, 1400]
  # using: :chrome

  # https://stackoverflow.com/questions/45483353/cannot-remove-port-from-capybara-while-visiting-external-link
  def visit_external(url)
    Capybara.always_include_port = false
    visit url
    Capybara.always_include_port = true
  end

end
