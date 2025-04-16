package repositories

import (
	"spoutbreeze/initializers"
)

func StoreRTMPURL(rtmpURL string) error {
	return initializers.RedisClient.Set(initializers.RedisContext, "rtmp_url", rtmpURL, 0).Err()
}

func StoreStreamURL(streamURL string) error {
	return initializers.RedisClient.Set(initializers.RedisContext, "twitch_stream_key", streamURL, 0).Err()
}