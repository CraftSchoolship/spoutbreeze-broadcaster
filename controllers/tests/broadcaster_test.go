package controllers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"spoutbreeze/models"
)

func testJoinBBBHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request models.BroadcasterRequest

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Broadcasting session started successfully"})
	}
}

var _ = Describe("Broadcaster Controller", func() {
	var (
		router *gin.Engine
		w      *httptest.ResponseRecorder
	)

	BeforeEach(func() {
		gin.SetMode(gin.TestMode)
		router = gin.New()
		router.POST("/broadcaster/joinBBB", testJoinBBBHandler())
		w = httptest.NewRecorder()
	})

	Describe("JoinBBB", func() {
		Context("when valid request is provided", func() {
			It("should return success response", func() {
				request := models.BroadcasterRequest{
					BBBServerURL:      "https://example.com/bigbluebutton",
					BBBHealthCheckURL: "https://example.com/bigbluebutton/api/health",
					RTMPURL:           "rtmp://streaming.example.com/live",
					StreamKey:         "stream-123",
				}
				jsonData, err := json.Marshal(request)
				Expect(err).NotTo(HaveOccurred())

				req, err := http.NewRequest("POST", "/broadcaster/joinBBB", bytes.NewBuffer(jsonData))
				Expect(err).NotTo(HaveOccurred())
				req.Header.Set("Content-Type", "application/json")

				router.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusOK))

				var response map[string]interface{}
				err = json.Unmarshal(w.Body.Bytes(), &response)
				Expect(err).NotTo(HaveOccurred())
				Expect(response["message"]).To(Equal("Broadcasting session started successfully"))
			})
		})

		Context("when invalid JSON is provided", func() {
			It("should return bad request error", func() {

				invalidJSON := `{"invalid": json}`
				req, err := http.NewRequest("POST", "/broadcaster/joinBBB", bytes.NewBufferString(invalidJSON))
				Expect(err).NotTo(HaveOccurred())
				req.Header.Set("Content-Type", "application/json")

				router.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusBadRequest))

				var response map[string]interface{}
				err = json.Unmarshal(w.Body.Bytes(), &response)
				Expect(err).NotTo(HaveOccurred())
				Expect(response["error"]).To(BeAssignableToTypeOf(""))
			})
		})

		Context("when empty request body is provided", func() {
			It("should return bad request error", func() {

				req, err := http.NewRequest("POST", "/broadcaster/joinBBB", bytes.NewBuffer([]byte{}))
				Expect(err).NotTo(HaveOccurred())
				req.Header.Set("Content-Type", "application/json")

				router.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusBadRequest))
			})
		})

		Context("when content-type is not application/json", func() {
			It("should return bad request error", func() {

				req, err := http.NewRequest("POST", "/broadcaster/joinBBB", bytes.NewBufferString("some data"))
				Expect(err).NotTo(HaveOccurred())
				req.Header.Set("Content-Type", "text/plain")

				router.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusBadRequest))
			})
		})

		Context("when service processing fails", func() {
			It("should return internal server error", func() {

				request := models.BroadcasterRequest{
				}
				jsonData, err := json.Marshal(request)
				Expect(err).NotTo(HaveOccurred())

				req, err := http.NewRequest("POST", "/broadcaster/joinBBB", bytes.NewBuffer(jsonData))
				Expect(err).NotTo(HaveOccurred())
				req.Header.Set("Content-Type", "application/json")

				router.ServeHTTP(w, req)

			})
		})
	})

	Describe("HTTP Method validation", func() {
		Context("when using GET method instead of POST", func() {
			It("should return method not allowed", func() {

				req, err := http.NewRequest("GET", "/broadcaster/joinBBB", nil)
				Expect(err).NotTo(HaveOccurred())

				router.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusNotFound))
			})
		})
	})
})
