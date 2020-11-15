package model

import "go.mongodb.org/mongo-driver/bson/primitive"

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
