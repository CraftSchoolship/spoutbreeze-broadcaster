package models

type BroadcasterRequest struct {
	BBBServerURL string `json:"bbb_server_url" binding:"required"`
	RTMPURL      string `json:"rtmp_url" binding:"required"`
	StreamURL    string `json:"stream_url" binding:"required"`
}

type BroadcasterResponse struct {
	Message string `json:"message"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}