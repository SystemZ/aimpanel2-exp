package model

import "go.mongodb.org/mongo-driver/bson/primitive"

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
