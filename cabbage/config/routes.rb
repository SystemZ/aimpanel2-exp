Rails.application.routes.draw do
  resources :websites
  resources :domains
  devise_for :users
  # Define your application routes per the DSL in https://guides.rubyonrails.org/routing.html

  # Defines the root path route ("/")
  # Ensure you have defined root_url to *something*
  # root "articles#index"
  get "s/:slug", to: "websites#site"

  # TODO routing of big number of dynamic domains
  # https://stackoverflow.com/questions/6054668/rails-3-routing-and-multiple-domains/6058737#6058737

  # main domain exp.lvlup.pro
  #root :to => "static#home", :constraints => { :domain => "exp.lvlup.pro" }

  # customer domains examples
  # shop.craftexample.com
  # craftmax.com
  # www.craftcraft.pro
  #root :to => "portfolios#show"
end
