require "application_system_test_case"

class DomainsRoutingsTest < ApplicationSystemTestCase
  setup do
    Rails.application.load_seed
  end

  test "visiting the index" do
    visit domains_url
    assert_selector "h1", text: "Domains"
    assert_text "shop.craftexample.com"
    assert_text "craftmax.com"
    assert_text "www.craftcraft.pro"

    visit_external "http://example.com"
    assert_selector "h1", text: "Example Domain"
    #sleep 1

    visit domains_url
    visit "http://exp.lvlup.pro/domains"
    assert_selector "h1", text: "Domains"

  end
end
