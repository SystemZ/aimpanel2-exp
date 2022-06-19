# This file should contain all the record creation needed to seed the database with its default values.
# The data can then be loaded with the bin/rails db:seed command (or created alongside the database with db:setup).
#
# Examples:
#
#   movies = Movie.create([{ name: "Star Wars" }, { name: "Lord of the Rings" }])
#   Character.create(name: "Luke", movie: movies.first)
include PagesHelper

websites = example_saas_pages
user = []
example_users.each { |u|
  user << User.create!(:email => u["email"], :password => u["password"])
}

websites.each { |website|
  new_website = Website.create!(user: user[website["user"]], name: website["domain"], slug: website["slug"])
  website["pages"].each { |page|
    Page.create!(website: new_website, user: user[website["user"]], language: page[:language], slug: page[:slug], title: page[:title], body: page[:body])
  }
  Domain.create!(name: website["domain"], user: user[website["user"]], website: new_website)
}
