package migrations

import (
	"context"
	"gitlab.com/systemz/aimpanel2/master/model"
)

func Migration2Up() (err error) {
	err = model.DB.Collection(model.PermissionCollection).Drop(context.TODO())
	if err != nil {
		return
	}

	err = model.DB.Collection(model.GameFileCollection).Drop(context.TODO())
	if err != nil {
		return
	}

	return
}

func Migration2Down() {

}
