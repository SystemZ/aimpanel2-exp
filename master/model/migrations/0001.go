package migrations

import (
	"gitlab.com/systemz/aimpanel2/master/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Migration1Up() (err error) {
	id, err := primitive.ObjectIDFromHex("000000000000000000000001")
	if err != nil {
		return err
	}
	doc := &model.Config{
		ID:             id,
		MigrateVersion: 1,
	}
	err = model.Put(doc)
	return
}

func Migration1Down() {

}
