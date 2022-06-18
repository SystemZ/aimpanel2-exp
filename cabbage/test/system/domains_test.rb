require "application_system_test_case"

class DomainsTest < ApplicationSystemTestCase
  setup do
    @domain = domains(:one)
  end

  test "visiting the index" do
    visit domains_url
    assert_selector "h1", text: "Domains"
  end

  test "should create domain" do
    visit domains_url
    click_on "New domain"

    fill_in "Name", with: @domain.name
    fill_in "User", with: @domain.user_id
    fill_in "Website", with: @domain.website_id
    click_on "Create Domain"

    assert_text "Domain was successfully created"
    click_on "Back"
  end

  test "should update Domain" do
    visit domain_url(@domain)
    click_on "Edit this domain", match: :first

    fill_in "Name", with: @domain.name
    fill_in "User", with: @domain.user_id
    fill_in "Website", with: @domain.website_id
    click_on "Update Domain"

    assert_text "Domain was successfully updated"
    click_on "Back"
  end

  test "should destroy Domain" do
    visit domain_url(@domain)
    click_on "Destroy this domain", match: :first

    assert_text "Domain was successfully destroyed"
  end
end
