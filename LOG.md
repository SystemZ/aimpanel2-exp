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
```
