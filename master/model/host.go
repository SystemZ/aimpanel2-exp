package model

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Host struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty" example:"1238206236281802752"`

	// User assigned name
	Name string `bson:"name" json:"name" example:"My Great Linux server"`

	// User ID
	UserId primitive.ObjectID `bson:"user_id" json:"user_id" example:"100112233-4455-6677-8899-aabbccddeeff"`

	// User assigned ip
	Ip string `bson:"ip" json:"ip" example:"192.51.100.128"`

	// Generated token for host
	Token string `bson:"token" json:"token" example:"VRAKKUBHNIMKLXSXLWTQAOFGOMSSCXOO"`

	// Metric frequency
	MetricFrequency int `bson:"metric_frequency" json:"metric_frequency" example:"30"`

	// Host OS
	OS string `bson:"os" json:"os" example:"linux"`

	// Host platform
	Platform string `bson:"platform" json:"platform" example:"ubuntu"`

	// Host platform family
	PlatformFamily string `bson:"platform_family" json:"platform_family" example:"debian"`

	// Host platform version
	PlatformVersion string `bson:"platform_version" json:"platform_version" example:"18.04"`

	// Host kernel version
	KernelVersion string `bson:"kernel_version" json:"kernel_version" example:"5.3.0-1-generic"`

	// Host arch
	KernelArch string `bson:"kernel_arch" json:"kernel_arch" example:"x86_64"`

	// State
	// 0 off, 1 running
	State uint `bson:"state" json:"state" example:"1"`
}

func (h *Host) GetCollectionName() string {
	return hostCollection
}

func (h *Host) GetID() primitive.ObjectID {
	return h.ID
}

func (h *Host) SetID(id primitive.ObjectID) {
	h.ID = id
}

func GetHosts() ([]Host, error) {
	var hosts []Host

	cur, err := DB.Collection(hostCollection).Find(context.TODO(),
		bson.D{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.TODO())

	for cur.Next(context.TODO()) {
		var host Host
		if err := cur.Decode(&host); err != nil {
			return nil, err
		}
		hosts = append(hosts, host)
	}

	return hosts, nil
}

func GetHostById(id primitive.ObjectID) (*Host, error) {
	var host Host

	err := DB.Collection(hostCollection).FindOne(context.TODO(), bson.D{{Key: "_id", Value: id}}).Decode(&host)
	if err != nil {
		return nil, err
	}

	return &host, nil
}

func GetHostTokenById(id primitive.ObjectID) (string, error) {
	var host Host

	err := DB.Collection(hostCollection).FindOne(context.TODO(), bson.D{{Key: "_id", Value: id}}).Decode(&host)
	if err != nil {
		return "", err
	}

	return host.Token, nil
}

func GetHostByToken(token string) (*Host, error) {
	var host Host

	err := DB.Collection(hostCollection).FindOne(context.TODO(), bson.D{{Key: "token", Value: token}}).Decode(&host)
	if err != nil {
		return nil, err
	}

	return &host, nil
}

func GetHostsByUserId(userId primitive.ObjectID) ([]Host, error) {
	var hosts []Host

	cur, err := DB.Collection(hostCollection).Find(context.TODO(),
		bson.D{{Key: "user_id", Value: userId}})
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.TODO())

	for cur.Next(context.TODO()) {
		var host Host
		if err := cur.Decode(&host); err != nil {
			return nil, err
		}
		hosts = append(hosts, host)
	}

	return hosts, nil
}

func GetHostMetricsByHostId(hostId primitive.ObjectID, limit int64) ([]MetricHost, error) {
	var metrics []MetricHost

	opts := options.Find()
	opts.SetLimit(limit)
	opts.SetSort(bson.D{{Key: "_id", Value: -1}})

	cur, err := DB.Collection(metricHostCollection).Find(context.TODO(),
		bson.D{{Key: "host_id", Value: hostId}}, opts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.TODO())

	for cur.Next(context.TODO()) {
		var metric MetricHost
		if err := cur.Decode(&metric); err != nil {
			return nil, err
		}
		metrics = append(metrics, metric)
	}

	return metrics, nil
}
