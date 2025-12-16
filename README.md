# Shortlink Service

A simple and efficient URL shortener service written in Go.

## Features

- ✅ Create short links from long URLs
- ✅ Redirect short links to original URLs
- ✅ Health check endpoint
- ✅ Background service management with Makefile
- ✅ Graceful shutdown
- ✅ In-memory storage
- ✅ Request logging

## Tech Stack

- **Go 1.22.0** - Programming language
- **Standard Library** - No external dependencies for core functionality
- **Testify** - Testing framework

## Installation

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd shortlink
   ```

2. **Install dependencies**
   ```bash
   make install
   ```

## Running the Service

### Direct Run
```bash
make run
```

### Background Start
```bash
make start
```

### Stop Service
```bash
make stop
```

### Restart Service
```bash
make restart
```

### Build Application
```bash
make build
```

### View Logs
```bash
make logs
```

### Clean Build Artifacts
```bash
make clean
```

## API Documentation

### Create Short Link

**Endpoint:** `POST /api/links`

**Request Body:**
```json
{
  "long_url": "https://example.com"
}
```

**Response:**
```json
{
  "short_code": "abc123"
}
```

**Status Codes:**
- `201 Created` - Short link created successfully
- `400 Bad Request` - Invalid request body
- `405 Method Not Allowed` - Only POST method is allowed
- `500 Internal Server Error` - Server error

### Redirect Short Link

**Endpoint:** `GET /{short_code}`

**Response:**
- `302 Found` - Redirect to original URL
- `404 Not Found` - Short code not found
- `405 Method Not Allowed` - Only GET method is allowed
- `500 Internal Server Error` - Server error

### Health Check

**Endpoint:** `GET /healthz`

**Response:**
```json
{
  "status": "ok"
}
```

**Status Codes:**
- `200 OK` - Service is healthy

## Usage Examples

### Create Short Link with cURL
```bash
curl -X POST -H "Content-Type: application/json" -d '{"long_url": "https://example.com"}' http://localhost:8080/api/links
```

### Redirect Short Link
```bash
curl -L http://localhost:8080/{short_code}
```

### Check Health
```bash
curl http://localhost:8080/healthz
```

## Project Structure

```
shortlink/
├── internal/
│   ├── api/
│   │   └── http/
│   │       ├── handler/        # HTTP request handlers
│   │       │   ├── handler.go
│   │       │   └── handler_test.go
│   │       └── server/         # HTTP server configuration
│   │           └── server.go
│   ├── config/                 # Configuration management
│   │   └── config.go
│   ├── idgen/                  # Short code generator
│   │   └── simple_hash_generator.go
│   ├── shortener/              # Shortlink business logic
│   │   └── service.go
│   └── storage/                # Storage implementation
│       └── memory_store.go
├── go.mod                      # Go module file
├── main.go                     # Application entry point
├── Makefile                    # Build and management scripts
└── README.md                   # This file
```

## Configuration

The service uses a simple configuration structure defined in `internal/config/config.go`. By default, the server listens on port `8080`.

## Development

### Run Tests

```bash
go test ./internal/api/http/handler -v
```

### Code Structure

- **API Layer**: Handles HTTP requests and responses
- **Service Layer**: Implements business logic
- **Storage Layer**: Manages data persistence
- **ID Generator**: Creates unique short codes
- **Config Layer**: Manages application configuration

## Logging

The service logs request information, errors, and service status to standard output. Logs include:
- Request method and path
- Client IP address
- Response status codes
- Error details
- Service start/stop events

## Graceful Shutdown

The service handles SIGILL and SIGTERM signals for graceful shutdown, allowing in-flight requests to complete before exiting.

## License

MIT License

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
