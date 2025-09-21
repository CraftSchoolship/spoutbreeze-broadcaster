package repositories_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"spoutbreeze/repositories"
)

var _ = Describe("Redis Repository", func() {
	Describe("StoreRTMPURL", func() {
		Context("when Redis client is not initialized", func() {
			It("should panic with nil client", func() {

				Expect(func() {
					repositories.StoreRTMPURL("rtmp://test.url/live")
				}).To(Panic())
			})
		})

		Context("when given valid RTMP URL", func() {
			It("should panic when trying to store the URL", func() {
				rtmpURL := "rtmp://streaming.example.com/live"

				Expect(func() {
					repositories.StoreRTMPURL(rtmpURL)
				}).To(Panic())
			})
		})

		Context("when given empty RTMP URL", func() {
			It("should panic when trying to store empty string", func() {
				Expect(func() {
					repositories.StoreRTMPURL("")
				}).To(Panic())
			})
		})
	})

	Describe("StoreStreamKey", func() {
		Context("when Redis client is not initialized", func() {
			It("should panic with nil client", func() {

				Expect(func() {
					repositories.StoreStreamKey("test-stream-key")
				}).To(Panic())
			})
		})

		Context("when given valid stream key", func() {
			It("should panic when trying to store the stream key", func() {
				streamKey := "stream-123-abc"
				Expect(func() {
					repositories.StoreStreamKey(streamKey)
				}).To(Panic())
			})
		})

		Context("when given empty stream key", func() {
			It("should panic when trying to store empty string", func() {
				Expect(func() {
					repositories.StoreStreamKey("")
				}).To(Panic())
			})
		})
	})

	Describe("Function signatures", func() {
		It("should have correct StoreRTMPURL function signature", func() {

			var fn func(string) error = repositories.StoreRTMPURL
			Expect(fn).NotTo(BeNil())
		})

		It("should have correct StoreStreamKey function signature", func() {

			var fn func(string) error = repositories.StoreStreamKey
			Expect(fn).NotTo(BeNil())
		})
	})
})

func TestRepositories(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Repositories Suite")
}
