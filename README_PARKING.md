# AI Car Parking - Quick Start Guide

## Quick Start

### 1. Launch the Application
```bash
# Build and run
go build -o main cmd/server/main.go
./main
```

### 2. Access the Parking Feature
Open your browser and navigate to:
```
http://localhost:8080/parking
```

### 3. Test the APIs

#### Find Parking Spots
```bash
curl -X POST http://localhost:8080/api/parking/find \
  -H "Content-Type: application/json" \
  -d '{
    "latitude": 22.6913,
    "longitude": 114.0448,
    "max_price": 20,
    "max_distance": 5,
    "limit": 5
  }'
```

#### Reserve a Space
```bash
curl -X POST http://localhost:8080/api/parking/reserve \
  -H "Content-Type: application/json" \
  -H "x-user-id: demo_user" \
  -d '{
    "parking_lot_id": "lot_1",
    "space_id": "space_1",
    "start_time": "2026-04-21T12:00:00Z",
    "end_time": "2026-04-21T14:00:00Z"
  }'
```

#### Start Parking Session
```bash
curl -X POST http://localhost:8080/api/parking/session/start \
  -H "Content-Type: application/json" \
  -H "x-user-id: demo_user" \
  -d '{
    "parking_lot_id": "lot_1",
    "space_id": "space_1"
  }'
```

## Feature Overview

### What It Does
- **Smart Search**: Find parking based on location, price, and preferences
- **Real-time Recommendations**: AI-powered scoring system
- **Easy Booking**: One-click space reservation
- **Navigation**: Turn-by-turn directions to parking spots
- **Session Management**: Track active parking sessions

### Key Benefits
- **Time Saving**: Find optimal parking quickly
- **Cost Effective**: Compare prices and find best deals
- **Convenient**: Reserve spots in advance
- **User Friendly**: Simple, intuitive interface

## Architecture

### Core Components
```
Frontend (React/Vue ready)
    |
    v
API Layer (Gon/REST)
    |
    v
Business Logic (Go Services)
    |
    v
Data Layer (Repository Pattern)
```

### Main Files
- `public/html/parking.html` - Web interface
- `internal/handlers/parking_handlers.go` - API endpoints
- `internal/domain/services/parking_recommendation_service.go` - Business logic
- `internal/domain/entities/parking.go` - Data models

## Configuration

### Environment Variables
```bash
# AMap Configuration (required for maps)
AMAP_JS_KEY=45109d104b3c8d03a2c84175a7749241
AMAP_SECURITY_CODE=c552677838e5f5e71de92ce532c936bc

# Server Configuration
PORT=8080
ENVIRONMENT=development
```

### Setup AMap Keys
1. Register at [AMap Console](https://console.amap.com/)
2. Create a new application
3. Get JavaScript API key and security code
4. Update environment variables

## Testing

### Manual Testing
1. Open `http://localhost:8080/parking`
2. Click "GPS" to get current location
3. Set price and distance filters
4. Click "Find Parking Spots"
5. View recommendations on map
6. Click "Reserve" on any spot
7. Fill reservation form and confirm

### API Testing
```bash
# Test all endpoints
./scripts/test_all_parking_apis.sh

# Test specific functionality
curl -s http://localhost:8080/api/parking/lots | jq '.parking_lots | length'
```

## Deployment

### Development
```bash
go run cmd/server/main.go
```

### Production
```bash
# Build
go build -o parking-server cmd/server/main.go

# Run with production config
ENVIRONMENT=production ./parking-server
```

### Docker
```bash
# Build image
docker build -t ai-parking .

# Run container
docker run -p 8080:8080 ai-parking
```

## Troubleshooting

### Common Issues

**Map Not Loading**
- Check AMap API keys in environment
- Verify network connectivity
- Check browser console for errors

**API Not Responding**
- Check if server is running: `ps aux | grep main`
- Verify port 8080 is available
- Check server logs for errors

**GPS Not Working**
- Enable location services in browser
- Use HTTPS in production (required for GPS)
- Check browser permissions

### Debug Mode
```bash
# Enable debug logging
export LOG_LEVEL=debug
./main

# Test with verbose curl
curl -v http://localhost:8080/api/parking/lots
```

## Next Steps

### For Production Use
1. **Real Data Integration**: Connect to actual parking APIs
2. **Payment System**: Add payment processing
3. **User Authentication**: Implement user accounts
4. **Mobile App**: Develop native mobile applications

### For Development
1. **Unit Tests**: Add comprehensive test coverage
2. **Database**: Replace mock data with real database
3. **Caching**: Add Redis for performance
4. **Monitoring**: Add metrics and logging

## Support

### Documentation
- Full documentation: `PARKING_FEATURE.md`
- API reference: Check endpoint responses
- Architecture guide: See domain design

### Getting Help
- Check server logs for errors
- Verify environment configuration
- Test APIs individually
- Check browser console for frontend issues

## Quick Demo

Try this 5-minute demo:

1. **Start**: `./main` then open `http://localhost:8080/parking`
2. **Location**: Click GPS button or enter "22.6913,114.0448"
3. **Search**: Set max price ¥20, distance 5km, click "Find Parking"
4. **Explore**: View 3 recommended spots with scores
5. **Book**: Click "Reserve" on first spot, set time, confirm
6. **Navigate**: See directions and session info

That's it! You've experienced the complete AI Car Parking workflow.
