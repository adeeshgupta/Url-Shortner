# URL Shortener Service

A high-performance URL shortening service built with Go and Gin, featuring rate limiting, analytics, and Redis persistence.

## 🚀 Quick Start

### Prerequisites
- Go 1.21+
- Redis

### Setup
```bash
# Install Redis (macOS)
brew install redis
brew services start redis

# Install dependencies
make deps

# Build and run
make build
make run
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

- ✅ Clean architecture with service layer
- ✅ Comprehensive error handling

## 🛠️ Development

```bash
# Format code
make fmt

# Clean build artifacts
make clean

## 🔍 Viewing Redis Data

```bash
# Connect to Redis CLI
redis-cli

# View URL mappings (Database 0)
SELECT 0
KEYS *
GET <short_code>

# View analytics (Database 1)
SELECT 1
KEYS *
GET counter
GET access:<short_code>
```
