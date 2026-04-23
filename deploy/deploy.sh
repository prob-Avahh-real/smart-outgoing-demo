#!/bin/bash

# Smart Vehicle System Deployment Script
# Usage: ./deploy.sh [dev|prod|staging]

set -e

# Configuration
ENVIRONMENT=${1:-prod}
PROJECT_NAME="smart-vehicle"
VERSION=$(date +%Y%m%d-%H%M%S)
REGISTRY=${REGISTRY:-"your-registry.com"}
IMAGE_TAG="${REGISTRY}/${PROJECT_NAME}:${VERSION}"

echo "=== Smart Vehicle System Deployment ==="
echo "Environment: ${ENVIRONMENT}"
echo "Version: ${VERSION}"
echo "Image Tag: ${IMAGE_TAG}"
echo "======================================="

# Function to check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Check prerequisites
check_prerequisites() {
    echo "Checking prerequisites..."
    
    if ! command_exists docker; then
        echo "Error: Docker is not installed"
        exit 1
    fi
    
    if ! command_exists docker-compose; then
        echo "Error: Docker Compose is not installed"
        exit 1
    fi
    
    echo "Prerequisites check passed"
}

# Build Docker image
build_image() {
    echo "Building Docker image..."
    
    if [ "$ENVIRONMENT" = "prod" ]; then
        docker build -f deploy/Dockerfile.prod -t "${IMAGE_TAG}" .
    else
        docker build -f Dockerfile -t "${IMAGE_TAG}" .
    fi
    
    echo "Docker image built successfully: ${IMAGE_TAG}"
}

# Run tests
run_tests() {
    echo "Running tests..."
    
    # Run unit tests
    docker run --rm "${IMAGE_TAG}" go test -v ./...
    
    # Run integration tests
    docker run --rm -p 8080:8080 "${IMAGE_TAG}" &
    SERVER_PID=$!
    
    # Wait for server to start
    sleep 5
    
    # Test API endpoints
    curl -f http://localhost:8080/api/config || {
        echo "API health check failed"
        kill $SERVER_PID
        exit 1
    }
    
    kill $SERVER_PID
    echo "Tests passed"
}

# Deploy to environment
deploy() {
    echo "Deploying to ${ENVIRONMENT}..."
    
    # Create necessary directories
    mkdir -p data logs deploy/ssl
    
    # Set up environment file
    if [ ! -f ".env" ]; then
        if [ -f "deploy/.env.${ENVIRONMENT}" ]; then
            cp "deploy/.env.${ENVIRONMENT}" .env
            echo "Environment file copied from deploy/.env.${ENVIRONMENT}"
        else
            echo "Warning: No environment file found at deploy/.env.${ENVIRONMENT}"
            echo "Please create .env file with proper configuration"
        fi
    fi
    
    # Deploy with Docker Compose
    if [ "$ENVIRONMENT" = "prod" ]; then
        docker-compose -f deploy/docker-compose.prod.yml up -d
    else
        docker-compose up -d
    fi
    
    echo "Deployment completed"
}

# Health check
health_check() {
    echo "Performing health check..."
    
    # Wait for services to start
    sleep 10
    
    # Check main service
    if curl -f http://localhost:8080/api/config > /dev/null 2>&1; then
        echo "Main service is healthy"
    else
        echo "Main service health check failed"
        docker-compose logs
        exit 1
    fi
    
    # Check Redis
    if docker-compose exec redis redis-cli ping > /dev/null 2>&1; then
        echo "Redis is healthy"
    else
        echo "Redis health check failed"
        exit 1
    fi
    
    echo "All services are healthy"
}

# Show deployment info
show_info() {
    echo "=== Deployment Information ==="
    echo "Service URLs:"
    echo "  Main App: http://localhost:8080"
    echo "  API Docs: http://localhost:8080/api/config"
    echo "  WebSocket: ws://localhost:8080/ws"
    echo "  Health: http://localhost/health"
    echo ""
    echo "Useful Commands:"
    echo "  View logs: docker-compose logs -f"
    echo "  Stop services: docker-compose down"
    echo "  Restart: docker-compose restart"
    echo "  Scale: docker-compose up -d --scale smart-vehicle-server=3"
    echo ""
    echo "Monitoring:"
    echo "  Docker stats: docker stats"
    echo "  Container status: docker-compose ps"
    echo "=============================="
}

# Cleanup old images
cleanup() {
    echo "Cleaning up old images..."
    
    # Remove old images (keep last 5)
    docker images "${REGISTRY}/${PROJECT_NAME}" --format "table {{.Repository}}:{{.Tag}}" | tail -n +2 | tail -n +6 | xargs -r docker rmi
    
    echo "Cleanup completed"
}

# Main deployment flow
main() {
    check_prerequisites
    build_image
    
    if [ "$ENVIRONMENT" != "prod" ] || [ "${SKIP_TESTS}" != "true" ]; then
        run_tests
    fi
    
    deploy
    health_check
    show_info
    cleanup
    
    echo "Deployment completed successfully!"
    echo "Version: ${VERSION}"
    echo "Environment: ${ENVIRONMENT}"
}

# Handle script arguments
case "${1:-}" in
    "dev"|"prod"|"staging")
        main
        ;;
    "build-only")
        check_prerequisites
        build_image
        ;;
    "test-only")
        build_image
        run_tests
        ;;
    "cleanup")
        cleanup
        ;;
    "help"|"-h"|"--help")
        echo "Usage: $0 [dev|prod|staging|build-only|test-only|cleanup|help]"
        echo ""
        echo "Commands:"
        echo "  dev         Deploy to development environment"
        echo "  prod        Deploy to production environment"
        echo "  staging     Deploy to staging environment"
        echo "  build-only  Build Docker image only"
        echo "  test-only   Run tests only"
        echo "  cleanup     Clean up old images"
        echo "  help        Show this help message"
        echo ""
        echo "Environment Variables:"
        echo "  REGISTRY    Docker registry URL (default: your-registry.com)"
        echo "  SKIP_TESTS  Skip tests (default: false)"
        exit 0
        ;;
    *)
        echo "Unknown command: ${1:-}"
        echo "Use '$0 help' for usage information"
        exit 1
        ;;
esac
