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
	"github.com/sheva0914/selenium/chrome"
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

	// Get environment variables for Selenium hub URL
    minikubeIP := os.Getenv("MINIKUBE_IP")
    moonPort := os.Getenv("MOON_PORT_4444")
	RedisPassword := os.Getenv("REDIS_PASSWORD")
	// Configure Moon options with environment variables
	moonOptions := map[string]interface{}{
		"enableVideo": false,
		"env": []string{"USER_REDIS_PASSWORD=" + RedisPassword},
	}
	// Configure Chrome options}
	
	chromeCaps := chrome.Capabilities{
		ExcludeSwitches: []string{"enable-automation"},
		Args: []string{
			"--start-maximized",
			"--use-fake-ui-for-media-stream",
			"--use-fake-device-for-media-stream",
			"--autoplay-policy=no-user-gesture-required",
		},
	}
	
	// Define browser capabilities
	caps := selenium.Capabilities{
		"browserName":    "chrome",
		"browserVersion": "0.0.1.3",
		"moon:options":   moonOptions,
		"goog:chromeOptions": chromeCaps,
	}

    
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

	// err = driver.MaximizeWindow("")
    // if err != nil {
    //     log.Printf("Warning: Failed to maximize window: %v", err)
    //     // Alternative approach if maximize doesn't work
    //     _, err = driver.ExecuteScript("window.resizeTo(screen.width, screen.height);", nil)
    //     if err != nil {
    //         log.Printf("Warning: Failed to resize window with JavaScript: %v", err)
    //     }
    // }
	
	// Navigate to BigBlueButton URL
	err = driver.Get(BBB_URL)
	if err != nil {
		log.Fatalf("Failed to navigate to BigBlueButton: %v", err)
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

	// Find and click the close button on the popup
	closeButton, err := driver.FindElement(selenium.ByCSSSelector, "button[aria-label='Close Session Details']")
	if err != nil {
		log.Printf("Warning: Failed to find close button: %v", err)
	} else {
		_, err = driver.ExecuteScript("arguments[0].click();", []interface{}{closeButton})
			if err != nil {
				log.Printf("Warning: Failed to click close button with JavaScript: %v", err)
			}
		time.Sleep(2 * time.Second)
	}

	// Find Users and messages close button
	usersAndMessagesButton, err := driver.FindElement(selenium.ByCSSSelector, "button[aria-label='Users and messages toggle']")
	if err != nil {
		log.Printf("Warning: Failed to find Users and messages button: %v", err)
	} else {
		_, err = driver.ExecuteScript("arguments[0].click();", []interface{}{usersAndMessagesButton})
		if err != nil {
			log.Printf("Warning: Failed to click Users and messages button with JavaScript: %v", err)
		}
		time.Sleep(2 * time.Second)
	}

	
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

	var meetingWasRunning bool
	// Check session status every 20 seconds
	for {
		meetingRunning := isMeetingRunning()
		if meetingRunning != meetingWasRunning {
			if meetingRunning {
				log.Println("Meeting is still running, keeping session alive...")
			} else {
				log.Println("Meeting has ended, terminating session...")
			}
			meetingWasRunning = meetingRunning
		}
		time.Sleep(20 * time.Second)
	}
	log.Println("Streaming completed successfully")
}