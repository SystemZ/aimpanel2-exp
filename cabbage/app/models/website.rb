class Website < ApplicationRecord
  belongs_to :user
  has_many :domains
  has_many :pages
end
