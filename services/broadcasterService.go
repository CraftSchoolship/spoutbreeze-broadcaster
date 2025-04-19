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
	// Configure Moon options with environment variables
	moonOptions := map[string]interface{}{
		"enableVideo": true,
	}
	
	chromeCaps := chrome.Capabilities{
		ExcludeSwitches: []string{"enable-automation"},
		Args: []string{
			"--start-maximized",
		},
	}
	
	// Define browser capabilities
	caps := selenium.Capabilities{
		"browserName":    "chrome",
		"browserVersion": "133.0.6943.98-6",
		"moon:options":   moonOptions,
		"goog:chromeOptions": chromeCaps,
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

	// Find and click the close button on the popup
	closeButton, err := driver.FindElement(selenium.ByCSSSelector, "button[aria-label='Close Session Details']")
	if err != nil {
		// Try alternate selectors if the first one doesn't work
		closeButton, err = driver.FindElement(selenium.ByCSSSelector, "button.sc-fhzFiK.dIxYkk")
		if err != nil {
			// Try by ID
			closeButton, err = driver.FindElement(selenium.ByID, "tippy-38")
			if err != nil {
				// Try by data-test attribute
				closeButton, err = driver.FindElement(selenium.ByCSSSelector, "button[data-test='closeModal']")
				if err != nil {
					log.Printf("Warning: Failed to find close button: %v", err)
				}
			}
		}
	}
	
	if closeButton != nil {
		err = closeButton.Click()
		if err != nil {
			log.Printf("Warning: Failed to click close button: %v", err)
			// Try using JavaScript to click the button
			_, err = driver.ExecuteScript("arguments[0].click();", []interface{}{closeButton})
			if err != nil {
				log.Printf("Warning: Failed to click close button with JavaScript: %v", err)
			}
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

	// Check session status every 30 seconds
	for {
		if isMeetingRunning() {
			log.Println("Meeting is still running, keeping session alive...")
		} else {
			log.Println("Meeting has ended, terminating session...")
			break
		}
		time.Sleep(20 * time.Second)
	}

	log.Println("Streaming completed successfully")
}