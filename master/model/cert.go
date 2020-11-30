package model

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Cert struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`

	Cert string `bson:"cert" json:"cert"`

	PrivateKey string `bson:"private_key" json:"private_key"`

	DomainId primitive.ObjectID `bson:"domain_id" json:"domain_id"`
}

func (c *Cert) GetCollectionName() string {
	return CertCollection
}

func (c *Cert) GetID() primitive.ObjectID {
	return c.ID
}

func (c *Cert) SetID(id primitive.ObjectID) {
	c.ID = id
}

func GetCertByDomainId(domainId primitive.ObjectID) (*Cert, error) {
	var cert Cert

	err := DB.Collection(CertCollection).FindOne(context.TODO(), bson.D{{Key: "domain_id", Value: domainId}}).Decode(&cert)
	if err != nil {
		return nil, err
	}

	return &cert, nil
}
