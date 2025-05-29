package controllers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"spoutbreeze/controllers"
)

var _ = Describe("Health Controller", func() {
	var (
		router          *gin.Engine
		healthController *controllers.HealthController
		w               *httptest.ResponseRecorder
	)

	BeforeEach(func() {
		gin.SetMode(gin.TestMode)
		router = gin.New()
		healthController = controllers.NewHealthController()
		
		// Setup routes
		router.GET("/health", healthController.HealthCheck)
		router.GET("/readiness", healthController.ReadinessCheck)
		
		w = httptest.NewRecorder()
	})

	Describe("HealthCheck", func() {
		Context("when health endpoint is called", func() {
			It("should return healthy status", func() {
				req, err := http.NewRequest("GET", "/health", nil)
				Expect(err).NotTo(HaveOccurred())

				router.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusOK))
				Expect(w.Header().Get("Content-Type")).To(ContainSubstring("application/json"))

				var response map[string]interface{}
				err = json.Unmarshal(w.Body.Bytes(), &response)
				Expect(err).NotTo(HaveOccurred())
				Expect(response["status"]).To(Equal("healthy"))
			})
		})

		Context("when using POST method", func() {
			It("should return method not allowed", func() {
				req, err := http.NewRequest("POST", "/health", nil)
				Expect(err).NotTo(HaveOccurred())

				router.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusNotFound)) 
			})
		})

		Context("when using invalid path", func() {
			It("should return not found", func() {
				req, err := http.NewRequest("GET", "/health/invalid", nil)
				Expect(err).NotTo(HaveOccurred())

				router.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusNotFound))
			})
		})
	})

	Describe("ReadinessCheck", func() {
		Context("when readiness endpoint is called", func() {
			It("should return ready status", func() {
				req, err := http.NewRequest("GET", "/readiness", nil)
				Expect(err).NotTo(HaveOccurred())

				router.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusOK))
				Expect(w.Header().Get("Content-Type")).To(ContainSubstring("application/json"))

				var response map[string]interface{}
				err = json.Unmarshal(w.Body.Bytes(), &response)
				Expect(err).NotTo(HaveOccurred())
				Expect(response["status"]).To(Equal("ready"))
			})
		})

		Context("when using POST method", func() {
			It("should return method not allowed", func() {
				req, err := http.NewRequest("POST", "/readiness", nil)
				Expect(err).NotTo(HaveOccurred())

				router.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusNotFound)) 
			})
		})

		Context("when controller is properly initialized", func() {
			It("should create a new health controller instance", func() {
				controller := controllers.NewHealthController()

				Expect(controller).NotTo(BeNil())
				Expect(controller).To(BeAssignableToTypeOf(&controllers.HealthController{}))
			})
		})
	})

	Describe("Response format validation", func() {
		Context("when health endpoint is called", func() {
			It("should return valid JSON structure", func() {
				req, err := http.NewRequest("GET", "/health", nil)
				Expect(err).NotTo(HaveOccurred())

				router.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusOK))
				
				var response map[string]interface{}
				err = json.Unmarshal(w.Body.Bytes(), &response)
				Expect(err).NotTo(HaveOccurred())
				
				Expect(response).To(HaveKey("status"))
				Expect(len(response)).To(Equal(1)) 
			})
		})

		Context("when readiness endpoint is called", func() {
			It("should return valid JSON structure", func() {
				req, err := http.NewRequest("GET", "/readiness", nil)
				Expect(err).NotTo(HaveOccurred())

				router.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusOK))
				
				var response map[string]interface{}
				err = json.Unmarshal(w.Body.Bytes(), &response)
				Expect(err).NotTo(HaveOccurred())
				
				Expect(response).To(HaveKey("status"))
				Expect(len(response)).To(Equal(1)) 
			})
		})
	})
})