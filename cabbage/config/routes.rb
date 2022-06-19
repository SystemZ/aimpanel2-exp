# Define your application routes per the DSL in https://guides.rubyonrails.org/routing.html
Rails.application.routes.draw do
  # alternative matching host: "exp.lvlup.pro"
  # https://guides.rubyonrails.org/routing.html#restricting-the-routes-created
  # https://guides.rubyonrails.org/routing.html#advanced-constraints
  # https://blog.appsignal.com/2020/03/04/building-a-rails-app-with-multiple-subdomains.html
  # https://stackoverflow.com/questions/6054668/rails-3-routing-and-multiple-domains/6058737#6058737
  constraints domain: "lvlup.pro" do
    devise_for :users
    resources :websites
    resources :domains
    resources :pages
  end

  root :to => "pages#page_by_slug_and_domain"
  match '*path' => 'pages#page_by_slug_and_domain', via: :get
end
