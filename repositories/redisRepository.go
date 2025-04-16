package repositories

import (
	"spoutbreeze/initializers"
)

func StoreRTMPURL(rtmpURL string) error {
	return initializers.RedisClient.Set(initializers.RedisContext, "rtmp_url", rtmpURL, 0).Err()
}

func StoreStreamKey(streamKEY string) error {
	return initializers.RedisClient.Set(initializers.RedisContext, "twitch_stream_key", streamKEY, 0).Err()
}