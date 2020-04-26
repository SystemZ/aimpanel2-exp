package model

import (
	"github.com/go-redis/redis"
	"time"
)

func SetGsRestart(redis *redis.Client, gsId string, state int) {
	redis.Set("gs_restart_id"+gsId, state, 24*time.Hour)
}

func GetGsRestart(redis *redis.Client, gsId string) (int64, error) {
	return redis.Get("gs_restart_id_" + gsId).Int64()
}

func DelGsRestart(redis *redis.Client, gsId string) {
	redis.Del("gs_restart_id_" + gsId)
}

func GetSlaveCommit(redis *redis.Client) (string, error) {
	return redis.Get("slave_commit").Result()
}

func GetSlaveUrl(redis *redis.Client) (string, error) {
	return redis.Get("slave_url").Result()
}

func SetGsStart(redis *redis.Client, gsId string, state int) {
	redis.Set("gs_start_id_"+gsId, state, 24*time.Hour)
}

func GetGsStart(redis *redis.Client, gsId string) (int64, error) {
	return redis.Get("gs_start_id_" + gsId).Int64()
}

func DelGsStart(redis *redis.Client, gsId string) {
	redis.Del("gs_start_id_" + gsId)
}

func SetAgentHeartbeat(redis *redis.Client, token string, timestamp int64) {
	redis.Set("agent_heartbeat_token_"+token, timestamp, 24*time.Hour)
}

func SetWrapperHeartbeat(redis *redis.Client, gsId string, timestamp int64) {
	redis.Set("wrapper_heartbeat_id_"+gsId, timestamp, 24*time.Hour)
}
