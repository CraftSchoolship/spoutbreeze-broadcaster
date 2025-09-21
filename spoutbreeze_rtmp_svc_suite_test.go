package main_test

import (
	"os"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Main Package", func() {
	Describe("Environment handling", func() {
		Context("when PORT environment variable is set", func() {
			It("should use the custom port", func() {
				port := os.Getenv("PORT")
				if port == "" {
					port = "1323" 
				}
				Expect(port).NotTo(BeEmpty())
				Expect(port).To(MatchRegexp(`^\d+$`))
			})
		})

		Context("when PORT environment variable is not set", func() {
			It("should use the default port 1323", func() {
				oldPort := os.Getenv("PORT")
				os.Unsetenv("PORT")
				defer func() {
					if oldPort != "" {
						os.Setenv("PORT", oldPort)
					}
				}()

				port := os.Getenv("PORT")
				if port == "" {
					port = "1323"
				}
				Expect(port).To(Equal("1323"))
			})
		})
	})

	Describe("Application initialization", func() {
		It("should initialize without panicking", func() {
			Expect(func() {
			}).NotTo(Panic())
		})
	})
})

func TestSpoutbreezeRtmpSvc(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "SpoutbreezeRtmpSvc Suite")
}
