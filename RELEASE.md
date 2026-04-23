# Smart Vehicle System - Release Notes

## Version 1.0.0 - 2026-04-20

### Overview
Smart Vehicle System is an intelligent vehicle scheduling and management platform with real-time tracking, route optimization, and traffic simulation capabilities.

### Features
- **Real-time Vehicle Tracking**: WebSocket-based live vehicle position updates
- **Route Optimization**: Advanced A* algorithm for optimal path planning
- **Traffic Simulation**: Python-based traffic flow simulation with safety controls
- **Map Integration**: AMap (Gaode Maps) integration for real-world mapping
- **RESTful API**: Complete API for vehicle management and system control
- **Docker Deployment**: Production-ready containerized deployment
- **Monitoring & Logging**: Comprehensive logging and health monitoring

### Architecture
- **Backend**: Go (Gin framework)
- **Frontend**: HTML5/JavaScript with real-time WebSocket
- **Database**: In-memory with Redis support
- **Maps**: AMap JavaScript API
- **Simulation**: Python with NumPy/Matplotlib
- **Deployment**: Docker + Docker Compose + Nginx

### Quick Start

#### Prerequisites
- Docker & Docker Compose
- AMap API Key (for map functionality)
- Git

#### Installation
```bash
# Clone repository
git clone <repository-url>
cd smart-outgoing-demo

# Set up environment
cp deploy/.env.prod .env
# Edit .env with your AMap API keys

# Deploy
./deploy/deploy.sh prod
```

#### Configuration
Edit `.env` file with your settings:
- `AMAP_JS_KEY`: Your AMap JavaScript API key
- `AMAP_SECURITY_CODE`: Your AMap security code
- `REDIS_URL`: Redis connection string
- `LOG_LEVEL`: Logging level (info, debug, warn, error)

### API Endpoints

#### Vehicle Management
- `GET /api/vehicles` - List all vehicles
- `POST /api/vehicles` - Create new vehicle
- `PUT /api/vehicles/:id/destination` - Set vehicle destination
- `DELETE /api/vehicles/:id` - Delete vehicle
- `POST /api/vehicles/import` - Import vehicles from CSV

#### Algorithm & Simulation
- `POST /api/algorithm/plan` - Plan route for single vehicle
- `POST /api/algorithm/schedule` - Batch vehicle scheduling
- `POST /api/simulation/run` - Run traffic simulation
- `GET /api/simulation/status` - Get simulation status
- `GET /api/simulation/results` - Get simulation results

#### System Management
- `GET /api/config` - Get system configuration
- `PUT /api/config` - Update system configuration
- `GET /api/cache/stats` - Get cache statistics
- `POST /api/cache/cleanup` - Clear cache

#### WebSocket
- `GET /ws` - Real-time vehicle updates

### Deployment

#### Development
```bash
./deploy/deploy.sh dev
```

#### Production
```bash
./deploy/deploy.sh prod
```

#### Custom Registry
```bash
REGISTRY=your-registry.com ./deploy/deploy.sh prod
```

### Monitoring

#### Health Checks
- Application: `http://localhost:8080/api/config`
- Nginx: `http://localhost/health`
- Redis: `docker-compose exec redis redis-cli ping`

#### Logs
```bash
# View all logs
docker-compose logs -f

# View specific service logs
docker-compose logs -f smart-vehicle-server
docker-compose logs -f nginx
docker-compose logs -f redis
```

#### Metrics
- Docker stats: `docker stats`
- Container status: `docker-compose ps`

### Security

#### Production Security
- Nginx reverse proxy with rate limiting
- Security headers (XSS, CSRF protection)
- TLS/SSL support (configure in nginx.conf)
- Non-root container execution
- Health monitoring and automatic restart

#### API Security
- Input validation and sanitization
- Error handling without information leakage
- Rate limiting on API endpoints
- WebSocket connection limits

### Performance

#### Optimization Features
- In-memory caching with Redis backend
- Gzip compression for static assets
- Connection pooling and keep-alive
- Efficient WebSocket message handling
- Background cache cleanup

#### Scaling
- Horizontal scaling with Docker Compose
- Load balancing via Nginx
- Redis for shared state
- Stateless application design

### Troubleshooting

#### Common Issues

**Map not loading (black screen)**
- Check AMap API keys in `.env`
- Verify network connectivity to AMap services
- Check browser console for API errors

**WebSocket connection issues**
- Ensure port 8080 is accessible
- Check firewall settings
- Verify Nginx WebSocket configuration

**High memory usage**
- Monitor vehicle count and cache size
- Adjust cache cleanup intervals
- Scale with additional instances

#### Debug Mode
```bash
# Enable debug logging
export LOG_LEVEL=debug
./deploy/deploy.sh dev
```

### Development

#### Local Development
```bash
# Install dependencies
go mod download

# Run locally
go run cmd/server/main.go

# Run tests
go test ./...

# Build
go build -o main cmd/server/main.go
```

#### Code Structure
```
cmd/server/          # Application entry point
internal/
  handlers/          # HTTP handlers
  simulation/        # Traffic simulation bridge
  websocket/         # WebSocket management
  store/            # In-memory data store
  algorithm/        # Route planning algorithms
  config/           # Configuration management
deploy/             # Deployment configurations
public/             # Static frontend files
tools/              # Utility tools
```

### Contributing

#### Development Workflow
1. Fork repository
2. Create feature branch
3. Make changes with tests
4. Run `go test ./...`
5. Submit pull request

#### Code Standards
- Go fmt and go vet compliance
- Unit tests for new features
- Documentation for public APIs
- Security review for sensitive changes

### License

This project is licensed under the MIT License - see the LICENSE file for details.

### Support

For support and questions:
- Create an issue in the repository
- Check the troubleshooting section
- Review the API documentation

### Changelog

#### v1.0.0 (2026-04-20)
- Initial release
- Core vehicle tracking and management
- AMap integration
- Traffic simulation system
- Docker deployment
- RESTful API
- WebSocket real-time updates
- Production-ready configuration

---

**Note**: This release requires AMap API keys for map functionality. Please configure them in the `.env` file before deployment.
