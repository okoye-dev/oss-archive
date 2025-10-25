.PHONY: dev dev-start dev-stop dev-backend dev-backend-supabase dev-logs minio help

# Default target
help:
	@echo "🛠️  OSS Archive Development Commands"
	@echo ""
	@echo "Setup (first time):"
	@echo "  cp .env.local.example .env.local"
	@echo "  cp .env.prod.example .env.prod"
	@echo "  # Edit .env.prod with your Supabase credentials"
	@echo ""
	@echo "Development (Fast iteration):"
	@echo "  make dev-start           Start MinIO + Frontend (hot reload)"
	@echo "  make dev-backend         Start Go backend with MinIO (.env.local)"
	@echo "  make dev-backend-sb      Start Go backend with Supabase (.env.prod)"
	@echo "  make dev-stop            Stop development services"
	@echo "  make dev-logs            View development logs"
	@echo ""
	@echo "Production:"
	@echo "  make prod                Build and run production container"
	@echo "  make prod-logs           View production logs"
	@echo "  make prod-stop           Stop production container"
	@echo ""
	@echo "Utilities:"
	@echo "  make minio               Open MinIO console"
	@echo "  make clean               Clean up Docker resources"

# Development commands
dev:
	@echo "🚀 Starting development environment..."
	docker compose -f docker-compose.dev.yml up -d
	@echo "✅ Services started:"
	@echo "   🗄️  MinIO Console: http://localhost:9001 (minioadmin/minioadmin)"
	@echo "   🌐 Frontend: http://localhost:3000"
	@echo ""
	@echo "💡 Next: run 'make dev-backend' in another terminal"

dev-backend:
	@echo "🔧 Starting Go backend with MinIO..."
	@export $$(cat .env.local | grep -v '^#' | xargs) && go run cmd/main.go

dev-backend-sb:
	@echo "🔧 Starting Go backend with Supabase..."
	@export $$(cat .env.prod | grep -v '^#' | xargs) && go run cmd/main.go

dev-stop:
	@echo "🛑 Stopping development environment..."
	docker compose -f docker-compose.dev.yml down

dev-logs:
	docker compose -f docker-compose.dev.yml logs -f

# Production commands
prod:
	@echo "🏗️  Building and starting production container..."
	@echo "⚠️  Note: Set environment variables manually or use Railway"
	@echo "💡 For local testing with Supabase: export variables from .env.prod first"
	docker compose up --build -d

prod-logs:
	docker compose logs -f app

prod-stop:
	@echo "🛑 Stopping production container..."
	docker compose down

# Utilities
minio:
	@echo "🗄️  Opening MinIO Console..."
	@open http://localhost:9001 2>/dev/null || echo "Open http://localhost:9001 in your browser"

clean:
	@echo "🧹 Cleaning up Docker resources..."
	docker compose -f docker-compose.dev.yml down -v
	docker compose down -v
	docker system prune -f
