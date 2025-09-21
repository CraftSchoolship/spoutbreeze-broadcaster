package models_test

import (
	"encoding/json"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"spoutbreeze/models"
)

var _ = Describe("Models", func() {
	Describe("BroadcasterRequest", func() {
		Context("when creating a new broadcaster request", func() {
			It("should have all required fields", func() {
				request := models.BroadcasterRequest{
					BBBServerURL:      "https://example.com/bigbluebutton",
					BBBHealthCheckURL: "https://example.com/health",
					RTMPURL:           "rtmp://streaming.example.com/live",
					StreamKey:         "stream-123",
				}

				Expect(request.BBBServerURL).To(Equal("https://example.com/bigbluebutton"))
				Expect(request.BBBHealthCheckURL).To(Equal("https://example.com/health"))
				Expect(request.RTMPURL).To(Equal("rtmp://streaming.example.com/live"))
				Expect(request.StreamKey).To(Equal("stream-123"))
			})

			It("should be serializable to JSON", func() {
				request := models.BroadcasterRequest{
					BBBServerURL:      "https://example.com/bigbluebutton",
					BBBHealthCheckURL: "https://example.com/health",
					RTMPURL:           "rtmp://streaming.example.com/live",
					StreamKey:         "stream-123",
				}

				jsonData, err := json.Marshal(request)
				Expect(err).NotTo(HaveOccurred())
				Expect(string(jsonData)).To(ContainSubstring("bbb_server_url"))
				Expect(string(jsonData)).To(ContainSubstring("bbb_health_check_url"))
				Expect(string(jsonData)).To(ContainSubstring("rtmp_url"))
				Expect(string(jsonData)).To(ContainSubstring("stream_key"))
			})

			It("should be deserializable from JSON", func() {
				jsonData := `{
					"bbb_server_url": "https://example.com/bigbluebutton",
					"bbb_health_check_url": "https://example.com/health",
					"rtmp_url": "rtmp://streaming.example.com/live",
					"stream_key": "stream-123"
				}`

				var request models.BroadcasterRequest
				err := json.Unmarshal([]byte(jsonData), &request)
				Expect(err).NotTo(HaveOccurred())
				Expect(request.BBBServerURL).To(Equal("https://example.com/bigbluebutton"))
				Expect(request.BBBHealthCheckURL).To(Equal("https://example.com/health"))
				Expect(request.RTMPURL).To(Equal("rtmp://streaming.example.com/live"))
				Expect(request.StreamKey).To(Equal("stream-123"))
			})
		})

		Context("when creating with empty values", func() {
			It("should handle empty strings", func() {
				request := models.BroadcasterRequest{
					BBBServerURL:      "",
					BBBHealthCheckURL: "",
					RTMPURL:           "",
					StreamKey:         "",
				}

				Expect(request.BBBServerURL).To(Equal(""))
				Expect(request.BBBHealthCheckURL).To(Equal(""))
				Expect(request.RTMPURL).To(Equal(""))
				Expect(request.StreamKey).To(Equal(""))
			})
		})
	})

	Describe("BroadcasterResponse", func() {
		Context("when creating a broadcaster response", func() {
			It("should have message field", func() {
				response := models.BroadcasterResponse{
					Message: "Broadcasting session started successfully",
				}

				Expect(response.Message).To(Equal("Broadcasting session started successfully"))
			})

			It("should be serializable to JSON", func() {
				response := models.BroadcasterResponse{
					Message: "Broadcasting session started successfully",
				}

				jsonData, err := json.Marshal(response)
				Expect(err).NotTo(HaveOccurred())
				Expect(string(jsonData)).To(ContainSubstring("message"))
				Expect(string(jsonData)).To(ContainSubstring("Broadcasting session started successfully"))
			})

			It("should be deserializable from JSON", func() {
				jsonData := `{"message": "Broadcasting session started successfully"}`

				var response models.BroadcasterResponse
				err := json.Unmarshal([]byte(jsonData), &response)
				Expect(err).NotTo(HaveOccurred())
				Expect(response.Message).To(Equal("Broadcasting session started successfully"))
			})
		})
	})

	Describe("ErrorResponse", func() {
		Context("when creating an error response", func() {
			It("should have message field", func() {
				errorResponse := models.ErrorResponse{
					Message: "An error occurred",
				}

				Expect(errorResponse.Message).To(Equal("An error occurred"))
			})

			It("should be serializable to JSON", func() {
				errorResponse := models.ErrorResponse{
					Message: "An error occurred",
				}

				jsonData, err := json.Marshal(errorResponse)
				Expect(err).NotTo(HaveOccurred())
				Expect(string(jsonData)).To(ContainSubstring("message"))
				Expect(string(jsonData)).To(ContainSubstring("An error occurred"))
			})

			It("should be deserializable from JSON", func() {
				jsonData := `{"message": "An error occurred"}`

				var errorResponse models.ErrorResponse
				err := json.Unmarshal([]byte(jsonData), &errorResponse)
				Expect(err).NotTo(HaveOccurred())
				Expect(errorResponse.Message).To(Equal("An error occurred"))
			})
		})
	})
})

func TestModels(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Models Suite")
}
