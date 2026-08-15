#!/usr/bin/env bash
# ==============================================================================
# Batam MedHub — Production Deployment Script
# ==============================================================================
set -Eeuo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"

cd "$ROOT_DIR"

ENV_FILE=".env.production"
COMPOSE_FILE="docker-compose.prod.yml"

echo "======================================================"
echo "🚀 Batam MedHub — Production Deployment"
echo "======================================================"

# 1. Check prerequisites
if ! command -v docker &> /dev/null; then
    echo "❌ Error: Docker is not installed. Please install Docker first."
    exit 1
fi

if ! docker compose version &> /dev/null; then
    echo "❌ Error: Docker Compose is not available. Please install Docker Compose v2."
    exit 1
fi

# 2. Check environment file
if [ ! -f "$ENV_FILE" ]; then
    if [ -f ".env" ]; then
        ENV_FILE=".env"
        echo "ℹ️ Using .env file for environment variables."
    else
        echo "❌ Error: $ENV_FILE not found!"
        echo "👉 Silakan buat file $ENV_FILE terlebih dahulu:"
        echo "   cp .env.production.example .env.production"
        echo "   nano .env.production"
        exit 1
    fi
fi

# 3. Pull / Build & Start containers
echo "📦 Building and starting containers with $COMPOSE_FILE..."
docker compose --env-file "$ENV_FILE" -f "$COMPOSE_FILE" up -d --build --remove-orphans

echo "⏳ Waiting for services to become healthy..."
sleep 5

# 4. Check container status
docker compose --env-file "$ENV_FILE" -f "$COMPOSE_FILE" ps

# 5. Verify Core Backend Health Check
echo "🔍 Verifying Core Backend health..."
MAX_RETRIES=12
RETRY_COUNT=0
HEALTHY=false

while [ $RETRY_COUNT -lt $MAX_RETRIES ]; do
    if docker compose --env-file "$ENV_FILE" -f "$COMPOSE_FILE" exec -T backend wget -q -O - http://127.0.0.1:8080/healthz > /dev/null 2>&1; then
        HEALTHY=true
        break
    fi
    echo "   Service starting... retry $((RETRY_COUNT+1))/$MAX_RETRIES"
    sleep 3
    RETRY_COUNT=$((RETRY_COUNT+1))
done

if [ "$HEALTHY" = true ]; then
    echo "======================================================"
    echo "✅ Batam MedHub is successfully deployed and healthy!"
    echo "======================================================"
    echo "Useful commands:"
    echo "  - View logs:    docker compose --env-file $ENV_FILE -f $COMPOSE_FILE logs -f"
    echo "  - Backend logs: docker compose --env-file $ENV_FILE -f $COMPOSE_FILE logs -f backend"
    echo "  - Caddy logs:   docker compose --env-file $ENV_FILE -f $COMPOSE_FILE logs -f caddy"
    echo "  - Stop stack:   docker compose --env-file $ENV_FILE -f $COMPOSE_FILE down"
    echo "  - Reset DB:     docker compose --env-file $ENV_FILE -f $COMPOSE_FILE restart"
    echo "======================================================"
else
    echo "⚠️ Warning: Health check timed out. Please check container logs with:"
    echo "   docker compose --env-file $ENV_FILE -f $COMPOSE_FILE logs backend"
    exit 1
fi
