package model

import (
	"context"
	"encoding/json"
	"github.com/go-kivik/kivik/v3"
	"time"
)

type Base struct {
	ID        string    `json:"_id"`
	Rev       string    `json:"_rev,omitempty"`
	CreatedAt time.Time `json:"created_at"`

	obj interface{}
}

func (b *Base) SetObj(obj interface{}) {
	b.obj = obj
}

func (b *Base) GetJSON() ([]byte, error) {
	return json.Marshal(b.obj)
}

func (b *Base) Put(obj interface{}) error {
	b.obj = &obj

	b.ID = Snowflake.Generate().String()
	b.CreatedAt = time.Now()

	rev, err := CouchDB.Put(context.TODO(), b.ID, b.obj)
	if err != nil {
		return err
	}
	b.Rev = rev

	return nil
}

func Count(query map[string]interface{}) (int, error) {
	rows, err := CouchDB.Find(context.TODO(), query)
	if err != nil {
		return 0, err
	}
	count := 0
	for rows.Next() {
		count++
	}

	return count, nil
}

func Get(query map[string]interface{}) (*kivik.Rows, error) {
	rows, err := CouchDB.Find(context.TODO(), query)
	if err != nil {
		return nil, err
	}
	return rows, nil
}
