json.extract! website, :id, :user_id, :name, :slug, :disabled_until, :abuse, :body, :created_at, :updated_at
json.url website_url(website, format: :json)
