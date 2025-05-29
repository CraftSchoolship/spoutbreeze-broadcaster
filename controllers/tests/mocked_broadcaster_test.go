package controllers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"spoutbreeze/models"
)

type BroadcasterServiceInterface interface {
	ProcessBroadcasterRequest(request *models.BroadcasterRequest) error
}

type MockBroadcasterService struct {
	ShouldFail bool
	Error      error
}

func (m *MockBroadcasterService) ProcessBroadcasterRequest(request *models.BroadcasterRequest) error {
	if m.ShouldFail {
		return m.Error
	}
	return nil
}

func joinBBBHandler(service BroadcasterServiceInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request models.BroadcasterRequest

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := service.ProcessBroadcasterRequest(&request)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Broadcasting session started successfully"})
	}
}

var _ = Describe("Broadcaster Controller (with mocking)", func() {
	var (
		router      *gin.Engine
		w           *httptest.ResponseRecorder
		mockService *MockBroadcasterService
	)

	BeforeEach(func() {
		gin.SetMode(gin.TestMode)
		router = gin.New()
		mockService = &MockBroadcasterService{}

		router.POST("/broadcaster/joinBBB", joinBBBHandler(mockService))
		w = httptest.NewRecorder()
	})

	Describe("JoinBBB with service mocking", func() {
		Context("when service returns success", func() {
			It("should return success response", func() {
	
				mockService.ShouldFail = false
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

		Context("when service returns error", func() {
			It("should return internal server error", func() {
				mockService.ShouldFail = true
				mockService.Error = errors.New("service processing failed")

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

				Expect(w.Code).To(Equal(http.StatusInternalServerError))

				var response map[string]interface{}
				err = json.Unmarshal(w.Body.Bytes(), &response)
				Expect(err).NotTo(HaveOccurred())
				Expect(response["error"]).To(ContainSubstring("service processing failed"))
			})
		})
	})

	DescribeTable("JoinBBB with various inputs",
		func(requestBody string, expectedStatus int, expectedContains string) {
			req, err := http.NewRequest("POST", "/broadcaster/joinBBB", bytes.NewBufferString(requestBody))
			Expect(err).NotTo(HaveOccurred())
			req.Header.Set("Content-Type", "application/json")

			router.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(expectedStatus))
			if expectedContains != "" {
				Expect(w.Body.String()).To(ContainSubstring(expectedContains))
			}
		},
		Entry("valid request", `{"bbb_server_url":"https://example.com/bigbluebutton","bbb_health_check_url":"https://example.com/bigbluebutton/api/health","rtmp_url":"rtmp://streaming.example.com/live","stream_key":"stream-123"}`, http.StatusOK, "successfully"),
		Entry("empty request", `{}`, http.StatusBadRequest, "error"),
		Entry("invalid JSON", `{"invalid": json}`, http.StatusBadRequest, "error"),
		Entry("malformed JSON", `{invalid}`, http.StatusBadRequest, "error"),
	)
})
