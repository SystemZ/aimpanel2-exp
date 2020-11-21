package model

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CertDomain struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`

	Name string `bson:"name" json:"name"`

	PoolSize int `bson:"pool_size" json:"pool_size"`
}

func (cd *CertDomain) GetCollectionName() string {
	return CertDomainCollection
}

func (cd *CertDomain) GetID() primitive.ObjectID {
	return cd.ID
}

func (cd *CertDomain) SetID(id primitive.ObjectID) {
	cd.ID = id
}

func GetCertDomains() ([]CertDomain, error) {
	var domains = make([]CertDomain, 0)

	cur, err := DB.Collection(CertDomainCollection).Find(context.TODO(),
		bson.D{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.TODO())

	for cur.Next(context.TODO()) {
		var domain CertDomain
		if err := cur.Decode(&domain); err != nil {
			return nil, err
		}
		domains = append(domains, domain)
	}

	return domains, nil
}

func GetCertDomainByName(name string) (*CertDomain, error) {
	var domain CertDomain

	err := DB.Collection(CertDomainCollection).FindOne(context.TODO(), bson.D{{Key: "name", Value: name}}).Decode(&domain)
	if err != nil {
		return nil, err
	}

	return &domain, nil
}
