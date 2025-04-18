package services

import (
	"log"
	"testing"
	"time"
	"os"
	"net/http"
	"io"
	"encoding/xml"
	"fmt"


	"github.com/tebeka/selenium"
	"spoutbreeze/models"
	"spoutbreeze/repositories"
)

func ProcessBroadcasterRequest(request *models.BroadcasterRequest) error {
	// Store RTMP URL and Stream URL in Redis
	err := repositories.StoreRTMPURL(request.RTMPURL)
	if err != nil {
		return err
	}
	
	err = repositories.StoreStreamKey(request.StreamKey)
	if err != nil {
		return err
	}
	
	// Launch selenium script in the background
	go launchSeleniumScript(request.BBBServerURL,request.BBBHealthCheckURL)
	
	return nil
}

func launchSeleniumScript(bbbURL string ,BBBHealthCheckURL string) {
	StreamBBBSession(nil, bbbURL, BBBHealthCheckURL)
}

func StreamBBBSession(t *testing.T, BBB_URL string, BBBHealthCheckURL string) {
	// Configure Moon options with environment variables
	moonOptions := map[string]interface{}{
		"enableVideo": true,
	}
	
	// Define browser capabilities
	caps := selenium.Capabilities{
		"browserName":    "chrome",
		"browserVersion": "133.0.6943.98-6",
		"moon:options":   moonOptions,
	}
	
	// Get environment variables for Selenium hub URL
    minikubeIP := os.Getenv("MINIKUBE_IP")
    moonPort := os.Getenv("MOON_PORT_4444")
    
    // Use default values if environment variables are not set
    if minikubeIP == "" {
        log.Println("MINIKUBE_IP environment variable not set. Using default:", minikubeIP)
    }
    if moonPort == "" {
        log.Println("MOON_PORT_4444 environment variable not set. Using default:", moonPort)
    }
    
    // Construct the Selenium hub URL
    seleniumHubURL := fmt.Sprintf("http://%s:%s/wd/hub", minikubeIP, moonPort)

	// Connect to Moon server
	driver, err := selenium.NewRemote(caps, seleniumHubURL)
	if err != nil {
		log.Fatalf("Error starting browser: %v", err)
	}
	defer driver.Quit()

	err = driver.MaximizeWindow("")
    if err != nil {
        log.Printf("Warning: Failed to maximize window: %v", err)
        // Alternative approach if maximize doesn't work
        _, err = driver.ExecuteScript("window.resizeTo(screen.width, screen.height);", nil)
        if err != nil {
            log.Printf("Warning: Failed to resize window with JavaScript: %v", err)
        }
    }
	
	// Navigate to BigBlueButton URL
	err = driver.Get(BBB_URL)
	if err != nil {
		log.Fatalf("Failed to navigate to YouTube: %v", err)
	}
	
	// Wait for page load
	time.Sleep(5 * time.Second)
	
	// Handle consent popup if exists
	consentButton, err := driver.FindElement(selenium.ByCSSSelector, "button.ytp-button[aria-label='Accept all']")
	if err == nil {
		consentButton.Click()
		time.Sleep(2 * time.Second)
	}
	
	// Click listen only button (bigbluebutton session)
	listenOnlyButton, err := driver.FindElement(selenium.ByCSSSelector, "button[aria-label='Listen only']")
	if err == nil {
		listenOnlyButton.Click()
		time.Sleep(2 * time.Second)
	}
	
	// End session after the meeting ends
	// Wait for the session to end 
	
	sessionID := driver.SessionID()
	if sessionID == "" {
		log.Fatalf("Failed to retrieve session ID")
	}
	
	// Create HTTP client
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Function to check if meeting is running
	isMeetingRunning := func() bool {
		resp, err := client.Get(BBBHealthCheckURL)
		if err != nil {
			log.Printf("Error checking meeting status: %v", err)
			return false
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Error reading response body: %v", err)
			return false
		}

		// Parse XML response
		var response struct {
			ReturnCode string `xml:"returncode"`
			Running    string `xml:"running"`
		}
		
		err = xml.Unmarshal(body, &response)
		if err != nil {
			log.Printf("Error parsing XML response: %v", err)
			return false
		}
		
		return response.ReturnCode == "SUCCESS" && response.Running == "true"
	}

	// Check session status every 30 seconds
	for {
		if isMeetingRunning() {
			log.Println("Meeting is still running, keeping session alive...")
		} else {
			log.Println("Meeting has ended, terminating session...")
			break
		}
		time.Sleep(30 * time.Second)
	}

	log.Println("Streaming completed successfully")
}