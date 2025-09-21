package initializers_test

import (
	"os"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"spoutbreeze/initializers"
)

var _ = Describe("LoadEnvVariables", func() {
	Describe("Environment variable loading", func() {
		Context("when .env file exists", func() {
			It("should load environment variables without panic", func() {
				Expect(func() {
					initializers.LoadEnvVariables()
				}).NotTo(Panic())
			})
		})

		Context("when .env file does not exist", func() {
			It("should handle missing .env file gracefully", func() {
				originalDir, _ := os.Getwd()
				tempDir, err := os.MkdirTemp("", "test-env")
				Expect(err).NotTo(HaveOccurred())
				defer os.RemoveAll(tempDir)
				defer os.Chdir(originalDir)

				os.Chdir(tempDir)

				Expect(func() {
					initializers.LoadEnvVariables()
				}).NotTo(Panic())
			})
		})
	})
})

func TestInitializers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Initializers Suite")
}
