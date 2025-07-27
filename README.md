
![Version](https://img.shields.io/github/v/release/gusdeyw/pgbridge-go?label=version)
![Go Version](https://img.shields.io/badge/go-1.18%2B-blue)
![License](https://img.shields.io/github/license/gusdeyw/pgbridge-go)
![Docker](https://img.shields.io/badge/docker-ready-blue)
![Build Status](https://img.shields.io/badge/build-passing-brightgreen)


# PGBridge-Go


## Version

This project follows [semantic versioning](https://semver.org/). See the badge above for the latest release.

PGBridge-Go is a lightweight, production-ready payment gateway bridge built with Go. It provides a simple, secure, and extensible way to integrate with payment providers, handle payment notifications, and manage payment flows. The project is containerized with Docker and includes a modern logging system, QR code generation, and a clean HTML UI for payment callbacks.

## Features

- ⚡ Fast and efficient Go backend (Fiber framework)
- 🔒 Secure authentication and payment flow
- 📦 Dockerized for easy deployment
- 📝 Structured logging with Zap
- 📄 Swagger UI for API documentation
- 🖼️ QR code generation for payment links
- 📬 Webhook/callback handling
- 🧩 Modular code structure for easy extension


## Supported Payment Gateways

- ✅ Midtrans
- ⏳ Xendit, DOKU, iPaymu, and more gateways coming soon (in progress)

## Tech Stack

- **Backend:** Go (Fiber)
- **Logging:** Zap
- **API Docs:** Swagger UI
- **Containerization:** Docker, Docker Compose
- **Web Server:** Nginx (for static/docs)
- **Database:** (Pluggable, add your own in `/src/database`)
- **QR Code:** github.com/skip2/go-qrcode

## Project Structure

```
.
├── src/
│   ├── main.go                # Entry point
│   ├── controllers/           # Business logic (auth, payment, callbacks, etc.)
│   ├── helper/                # Utility functions (auth, QR, etc.)
│   ├── logger/                # Zap logger setup
│   ├── models/                # Data models
│   ├── routes/                # API routes
│   ├── views/                 # HTML templates
│   ├── config/                # Configuration
│   └── ...
├── nginx/                     # Nginx config for static/docs
├── swagger/                   # Swagger UI and OpenAPI spec
├── docker-compose.yaml
├── README.md
└── LICENSE
```

## Getting Started


### Prerequisites

- [Go](https://golang.org/dl/) 1.18+
- [Docker](https://www.docker.com/) (optional, for containerized production)
- [Docker Compose](https://docs.docker.com/compose/) (optional)

### Running Without Docker (Development)

1. Copy `.env.example` to `.env` and adjust as needed.
2. Install Go dependencies:
   ```sh
   go mod download
   ```
3. Run the application:
   ```sh
   go run main.go
   ```
   or, for hot reload during development:
   ```sh
   CompileDaemon --command="go run main.go"
   ```
4. The backend will be available at `http://localhost:5000`.


### Running with Docker (Development)

Use the Docker Compose file in the root folder for development:

```sh
docker-compose up --build
```

The backend will be available at `http://localhost:5000`.

### Running with Docker (Production)

The production Docker Compose file is located in the `src/` folder:

```sh
cd src
docker-compose up --build
```

The backend will be available at `http://localhost:5000`.

### API Documentation

Swagger UI is available at:  
`http://localhost:8080/swagger/`

### Environment Variables

Copy `.env.example` to `.env` and adjust as needed.

## Development

- Hot reload is enabled via [CompileDaemon](https://github.com/githubnemo/CompileDaemon).
- Logs are written to `src/app.log`.

## License

This project is licensed under the [MIT License](LICENSE).

---

Made with ❤️ by Gusde Widnyana

