Rails.application.routes.draw do
  resources :domains
  devise_for :users
  # Define your application routes per the DSL in https://guides.rubyonrails.org/routing.html

  # Defines the root path route ("/")
  # Ensure you have defined root_url to *something*
  # root "articles#index"
end
