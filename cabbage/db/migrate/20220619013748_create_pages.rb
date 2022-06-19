class CreatePages < ActiveRecord::Migration[7.0]
  def change
    create_table :pages, id: :uuid do |t|
      t.belongs_to :website, null: false, foreign_key: true, type: :uuid
      t.belongs_to :user, null: false, foreign_key: true, type: :uuid
      t.integer :breed, default: 1
      t.string :language, limit: 2
      t.datetime :active_from
      t.string :slug
      t.string :title
      t.json :body

      t.timestamps
    end
  end
end
