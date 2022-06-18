class Website < ApplicationRecord
  belongs_to :user
  has_many :domains
end
