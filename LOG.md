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
rake db:migrate
# You can copy Devise views (for customization) to your app by running: rails g devise:views
```
