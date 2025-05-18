# Chirpy

A simple social network app built in Go to refresh Golang skills.

## Overview

Chirpy is a web application that serves static content and provides several API endpoints for health checks and metrics.

## Prerequisites

- Go 1.24+ for local development
- Docker and Docker Compose for containerized deployment

## Running Locally

1. Clone the repository
2. Navigate to the project directory
3. Run the application:

```bash
go run main.go
```

The server will start on port 8080.

## Running with Docker

1. Build and start the container:

```bash
docker compose up --build
```

2. The application will be available at `http://localhost:8080`

## API Endpoints

- `/app/` - Serves the static web interface
- `/healthz` - Health check endpoint (returns "OK")
- `/metrics` - View hit counter metrics
- `/reset` - Reset the hit counter

## Project Structure

```
.
├── app/            # Static files served by the application
│   └── index.html  # Main landing page
├── Dockerfile      # Container definition
├── docker-compose.yml
├── go.mod          # Go module definition
├── main.go         # Application entry point
└── README.md       # This file
```

## Development Notes

The application uses middleware chains for request processing, including logging, header addition, and metrics tracking. The file server serves content from the root directory.