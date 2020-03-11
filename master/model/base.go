package model

import (
	"context"
	"encoding/json"
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

//Get all documents by specified query
func Get(out []interface{}, query map[string]interface{}) error {
	rows, err := CouchDB.Find(context.TODO(), query)
	if err != nil {
		return err
	}
	var result []interface{}
	for rows.Next() {
		var doc interface{}
		if err := rows.ScanDoc(&doc); err != nil {
			return err
		}
		result = append(result, doc)
	}

	out = result

	return nil
}

//Get all documents by specified selector
func GetS(out []interface{}, selector map[string]string) error {
	err := Get(out, map[string]interface{}{
		"selector": selector,
	})
	if err != nil {
		return err
	}
	return nil
}

//Get one document by specified query
func GetOne(out interface{}, query map[string]interface{}) error {
	rows, err := CouchDB.Find(context.TODO(), query)
	if err != nil {
		return err
	}

	for rows.Next() {
		if err := rows.ScanDoc(&out); err != nil {
			return err
		}
		break
	}
	return nil
}

//Get one document by specified selector
func GetOneS(out interface{}, selector map[string]string) error {
	err := GetOne(out, map[string]interface{}{
		"selector": selector,
	})
	if err != nil {
		return err
	}
	return nil
}
