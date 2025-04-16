package utils

import (
	"io/ioutil"
	"log"
	"text/template"
)

func CreateSeleniumScript(bbbURL string) string {
	// Create a temporary file for the selenium script
	file, err := ioutil.TempFile("", "selenium_script_*.go")
	if err != nil {
		log.Fatalf("Failed to create temporary file: %v", err)
	}
	defer file.Close()

	// Template for the selenium script
	scriptTemplate := `package main

import (
	"log"
	"time"

	"github.com/tebeka/selenium"
)

func main() {
	BBB_URL := "{{.BBBURL}}"
	TestRecording(nil, BBB_URL)
}

func TestRecording(t testing.TB, BBB_URL string) {
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

	// click listen only button (bigbluebutton session)
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
}`

	// Create a template and execute it with the BBB URL
	tmpl, err := template.New("selenium").Parse(scriptTemplate)
	if err != nil {
		log.Fatalf("Failed to parse template: %v", err)
	}

	data := struct {
		BBBURL string
	}{
		BBBURL: bbbURL,
	}

	err = tmpl.Execute(file, data)
	if err != nil {
		log.Fatalf("Failed to execute template: %v", err)
	}

	return file.Name()
}
