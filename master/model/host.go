package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Host struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty" example:"1238206236281802752"`

	// User assigned name
	Name string `json:"name" example:"My Great Linux server"`

	// User ID
	UserId primitive.ObjectID `json:"user_id" example:"100112233-4455-6677-8899-aabbccddeeff"`

	// User assigned ip
	Ip string `json:"ip" example:"192.51.100.128"`

	// Generated token for host
	Token string `json:"token" example:"VRAKKUBHNIMKLXSXLWTQAOFGOMSSCXOO"`

	// Metric frequency
	MetricFrequency int `json:"metric_frequency" example:"30"`

	// Host OS
	OS string `json:"os" example:"linux"`

	// Host platform
	Platform string `json:"platform" example:"ubuntu"`

	// Host platform family
	PlatformFamily string `json:"platform_family" example:"debian"`

	// Host platform version
	PlatformVersion string `json:"platform_version" example:"18.04"`

	// Host kernel version
	KernelVersion string `json:"kernel_version" example:"5.3.0-1-generic"`

	// Host arch
	KernelArch string `json:"kernel_arch" example:"x86_64"`

	// State
	// 0 off, 1 running
	State uint `json:"state" example:"1"`
}

func (h *Host) GetCollectionName() string {
	return "hosts"
}

func GetHosts() []Host {
	var hosts []Host

	err := GetS(&hosts, map[string]interface{}{
		"doc_type": "host",
	})
	if err != nil {
		return nil
	}

	return hosts
}

func GetHost(hostId primitive.ObjectID) *Host {
	var host Host
	err := GetOneS(&host, map[string]interface{}{
		"doc_type": "host",
		"_id":      hostId,
	})
	if err != nil {
		return nil
	}

	return &host
}

func GetHostToken(hostId primitive.ObjectID) string {
	var host Host
	err := GetOneS(&host, map[string]interface{}{
		"doc_type": "host",
		"_id":      hostId,
	})
	if err != nil {
		return ""
	}

	return host.Token
}

func GetHostByToken(token string) *Host {
	var host Host
	err := GetOneS(&host, map[string]interface{}{
		"doc_type": "host",
		"token":    token,
	})
	if err != nil {
		return nil
	}

	return &host
}

func GetHostsByUserId(userId primitive.ObjectID) []Host {
	var hosts []Host

	err := GetS(&hosts, map[string]interface{}{
		"doc_type": "host",
		"user_id":  userId,
	})
	if err != nil {
		return nil
	}

	return hosts
}

func GetHostMetrics(hostId primitive.ObjectID, limit int) []MetricHost {
	var metrics []MetricHost

	err := Get(&metrics, map[string]interface{}{
		"selector": map[string]string{
			"doc_type": "metric_host",
			"host_id":  hostId.String(),
		},
		"limit": limit,
		"sort": []map[string]interface{}{
			{"created_at": "desc"},
		},
	})
	if err != nil {
		return nil
	}

	return metrics
}
