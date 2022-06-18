# Define your application routes per the DSL in https://guides.rubyonrails.org/routing.html
Rails.application.routes.draw do

  # alternative matching host: "exp.lvlup.pro"
  # https://guides.rubyonrails.org/routing.html#restricting-the-routes-created
  # https://blog.appsignal.com/2020/03/04/building-a-rails-app-with-multiple-subdomains.html
  # https://stackoverflow.com/questions/6054668/rails-3-routing-and-multiple-domains/6058737#6058737
  constraints domain: "lvlup.pro" do
    devise_for :users
    resources :websites
    resources :domains
  end

  # Defines the root path route ("/")
  # Ensure you have defined root_url to *something*
  # root "articles#index"
  root :to => "websites#saas"
  #get "s/:slug", to: "websites#site"
end
