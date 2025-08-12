# URL Shortener Service

A high-performance URL shortening service built with Go and Gin, featuring rate limiting, analytics, and Redis persistence.

## 🚀 Quick Start

### Prerequisites
- Go 1.21+
- Docker & Docker Compose
- Redis

### Local Development
```bash
# Install dependencies
make deps

# Build and run
make build
make run
```

### Docker Deployment
```bash
# Build and run with Docker
make docker-build
make docker-run-detached

# View logs
make docker-logs

# Stop services
make docker-stop
```

## 📡 API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/:url` | Redirect to original URL |
| `POST` | `/api/v1` | Create shortened URL |
| `GET` | `/api/v1/analytics` | Get total redirect count |
| `GET` | `/api/v1/analytics/:url` | Get URL-specific analytics |

### Example Usage
```bash
# Shorten a URL
curl -X POST http://localhost:3000/api/v1 \
  -H "Content-Type: application/json" \
  -d '{"url": "https://example.com/very-long-url"}'

# Access shortened URL
curl http://localhost:3000/abc123

# Get analytics
curl http://localhost:3000/api/v1/analytics

# Get analytics for a specific url
curl http://localhost:3000/api/v1/analytics/abc123
```

## 🔧 Configuration

Environment variables (optional - defaults provided):
- `DB_ADDR`: Redis address (default: localhost:6379)
- `DB_PASS`: Redis password (default: empty)
- `DOMAIN`: Application domain (default: http://localhost:3000)
- `API_QUOTA`: Rate limit quota (default: 20)
- `RATE_LIMIT_MINUTES`: Rate limit window (default: 30)

## 🏗️ Architecture

- **Web Framework**: Gin (high-performance HTTP framework)
- **Database**: Redis (in-memory data store)
- **Rate Limiting**: Per-IP request limiting
- **Analytics**: URL access tracking and counters

## 📊 Features

- ✅ URL shortening with custom short codes
- ✅ Rate limiting (20 requests per 30 minutes)
- ✅ Analytics tracking
- ✅ Redis persistence
- ✅ Docker support
- ✅ Clean architecture with service layer
- ✅ Comprehensive error handling
