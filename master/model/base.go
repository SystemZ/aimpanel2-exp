package model

import (
	"context"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"time"
)

type Base struct {
	ID      string `json:"_id" example:"1238206236281802752"`
	Rev     string `json:"_rev,omitempty"`
	DocType string `json:"doc_type"`

	// Date of creation
	CreatedAt time.Time `json:"created_at" example:"2019-09-29T03:16:27+02:00"`

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

	rev, err := DB.Put(context.TODO(), b.ID, b.obj)
	if err != nil {
		return err
	}
	b.Rev = rev

	return nil
}

func (b *Base) Update(obj interface{}) error {
	b.obj = &obj

	rev, err := DB.Put(context.TODO(), b.ID, b.obj)
	if err != nil {
		return err
	}
	b.Rev = rev

	return nil
}

func Count(query map[string]interface{}) (int, error) {
	rows, err := DB.Find(context.TODO(), query)
	if err != nil {
		return 0, err
	}
	count := 0
	// release resources
	defer rows.Close()
	for rows.Next() {
		count++
	}

	return count, nil
}

//Get all documents by specified query
func Get(out interface{}, query map[string]interface{}) error {
	rows, err := DB.Find(context.TODO(), query)
	if err != nil {
		return err
	}
	var result []interface{}
	// release resources
	defer rows.Close()
	for rows.Next() {
		var doc interface{}
		if err := rows.ScanDoc(&doc); err != nil {
			return err
		}
		result = append(result, doc)
	}

	jsonStr, _ := json.Marshal(result)
	err = json.Unmarshal(jsonStr, &out)
	if err != nil {
		return err
	}

	return nil
}

//Get all documents by specified selector
func GetS(out interface{}, selector map[string]interface{}) error {
	err := Get(out, map[string]interface{}{
		"selector": selector,
	})
	if err != nil {
		return err
	}
	return nil
}

//Get all documents by specified selector and limit
func GetSLimit(out interface{}, limit int, selector map[string]string) error {
	err := Get(out, map[string]interface{}{
		"selector": selector,
		"limit":    limit,
	})
	if err != nil {
		return err
	}
	return nil
}

//Get one document by specified query
func GetOne(out interface{}, query map[string]interface{}) error {
	rows, err := DB.Find(context.TODO(), query)
	if err != nil {
		return err
	}
	// release resources
	defer rows.Close()
	for rows.Next() {
		if err := rows.ScanDoc(&out); err != nil {
			return err
		}
		break
	}
	return nil
}

//Get one document by specified selector
func GetOneS(out interface{}, selector map[string]interface{}) error {
	err := GetOne(out, map[string]interface{}{
		"selector": selector,
	})
	if err != nil {
		return err
	}
	return nil
}

func Delete(id string, rev string) error {
	_, err := DB.Delete(context.TODO(), id, rev)
	if err != nil {
		logrus.Info(err)
		return err
	}
	return nil
}
