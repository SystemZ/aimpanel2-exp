module PagesHelper
  def default_panel_domain
    "exp.lvlup.pro"
  end
  def example_users
    [
      {"email" => "user1@example.com", "password" => "password"},
      {"email" => "user2@example.com", "password" => "password"},
      {"email" => "user3@example.com", "password" => "password"},
    ]
  end
  def example_saas_pages
    [
      {
        "user" => 0,
        "domain" => "shop.craftexample.com",
        "slug" => "shop-craftexample-com",
        "pages" => [
          {
            "language": "en",
            "slug": nil,
            "title": "Homepage",
            "body": "Welcome to shop.craftexample.com"
          },
          {
            "language": "en",
            "slug": "rulez",
            "title": "Rules of server",
            "body": "some rules here"
          },
          {
            "language": "en",
            "slug": "aboutz",
            "title": "About team",
            "body": "some team info will be here in the future"
          },
        ]
      },
      {
        "user" => 0,
        "domain" => "craftmax.com",
        "slug" => "craftmax-com",
        "pages" => [
          {
            "language": "en",
            "slug": nil,
            "title": "Homepage",
            "body": "Welcome to craftmax.com"
          },
          {
            "language": "en",
            "slug": "tos",
            "title": "ToS",
            "body": "Long ToS"
          },
        ]
      },
      {
        "user" => 1,
        "domain" => "www.craftcraft.pro",
        "slug" => "www-craftcraft-pro",
        "pages" => [
          {
            "language": "en",
            "slug": nil,
            "title": "Homepage",
            "body": "Welcome to www.craftcraft.pro"
          },
        ]
      },
    ]
  end
end
