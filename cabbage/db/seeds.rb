# This file should contain all the record creation needed to seed the database with its default values.
# The data can then be loaded with the bin/rails db:seed command (or created alongside the database with db:setup).
#
# Examples:
#
#   movies = Movie.create([{ name: "Star Wars" }, { name: "Lord of the Rings" }])
#   Character.create(name: "Luke", movie: movies.first)

user1 = User.create!(:email => 'user1@example.com', :password => 'password')

website1 = Website.create!(user: user1, name: "shop.craftexample.com", slug: "shop-craftexample-com")
Page.create!(website: website1, user: user1, language: "en", slug: nil, title: "Homepage", body: "Welcome to shop.craftexample.com")
Page.create!(website: website1, user: user1, language: "en", slug: "rulez", title: "Rules of server", body: "some rules here")
Page.create!(website: website1, user: user1, language: "en", slug: "aboutz", title: "About team", body: "some team info will be here in the future")
Domain.create!(name: "shop.craftexample.com", user: user1, website: website1)

website2 = Website.create!(user: user1, name: "craftmax.com", slug: "craftmax-com")
Page.create!(website: website2, user: user1, language: "en", slug: nil, title: "Homepage", body: "Welcome to craftmax.com")
Page.create!(website: website2, user: user1, language: "en", slug: "tos", title: "ToS", body: "Long ToS")
Domain.create!(name: "craftmax.com", user: user1, website: website2)

user2 = User.create!(:email => 'user2@example.com', :password => 'password')

website3 = Website.create!(user: user2, name: "www.craftcraft.pro", slug: "www-craftcraft-pro")
Page.create!(website: website3, user: user2, language: "en", slug: nil, title: "Homepage", body: "Welcome to www.craftcraft.pro")
Domain.create!(name: "www.craftcraft.pro", user: user2, website: website3)
