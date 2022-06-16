require "application_system_test_case"

class WebsitesTest < ApplicationSystemTestCase
  setup do
    @website = websites(:one)
  end

  test "visiting the index" do
    visit websites_url
    assert_selector "h1", text: "Websites"
  end

  test "should create website" do
    visit websites_url
    click_on "New website"

    check "Abuse" if @website.abuse
    fill_in "Body", with: @website.body
    fill_in "Disabled until", with: @website.disabled_until
    fill_in "Name", with: @website.name
    fill_in "Slug", with: @website.slug
    fill_in "User", with: @website.user_id
    click_on "Create Website"

    assert_text "Website was successfully created"
    click_on "Back"
  end

  test "should update Website" do
    visit website_url(@website)
    click_on "Edit this website", match: :first

    check "Abuse" if @website.abuse
    fill_in "Body", with: @website.body
    fill_in "Disabled until", with: @website.disabled_until
    fill_in "Name", with: @website.name
    fill_in "Slug", with: @website.slug
    fill_in "User", with: @website.user_id
    click_on "Update Website"

    assert_text "Website was successfully updated"
    click_on "Back"
  end

  test "should destroy Website" do
    visit website_url(@website)
    click_on "Destroy this website", match: :first

    assert_text "Website was successfully destroyed"
  end
end
