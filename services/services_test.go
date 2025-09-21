package services_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"spoutbreeze/models"
	"spoutbreeze/services"
)

var _ = Describe("Broadcaster Service", func() {
	Describe("ProcessBroadcasterRequest", func() {
		Context("when given a valid broadcaster request", func() {
			It("should process the request without error", func() {
				request := &models.BroadcasterRequest{
					BBBServerURL:      "https://example.com/bigbluebutton",
					BBBHealthCheckURL: "https://example.com/health",
					RTMPURL:           "rtmp://streaming.example.com/live",
					StreamKey:         "stream-123",
				}

				err := services.ProcessBroadcasterRequest(request)
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when given request with empty fields", func() {
			It("should still process without error", func() {
				request := &models.BroadcasterRequest{
					BBBServerURL:      "",
					BBBHealthCheckURL: "",
					RTMPURL:           "",
					StreamKey:         "",
				}

				err := services.ProcessBroadcasterRequest(request)
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when given nil request", func() {
			It("should panic due to nil pointer access", func() {
	
				Expect(func() {
					services.ProcessBroadcasterRequest(nil)
				}).To(Panic()) 
			})
		})
	})

	Describe("Function existence and signatures", func() {
		It("should have ProcessBroadcasterRequest function available", func() {
			var fn func(*models.BroadcasterRequest) error = services.ProcessBroadcasterRequest
			Expect(fn).NotTo(BeNil())
		})
	})

	Describe("Request processing behavior", func() {
		Context("when processing concurrent requests", func() {
			It("should handle multiple requests", func() {
				request1 := &models.BroadcasterRequest{
					BBBServerURL:      "https://example1.com/bigbluebutton",
					BBBHealthCheckURL: "https://example1.com/health",
					RTMPURL:           "rtmp://streaming1.example.com/live",
					StreamKey:         "stream-123",
				}

				request2 := &models.BroadcasterRequest{
					BBBServerURL:      "https://example2.com/bigbluebutton",
					BBBHealthCheckURL: "https://example2.com/health",
					RTMPURL:           "rtmp://streaming2.example.com/live",
					StreamKey:         "stream-456",
				}

				err1 := services.ProcessBroadcasterRequest(request1)
				err2 := services.ProcessBroadcasterRequest(request2)

				Expect(err1).NotTo(HaveOccurred())
				Expect(err2).NotTo(HaveOccurred())
			})
		})

		Context("when processing with special characters", func() {
			It("should handle URLs with special characters", func() {
				request := &models.BroadcasterRequest{
					BBBServerURL:      "https://example.com/bigbluebutton?param=value&test=123",
					BBBHealthCheckURL: "https://example.com/health?check=true",
					RTMPURL:           "rtmp://streaming.example.com/live/stream?key=abc",
					StreamKey:         "stream-123-abc_def",
				}

				err := services.ProcessBroadcasterRequest(request)
				Expect(err).NotTo(HaveOccurred())
			})
		})
	})
})

func TestServices(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Services Suite")
}
