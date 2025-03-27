# Chatbot Backend (Go + Gin + Groq)

## Features

- ✅ Gin-based API with stream response
- ✅ Middleware: Logging, Auth, Rate Limiting
- ✅ .env config via godotenv
- ✅ Dockerfile + Docker Compose
- ✅ CI/CD: GitHub Actions
- ✅ Unit tests + metrics + profiling

## Development

```bash
git clone https://github.com/jt3angga/test-chatbot-be.git
cd chatbot-backend
cp .env.example .env
docker-compose up --build
```

## Available Endpoints

- `POST /chat` - stream response
- `GET /healthz` - health check
- `GET /metrics` - Prometheus metrics
- `GET /debug/pprof` - Go profiler (pprof)

## License

-
