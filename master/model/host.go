package model

import (
	"context"
	"gitlab.com/systemz/aimpanel2/master/events"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strings"
	"time"
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

	MetricMaxS int `bson:"metric_max_s" json:"metric_max_s" example:"30"`

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

	// Last successful heartbeat received
	LastHeartbeat time.Time `bson:"last_heartbeat" json:"last_heartbeat" example:"2019-09-29T03:16:27+02:00"`

	HwId string `bson:"hw_id" json:"hw_id" example:"0ea13526-5349-4671-9b6e-626d5678cc5f"`

	//Domain for internal use (certs, slave https server)
	Domain string `bson:"domain" json:"-"`
}

func (h *Host) GetCollectionName() string {
	return HostCollection
}

func (h *Host) GetID() primitive.ObjectID {
	return h.ID
}

func (h *Host) SetID(id primitive.ObjectID) {
	h.ID = id
}

func GetHosts() ([]Host, error) {
	var hosts = make([]Host, 0)

	cur, err := DB.Collection(HostCollection).Find(context.TODO(),
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

	err := DB.Collection(HostCollection).FindOne(context.TODO(), bson.D{{Key: "_id", Value: id}}).Decode(&host)
	if err != nil {
		return nil, err
	}

	return &host, nil
}

func GetHostTokenById(id primitive.ObjectID) (string, error) {
	var host Host

	err := DB.Collection(HostCollection).FindOne(context.TODO(), bson.D{{Key: "_id", Value: id}}).Decode(&host)
	if err != nil {
		return "", err
	}

	return host.Token, nil
}

func GetHostByToken(token string) (*Host, error) {
	var host Host

	err := DB.Collection(HostCollection).FindOne(context.TODO(), bson.D{{Key: "token", Value: token}}).Decode(&host)
	if err != nil {
		return nil, err
	}

	return &host, nil
}

func GetHostsByUser(user User) ([]Host, error) {
	var hosts = make([]Host, 0)

	// show hosts that belong to user
	searchTerms := bson.D{{Key: "user_id", Value: user.ID}}
	// admin have access to all hosts
	if user.Admin {
		searchTerms = bson.D{}
	}

	cur, err := DB.Collection(HostCollection).Find(context.TODO(), searchTerms)
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.TODO())

	// get currently connected hosts (and browsers)
	sseChannels := events.SSE.Channels()

	// loop through all hosts from DB for this user
	for cur.Next(context.TODO()) {
		var host Host
		if err := cur.Decode(&host); err != nil {
			return nil, err
		}
		// by default mark host as disconnected
		host.State = 0
		for _, chRaw := range sseChannels {
			ch := strings.Replace(chRaw, "/v1/events/", "", 1)
			if ch == host.Token {
				// SSE channel with host's token means that host is currently connected
				host.State = 1
			}
		}
		hosts = append(hosts, host)
	}

	return hosts, nil
}

func GetHostMetricsByHostId(hostId primitive.ObjectID, limit int64) ([]MetricHost, error) {
	var metrics = make([]MetricHost, 0)

	opts := options.Find()
	opts.SetLimit(limit)
	opts.SetSort(bson.D{{Key: "_id", Value: -1}})

	cur, err := DB.Collection(MetricHostCollection).Find(context.TODO(),
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
