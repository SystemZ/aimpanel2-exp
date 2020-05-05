package model

import (
	"github.com/go-redis/redis"
	"gitlab.com/systemz/aimpanel2/lib/filemanager"
)

func GetSlaveCommit(redis *redis.Client) (string, error) {
	return redis.Get("slave_commit").Result()
}

func GetSlaveUrl(redis *redis.Client) (string, error) {
	return redis.Get("slave_url").Result()
}

func GsFilesSubscribe(redis *redis.Client, gsId string) (*redis.PubSub, error) {
	pubsub := redis.Subscribe("gs_files_" + gsId)

	_, err := pubsub.Receive()
	if err != nil {
		return nil, err
	}

	return pubsub, nil
}

func GsFilesPublish(redis *redis.Client, gsId string, files *filemanager.Node) error {
	err := redis.Publish("gs_files_"+gsId, files.String()).Err()
	if err != nil {
		return err
	}

	return nil
}
