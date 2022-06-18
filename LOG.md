```bash
# rails bootstrap
# turn off IDE first
rails new --database=postgresql --javascript=importmap --css=tailwind cabbage
rm -rf cabbage/.git
# if you didn't close IDE, invalidate cache and restart it, remove new project dir as git repository

# database setup
cd cabbage
docker-compose up
bin/rails db:create

# auth setup
# https://betterprogramming.pub/devise-auth-setup-in-rails-7-44240aaed4be
bundle add devise
bundle install
rails g devise:install
rails g devise user
# you can edit id from int to UUID before migration
rake db:migrate
# You can copy Devise views (for customization) to your app by running: rails g devise:views
# show new routes for auth
rails routes

# domains
rails generate scaffold domain name:string user:belongs_to
# add `has_many :domains` to user model
# add type: :uuid to domains migrations
# https://stackoverflow.com/questions/40977506/rails-model-how-to-do-the-relation
rake db:migrate

# websites
rails generate scaffold website user:belongs_to name:string slug:string disabled_until:datetime abuse:boolean body:json
rake db:migrate

# domains change
rails generate scaffold domain name:string user:belongs_to website:belongs_to
# add `has_many :domains` to website model
# add `type: :uuid` to domains migrations
# add `null: true` to domains website migration
rake db:migrate

# model graphs
apt-get install -y graphviz
rake erd

# system tests
rails generate system_test domains_routing
rails test:system
rails test test/system/domains_routings_test.rb
```
