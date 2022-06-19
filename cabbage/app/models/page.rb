class Page < ApplicationRecord
  belongs_to :website
  belongs_to :user
  enum :breed, { none: 0, markdown: 1, gallery: 2 }, suffix: true, default: :markdown
end
