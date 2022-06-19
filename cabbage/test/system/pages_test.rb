require "application_system_test_case"

class PagesTest < ApplicationSystemTestCase
  setup do
    @page = pages(:one)
  end

  test "visiting the index" do
    visit pages_url
    assert_selector "h1", text: "Pages"
  end

  test "should create page" do
    visit pages_url
    click_on "New page"

    fill_in "Active from", with: @page.active_from
    fill_in "Body", with: @page.body
    fill_in "Breed", with: @page.breed
    fill_in "Language", with: @page.language
    fill_in "Slug", with: @page.slug
    fill_in "Title", with: @page.title
    fill_in "User", with: @page.user_id
    fill_in "Website", with: @page.website_id
    click_on "Create Page"

    assert_text "Page was successfully created"
    click_on "Back"
  end

  test "should update Page" do
    visit page_url(@page)
    click_on "Edit this page", match: :first

    fill_in "Active from", with: @page.active_from
    fill_in "Body", with: @page.body
    fill_in "Breed", with: @page.breed
    fill_in "Language", with: @page.language
    fill_in "Slug", with: @page.slug
    fill_in "Title", with: @page.title
    fill_in "User", with: @page.user_id
    fill_in "Website", with: @page.website_id
    click_on "Update Page"

    assert_text "Page was successfully updated"
    click_on "Back"
  end

  test "should destroy Page" do
    visit page_url(@page)
    click_on "Destroy this page", match: :first

    assert_text "Page was successfully destroyed"
  end
end
