require "application_system_test_case"

class NamespacesTest < ApplicationSystemTestCase
  include PagesHelper

  setup do
    Rails.application.load_seed
    Capybara.raise_server_errors = false
  end

  teardown do
    Capybara.raise_server_errors = true
  end

  test "namespaced pages" do
    # warmup
    visit "http://" + default_panel_domain

    example_saas_pages.each{|website|
      website["pages"].each{|page|
        puts "Checking domain: " + website["domain"] + " with slug: " + website["slug"] + " for content: " + page[:body]
        url = "http://" + default_panel_domain + "/p/" + website["slug"] + "/" + page[:slug].to_s.strip
        puts "Checking slug: " + page[:slug].to_s + " with URL: " + url
        visit url
        assert_text page[:body]
        puts "OK"
      }
    }
  end


end
