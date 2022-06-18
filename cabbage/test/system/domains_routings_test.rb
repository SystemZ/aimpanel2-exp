require "application_system_test_case"

class DomainsRoutingsTest < ApplicationSystemTestCase
  saas_domains = %w[shop.craftexample.com craftmax.com www.craftcraft.pro]

  setup do
    Rails.application.load_seed
    Capybara.raise_server_errors = false
  end

  teardown do
    Capybara.raise_server_errors = true
  end

  test "list of domains on panel domain" do
    #visit domains_url
    visit "http://exp.lvlup.pro/domains"
    assert_selector "h1", text: "Domains"
    assert_text "shop.craftexample.com"
    assert_text "craftmax.com"
    assert_text "www.craftcraft.pro"
  end

  test "list of websites on panel domain" do
    visit "http://exp.lvlup.pro/websites"
    assert_text "shop.craftexample.com"
    assert_text "craftmax.com"
    assert_text "www.craftcraft.pro"
  end

  test "saas domains body" do
    saas_domains.each{|domain|
      puts domain
      visit "http://" + domain
      assert_text "Welcome to " + domain
    }
  end

  test "saas domains lack of main panel" do
    saas_domains.each{|domain|
      puts domain
      visit "http://" + domain + "/domains"
      assert_no_selector "h1", text: "Domains"
      visit "http://" + domain + "/websites"
      assert_no_selector "h1", text: "Websites"
    }
  end

  #visit_external "http://example.com"
end
