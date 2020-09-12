package migrations

import (
	"context"
	"gitlab.com/systemz/aimpanel2/master/model"
)

func Migration3Up() (err error) {
	err = model.DB.Collection(model.EventCollection).Drop(context.TODO())
	if err != nil {
		return
	}

	return
}

func Migration3Down() {

}
