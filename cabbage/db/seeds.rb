# This file should contain all the record creation needed to seed the database with its default values.
# The data can then be loaded with the bin/rails db:seed command (or created alongside the database with db:setup).
#
# Examples:
#
#   movies = Movie.create([{ name: "Star Wars" }, { name: "Lord of the Rings" }])
#   Character.create(name: "Luke", movie: movies.first)

user1 = User.create!(:email => 'user1@example.com', :password => 'password')
user2 = User.create!(:email => 'user2@example.com', :password => 'password')
Domain.create!(name: "shop.craftexample.com", user: user1)
Domain.create!(name: "craftmax.com", user: user1)
Domain.create!(name: "www.craftcraft.pro", user: user2)
