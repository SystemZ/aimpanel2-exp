require "application_system_test_case"

class CmsDomainsTest < ApplicationSystemTestCase
  include PagesHelper

  setup do
    Rails.application.load_seed
    Capybara.raise_server_errors = false
  end

  teardown do
    Capybara.raise_server_errors = true
  end

  test "list of domains in panel" do
    visit "http://" + default_panel_domain + "/domains"
    assert_selector "h1", text: "Domains"
    example_saas_pages.each{|website|
      assert_text website["domain"]
    }
  end

  test "list of websites in panel" do
    visit "http://" + default_panel_domain + "/websites"
    assert_selector "h1", text: "Websites"
    example_saas_pages.each{|website|
      assert_text website["domain"]
    }
  end

  test "domains with websites and pages" do
    example_saas_pages.each{|website|
      website["pages"].each{|page|
        puts "Checking domain: " + website["domain"] + " with slug: " + website["slug"] + " for content: " + page[:body]
        url = "http://" + website["domain"] + "/" + page[:slug].to_s.strip
        puts "Checking slug: " + page[:slug].to_s + " with URL: " + url
        visit url
        assert_text page[:body]
        puts "OK"
      }
    }
  end

  test "domains without panel" do
    example_saas_pages.each{|website|
      website["pages"].each{|page|
        puts "Checking domain: " + website["domain"] + " with slug: " + website["slug"]

        url = "http://" + website["domain"] + "/domains"
        puts "Checking URL: " + url
        visit url
        assert_no_selector "h1", text: "Domains"
        puts "OK"

        url = "http://" + website["domain"] + "/websites"
        puts "Checking URL: " + url
        visit url
        assert_no_selector "h1", text: "Websites"
        puts "OK"
      }
    }
  end

end
