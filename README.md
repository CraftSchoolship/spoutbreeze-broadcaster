# SpoutBreeze BBB Broadcaster

## Project Description

SpoutBreeze BBB Broadcaster is a Go-based backend service that enables broadcasting BigBlueButton (BBB) sessions via RTMP streams. It provides a simple REST API endpoint that accepts streaming parameters and automatically joins a BBB session as a moderator to start broadcasting.

The service uses Selenium WebDriver to automate browser interactions with BigBlueButton sessions, making it ideal for automated broadcasting, recording, or stream redistribution workflows.

## Features

- RESTful API endpoint for initiating BBB broadcasting sessions
- Automated browser interaction with BigBlueButton using Selenium
- Containerization support with Docker and Kubernetes/Minikube
- Environment variable configuration

## Technology Stack

- **Backend**: Go (Golang) with Gin framework
- **Browser Automation**: Selenium WebDriver with Chrome
- **Container Orchestration**: Kubernetes/Minikube
- **Deployment**: Docker

## Project Structure

```
spoutbreeze/
├── .env                         # Environment variables
├── main.go                      # Entry point that initializes and runs the server
├── controllers/
│   └── broadcasterController.go # Handles HTTP requests
├── initializers/
│   └── loadEnvVariables.go      # Loads environment variables
├── models/
│   └── broadcaster.go           # Data structures for the API
├── routes/
│   └── routes.go                # API route definitions
├── services/
│   └── broadcasterService.go    # Business logic for BBB streaming
└── Dockerfile                   # Docker configuration for containerization
```

## Prerequisites

- Go 1.18 or later
- Kubernetes/Minikube with Moon Selenium Grid installed
- Chrome browser in Moon Selenium Grid

## Installation

### 1. Clone the repository

```bash
git clone https://github.com/Bouchiba43/spoutbreeze-rtmp-svc.git
cd spoutbreeze
```

### 2. Set up environment variables

Create a `.env` file in the project root:

```
# The IP address of your Minikube or Kubernetes cluster, used to access services externally
CLUSTER_IP=your_minikube_or_k8s_cluster_ip

# The port on which the Moon Selenium Hub is exposed (e.g., 4444)
MOON_SELENIUM_PORT=4444

# Default gin server port
PORT=1323
```

### 3. Install dependencies

```bash
go mod tidy
```

### 4. Build the application

```bash
go build -o spoutbreeze
```

### 5. Run the application

```bash
./spoutbreeze
```

## Testing

The project uses Ginkgo and Gomega testing frameworks for comprehensive testing of the application.

### Installing Test Dependencies

```bash
# Install Ginkgo and Gomega
go get -u github.com/onsi/ginkgo/v2/ginkgo
go get -u github.com/onsi/gomega
go install github.com/onsi/ginkgo/v2/ginkgo@latest
```

### Running Tests

The project includes several testing options via the Makefile:

```bash
# Run all standard tests (excluding Selenium-dependent tests)
make test

# Run tests with verbose output
make test-verbose

# Run only controller tests 
make test-controllers

# Run tests with code coverage report
make test-coverage

# Continuously run tests when files change (watch mode)
make test-watch

# Run tests for CI environments with JUnit reports
make test-ci
```

### Test Coverage

To generate a test coverage report:

```bash
make test-coverage
```

This will create an HTML coverage report (`coverage.html`) that you can view in any browser.

### Test Structure

The tests are organized by component:

1. **Controller Tests**: Test the API endpoints and request handling
   - Located in `controllers/tests/`
   - Use mocked services to isolate from external dependencies

2. **Service Tests**: Test the business logic
   - Located in `services/tests/`
   - Some tests may require Selenium and are skipped by default

3. **Integration Tests**: Test the full flow from request to processing
   - Located at project root level

### Writing New Tests

To generate new test files for a package:

```bash
# Replace 'package_name' with the target package
ginkgo generate package_name
```

For more complex packages that require test suites:

```bash
cd package_name
ginkgo bootstrap
ginkgo generate file_name
```

## Docker Deployment

### 1. Build the Docker image

```bash
docker build -t spoutbreeze:latest .
```

### 2. Run the container

```bash
docker run -p 8080:8080 --env-file .env spoutbreeze:latest
```

## API Documentation

### Join BBB Session for Broadcasting

Initiates a BigBlueButton session join and prepares for broadcasting.

**Endpoint:** `POST /broadcaster/joinBBB`

**Request Body:**

```json
{
  "bbb_server_url": "https://bbb-server.com/bigbluebutton/api/join?fullName=User&meetingID=meeting1&password=mp&redirect=true&checksum=abcdef123456",
  "rtmp_url": "rtmp://streaming-server.com/live/{$stream-key}",
  "stream_url": "stream-key"
}
```

**Parameters:**
- `bbb_server_url` (string, required): The BigBlueButton server URL with join parameters and checksum
- `rtmp_url` (string, required): RTMP URL for streaming
- `stream_url` (string, required): Public stream URL for viewers

**Response:**
- Success (200 OK):
  ```json
  {
    "message": "Broadcasting session started successfully"
  }
  ```
- Error (400 Bad Request):
  ```json
  {
    "error": "error message"
  }
  ```
- Error (500 Internal Server Error):
  ```json
  {
    "error": "error message"
  }
  ```

## Implementation Details

### Key Components

#### 1. Initializers

**Environment Variables:**
- Loads configuration from the `.env` file

#### 2. Models

Defines the data structure for the broadcaster request:
- `BBBServerURL`: URL to join the BigBlueButton session
- `RTMPURL`: RTMP URL for streaming
- `StreamURL`: Public stream URL

#### 3. Controllers

Handles HTTP requests to the API endpoint:
- Validates the incoming JSON request
- Calls the service layer to process the request

#### 4. Services

Contains the core business logic:
- Launches a Selenium script in the background to join the BBB session
- Implements the `StreamBBBSession` function that automates the browser interaction with BBB

#### 5. Routes

Defines API routes:
- Configures the `/broadcaster/joinBBB` endpoint

### Selenium Automation

The `StreamBBBSession` function performs the following actions:
1. Connects to the Moon Selenium Grid on Minikube
2. Launches a Chrome browser with video capability enabled
3. Navigates to the BBB session URL
4. Waits for the page to load
5. Handles any consent popups
6. Clicks the "Listen only" button in the BBB session
7. Maintains the session for 10 minutes (600 seconds)

## Troubleshooting

### Common Issues

#### 1. Selenium Connection Issues
**Symptoms:** "Error starting browser" messages.
**Solution:**
- Verify the Moon Selenium Grid is running
- Check the Moon service endpoint (e.g., 192.168.49.2:32440)
- Ensure Chrome browser is available in the Moon container

```bash
kubectl get pods -n moon
kubectl logs <moon-pod-name> -n moon
```

#### 2. BigBlueButton Join Problems
**Symptoms:** "Failed to navigate" or no "Listen only" button found.
**Solution:**
- Verify the BBB URL is valid and contains required parameters
- Check BBB server health
- Verify the session is actually running on the BBB server

#### 3. Test Failures
**Symptoms:** Tests failing with connection errors to Selenium.
**Solution:**
- For basic tests, use `make test` which skips Selenium-dependent tests
- For Selenium tests, ensure the Moon Selenium Grid is running and accessible
- Mock external services in tests using the provided test helpers

```bash
# Check if the test environment variables are set correctly
cat .env
# Run tests that don't rely on external dependencies
make test-controllers
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## Acknowledgments

- [Gin Web Framework](https://github.com/gin-gonic/gin)
- [Selenium](https://github.com/tebeka/selenium)
- [BigBlueButton](https://bigbluebutton.org/)
