class CreateWebsites < ActiveRecord::Migration[7.0]
  def change
    create_table :websites, id: :uuid do |t|
      t.belongs_to :user, null: false, foreign_key: true, type: :uuid
      t.string :name
      t.string :slug
      t.datetime :disabled_until
      t.boolean :abuse
      t.json :body

      t.timestamps
    end
  end
end
