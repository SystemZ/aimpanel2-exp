class CreateDomains < ActiveRecord::Migration[7.0]
  def change
    create_table :domains, id: :uuid do |t|
      t.string :name
      t.belongs_to :user, null: false, foreign_key: true, type: :uuid
      t.belongs_to :website, null: false, foreign_key: true, type: :uuid

      t.timestamps
    end
  end
end
