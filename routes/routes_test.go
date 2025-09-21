package routes_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"spoutbreeze/routes"
)

var _ = Describe("Routes", func() {
	var (
		router *gin.Engine
		w      *httptest.ResponseRecorder
	)

	BeforeEach(func() {
		gin.SetMode(gin.TestMode)
		router = routes.SetupRouter()
		w = httptest.NewRecorder()
	})

	Describe("SetupRouter", func() {
		Context("when setting up routes", func() {
			It("should create a valid Gin engine", func() {
				Expect(router).NotTo(BeNil())
				Expect(router).To(BeAssignableToTypeOf(&gin.Engine{}))
			})

			It("should have broadcaster routes configured", func() {
				req, err := http.NewRequest("POST", "/broadcaster/joinBBB", nil)
				Expect(err).NotTo(HaveOccurred())

				router.ServeHTTP(w, req)

				Expect(w.Code).NotTo(Equal(http.StatusNotFound))
			})

			It("should have health check routes configured", func() {

				req, err := http.NewRequest("GET", "/health", nil)
				Expect(err).NotTo(HaveOccurred())

				router.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusOK))
			})

			It("should have readiness check routes configured", func() {

				req, err := http.NewRequest("GET", "/readiness", nil)
				Expect(err).NotTo(HaveOccurred())

				router.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusOK))
			})

			It("should return 404 for undefined routes", func() {
				req, err := http.NewRequest("GET", "/undefined", nil)
				Expect(err).NotTo(HaveOccurred())

				router.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusNotFound))
			})
		})

		Context("when testing route groups", func() {
			It("should have broadcaster group configured properly", func() {

				req, err := http.NewRequest("GET", "/broadcaster/undefined", nil)
				Expect(err).NotTo(HaveOccurred())

				router.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusNotFound))
			})

			It("should handle different HTTP methods appropriately", func() {

				req, err := http.NewRequest("GET", "/broadcaster/joinBBB", nil)
				Expect(err).NotTo(HaveOccurred())

				router.ServeHTTP(w, req)

				Expect(w.Code).To(BeElementOf([]int{http.StatusMethodNotAllowed, http.StatusNotFound}))
			})
		})
	})
})

func TestRoutes(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Routes Suite")
}
