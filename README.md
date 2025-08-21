# Fuzzy

A clean, simple, and efficient Go web server with a home page.

## Features

- **Clean Architecture**: Well-structured Go code with proper error handling
- **Simple Design**: Minimalist approach focused on functionality
- **Efficient Performance**: Lightweight HTTP server using Go's standard library
- **English Documentation**: All code comments and documentation in English
- **Responsive UI**: Clean HTML interface that works on all devices

## Getting Started

### Prerequisites

- Go 1.18 or higher

### Installation

1. Clone the repository:
```bash
git clone https://github.com/0x72EB7BA9B69BF6AB/fuzzy-guide.git
cd fuzzy-guide
```

2. Build the application:
```bash
go build -o fuzey main.go
```

3. Run the server:
```bash
./fuzey
```

Or run directly with Go:
```bash
go run main.go
```

### Usage

Once the server is running, you can access:

- **Home Page**: http://localhost:8080/
- **Health Check**: http://localhost:8080/health

The server will display startup information in the console, including the URLs where the application is accessible.

## API Endpoints

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/` | GET | Home page with welcome message |
| `/health` | GET | Health check endpoint (JSON response) |

## Development

### Running in Development Mode

```bash
go run main.go
```

The server will start on port 8080 by default.

### Building for Production

```bash
go build -ldflags="-s -w" -o fuzey main.go
```

This creates an optimized binary with reduced size.

## Project Structure

```
fuzzy-guide/
├── main.go              # Main server application with routing
├── handlers/            # HTTP request handlers
│   ├── home.go         # Home page handler
│   └── health.go       # Health check handler
├── models/             # Data structures
│   └── page.go         # Page data models
├── templates/          # HTML templates
│   └── home.html      # Home page template
├── go.mod              # Go module file
├── README.md           # Project documentation
└── LICENSE             # License file
```

## License

This project is released into the public domain under the Unlicense. See the [LICENSE](LICENSE) file for details.