package model

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Document interface {
	GetCollectionName() string
}

type Base struct{}

func Put(d Document) error {
	_, err := DB.Collection(d.GetCollectionName()).InsertOne(context.TODO(), d)
	if err != nil {
		return err
	}

	return nil
}

func Update(d Document) error {
	return nil
}

func (b *Base) Put(obj Document) error {
	//b.obj = &obj
	//b.CreatedAt = time.Now()
	//
	//res, err := DB.Collection(obj.GetCollectionName()).InsertOne(context.TODO(), b.obj)
	//if err != nil {
	//	return err
	//}
	//
	//logrus.Info(res.InsertedID)

	return nil
}

func (b *Base) Update(obj interface{}) error {
	//b.obj = &obj
	//
	//rev, err := DB.Put(context.TODO(), b.ID, b.obj)
	//if err != nil {
	//	return err
	//}
	//b.Rev = rev

	return nil
}

func Count(query map[string]interface{}) (int, error) {
	//rows, err := DB.Find(context.TODO(), query)
	//if err != nil {
	//	return 0, err
	//}
	//count := 0
	//// release resources
	//defer rows.Close()
	//for rows.Next() {
	//	count++
	//}

	return 0, nil
}

//Get all documents by specified query
func Get(out interface{}, query map[string]interface{}) error {
	//rows, err := DB.Find(context.TODO(), query)
	//if err != nil {
	//	return err
	//}
	//var result []interface{}
	//// release resources
	//defer rows.Close()
	//for rows.Next() {
	//	var doc interface{}
	//	if err := rows.ScanDoc(&doc); err != nil {
	//		return err
	//	}
	//	result = append(result, doc)
	//}
	//
	//jsonStr, _ := json.Marshal(result)
	//err = json.Unmarshal(jsonStr, &out)
	//if err != nil {
	//	return err
	//}

	return nil
}

//Get all documents by specified selector
func GetS(out interface{}, selector map[string]interface{}) error {
	//err := Get(out, map[string]interface{}{
	//	"selector": selector,
	//})
	//if err != nil {
	//	return err
	//}
	return nil
}

//Get all documents by specified selector and limit
func GetSLimit(out interface{}, limit int, selector map[string]string) error {
	//err := Get(out, map[string]interface{}{
	//	"selector": selector,
	//	"limit":    limit,
	//})
	//if err != nil {
	//	return err
	//}
	return nil
}

//Get one document by specified query
func GetOne(out interface{}, query map[string]interface{}) error {
	//rows, err := DB.Find(context.TODO(), query)
	//if err != nil {
	//	return err
	//}
	//// release resources
	//defer rows.Close()
	//for rows.Next() {
	//	if err := rows.ScanDoc(&out); err != nil {
	//		return err
	//	}
	//	break
	//}
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

func Delete(id primitive.ObjectID, rev string) error {
	//_, err := DB.Delete(context.TODO(), id, rev)
	//if err != nil {
	//	logrus.Info(err)
	//	return err
	//}
	return nil
}
