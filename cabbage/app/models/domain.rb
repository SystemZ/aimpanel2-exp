class Domain < ApplicationRecord
  belongs_to :user
  belongs_to :website, optional: true
end
