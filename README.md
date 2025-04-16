# SpoutBreeze BBB Broadcaster

## Project Description

SpoutBreeze BBB Broadcaster is a Go-based backend service that enables broadcasting BigBlueButton (BBB) sessions via RTMP streams. It provides a simple REST API endpoint that accepts streaming parameters and automatically joins a BBB session as a moderator to start broadcasting.

The service integrates with Redis for storing stream configuration data and uses Selenium WebDriver to automate browser interactions with BigBlueButton sessions, making it ideal for automated broadcasting, recording, or stream redistribution workflows.

## Features

- RESTful API endpoint for initiating BBB broadcasting sessions
- Redis integration for storing stream configuration
- Automated browser interaction with BigBlueButton using Selenium
- Containerization support with Docker and Kubernetes/Minikube
- Environment variable configuration

## Technology Stack

- **Backend**: Go (Golang) with Gin framework
- **Storage**: Redis
- **Browser Automation**: Selenium WebDriver with Chrome
- **Container Orchestration**: Kubernetes/Minikube
- **Deployment**: Docker

## Project Structure

```
spoutbreeze/
├── .env                         # Environment variables including Redis password
├── main.go                      # Entry point that initializes and runs the server
├── controllers/
│   └── broadcasterController.go # Handles HTTP requests
├── initializers/
│   ├── loadEnvVariables.go      # Loads environment variables
│   └── redis.go                 # Sets up Redis connection
├── models/
│   └── broadcaster.go           # Data structures for the API
├── repositories/
│   └── redisRepository.go       # Data access layer for Redis
├── routes/
│   └── routes.go                # API route definitions
├── services/
│   └── broadcasterService.go    # Business logic for BBB streaming
└── Dockerfile                   # Docker configuration for containerization
```

## Prerequisites

- Go 1.18 or later
- Redis server (running on Kubernetes/Minikube at 192.168.49.2:6379)
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
REDIS_PASSWORD=your_redis_password_here
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

## Docker Deployment

### 1. Build the Docker image

```bash
docker build -t spoutbreeze:latest .
```

### 2. Run the container

```bash
docker run -p 8080:8080 --env-file .env spoutbreeze:latest
```

## Kubernetes/Minikube Deployment

### 1. Create Kubernetes deployment file

Create a file named `deployment.yaml`:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: spoutbreeze
spec:
  replicas: 1
  selector:
    matchLabels:
      app: spoutbreeze
  template:
    metadata:
      labels:
        app: spoutbreeze
    spec:
      containers:
      - name: spoutbreeze
        image: spoutbreeze:latest
        ports:
        - containerPort: 8080
        env:
        - name: REDIS_PASSWORD
          valueFrom:
            secretKeyRef:
              name: redis-credentials
              key: password
---
apiVersion: v1
kind: Service
metadata:
  name: spoutbreeze
spec:
  selector:
    app: spoutbreeze
  ports:
  - port: 8080
    targetPort: 8080
  type: LoadBalancer
```

### 2. Create a secret for Redis credentials

```bash
kubectl create secret generic redis-credentials --from-literal=password=your_redis_password_here
```

### 3. Deploy to Kubernetes

```bash
kubectl apply -f deployment.yaml
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

**Redis Connection:**
- Establishes a connection to the Redis server running on Minikube (192.168.49.2:6379)
- Uses the password from environment variables for authentication

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
- Stores RTMP and Stream URLs in Redis
- Launches a Selenium script in the background to join the BBB session
- Implements the `StreamBBBSession` function that automates the browser interaction with BBB

#### 5. Repositories

Manages data storage operations:
- Handles Redis SET operations for storing stream configuration

#### 6. Routes

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

#### 1. Redis Connection Failure
**Symptoms:** Error messages about Redis connection failures.
**Solution:** 
- Verify the Redis service is running in Minikube
- Check the Redis password in the `.env` file
- Ensure network connectivity to the Minikube Redis service

```bash
kubectl get svc -n redis
```

#### 2. Selenium Connection Issues
**Symptoms:** "Error starting browser" messages.
**Solution:**
- Verify the Moon Selenium Grid is running
- Check the Moon service endpoint (192.168.49.2:32440)
- Ensure Chrome browser is available in the Moon container

```bash
kubectl get pods -n moon
kubectl logs <moon-pod-name> -n moon
```

#### 3. BigBlueButton Join Problems
**Symptoms:** "Failed to navigate" or no "Listen only" button found.
**Solution:**
- Verify the BBB URL is valid and contains required parameters
- Check BBB server health
- Verify the session is actually running on the BBB server

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgments

- [Gin Web Framework](https://github.com/gin-gonic/gin)
- [go-redis](https://github.com/redis/go-redis)
- [Selenium](https://github.com/tebeka/selenium)
- [BigBlueButton](https://bigbluebutton.org/)