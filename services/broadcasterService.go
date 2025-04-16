package services

import (
	"log"
	"testing"
	"time"

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
	
	err = repositories.StoreStreamURL(request.StreamURL)
	if err != nil {
		return err
	}
	
	// Launch selenium script in the background
	go launchSeleniumScript(request.BBBServerURL)
	
	return nil
}

func launchSeleniumScript(bbbURL string) {
	StreamBBBSession(nil, bbbURL)
}

func StreamBBBSession(t *testing.T, BBB_URL string) {
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
	
	// Connect to Moon server
	driver, err := selenium.NewRemote(caps, "http://192.168.49.2:32440/wd/hub")
	if err != nil {
		log.Fatalf("Error starting browser: %v", err)
	}
	defer driver.Quit()
	
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
	
	time.Sleep(600 * time.Second)
	
	sessionID := driver.SessionID()
	if sessionID == "" {
		log.Fatalf("Failed to retrieve session ID")
	}
	
	log.Println("Streaming completed successfully")
}