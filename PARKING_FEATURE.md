# AI Car Parking Feature Documentation

## Overview

The AI Car Parking (AI car parking) feature is an intelligent parking solution that helps users find, reserve, and navigate to optimal parking spots. This feature addresses the common problem of parking difficulty in urban areas by providing real-time recommendations and seamless booking capabilities.

## Core Problem Solved

**"AI car parking" - One-click parking solution**
- **Issue**: Finding parking spots in crowded areas is time-consuming and stressful
- **Solution**: Intelligent recommendations based on location, price, availability, and user preferences
- **Benefit**: "Get in the car and drive straight into the parking lot" - seamless parking experience

## Architecture

### Domain-Driven Design Structure

```
internal/domain/
|-- entities/
|   |-- parking.go              # Core parking entities
|-- services/
|   |-- parking_recommendation_service.go  # Business logic
|-- repositories/
|   |-- parking_repositories.go  # Data access interfaces
```

### Key Entities

1. **ParkingLot** - Parking facility information
2. **ParkingSpace** - Individual parking spots
3. **ParkingReservation** - Time-based bookings
4. **ParkingSession** - Active parking sessions
5. **ParkingRecommendation** - AI-powered recommendations
6. **ParkingRoute** - Navigation instructions

### API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/parking/find` | Find best parking spots |
| POST | `/api/parking/reserve` | Reserve parking space |
| POST | `/api/parking/session/start` | Start parking session |
| GET | `/api/parking/lots` | List all parking lots |
| GET | `/api/parking/lots/:id` | Get parking lot details |
| GET | `/api/parking/lots/:id/spaces` | Get parking spaces |

## Features

### 1. Smart Search & Recommendations
- **Location-based**: GPS or manual address input
- **Intelligent Scoring**: Distance, price, availability, features
- **Personalization**: User preferences and vehicle type
- **Real-time**: Current availability and pricing

### 2. Interactive Map Interface
- **AMap Integration**: High-quality Chinese maps
- **Visual Markers**: Color-coded availability indicators
- **Real-time Updates**: Live parking status
- **Route Visualization**: Turn-by-turn navigation

### 3. Reservation System
- **Time-based Booking**: Specify start/end times
- **Space Selection**: Choose specific parking spots
- **Confirmation**: Instant booking confirmation
- **Management**: View and manage reservations

### 4. Session Management
- **One-click Start**: Begin parking session
- **Active Tracking**: Monitor current parking
- **Cost Calculation**: Real-time cost tracking
- **Session History**: Past parking records

### 5. Underground Automated Valet Parking (AVP)
- **Drop-off & Summon**: Driver exits at a designated zone, vehicle parks autonomously and can be summoned back
- **Vehicle-Infrastructure Collaboration**: Combines on-vehicle low-speed autonomy with parking-lot-side scheduling and sensing
- **Cross-level Routing**: Supports ramp traversal, lane-level path planning, and dynamic rerouting
- **Safety-First Design**: Conservative yielding logic for pedestrians, emergency-stop capability, and minimum risk stop strategy
- **Operational Readiness**: Includes dispatch policies, telemetry, replay, alerting, and staged rollout support

#### AVP Scope (Current Project Integration)
- **Phase 1 (Documentation + API contract)**: Define AVP workflow and interfaces for simulation
- **Phase 2 (Demo orchestration)**: Add backend orchestration endpoints and mock state transitions
- **Phase 3 (Production integration)**: Connect to real parking infrastructure, edge compute, and remote operations

#### AVP Functional Flow
1. User selects **Auto Park** in app/vehicle UI at drop-off zone
2. System validates vehicle status and available AVP area
3. Dispatch service allocates target slot and sends route task
4. Vehicle executes low-speed navigation with continuous obstacle checks
5. On successful parking, session status switches to `parked_by_avp`
6. User can trigger **Summon**, dispatching reverse route back to pickup zone

#### AVP Service Boundaries (DDD Mapping)
- **Entity extension candidates**
  - `ParkingSession`: add AVP mode, status transitions, safety events
  - `ParkingRoute`: add lane graph metadata and cross-floor waypoints
  - `ParkingSpace`: add AVP compatibility flags (width, pillar risk, charger proximity)
- **Domain service extension candidates**
  - New service: `avp_dispatch_service.go` for task assignment and conflict resolution
  - Existing recommendation service: include AVP-eligible score dimensions
- **Repository extension candidates**
  - AVP task repository for mission lifecycle and telemetry snapshots
  - Infrastructure capability repository for zone/ramp/closure status

#### AVP API Draft (for next implementation phase)
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/parking/avp/start` | Start AVP parking task from drop-off zone |
| POST | `/api/parking/avp/summon` | Summon vehicle to designated pickup zone |
| GET | `/api/parking/avp/tasks/:id` | Query AVP task state and latest progress |
| POST | `/api/parking/avp/tasks/:id/cancel` | Cancel AVP task and trigger safe stop |

#### AVP Key Metrics
- AVP parking success rate
- Average drop-off to parked duration
- Summon wait time
- Human intervention rate
- Safety incident count (target: zero collision)

## User Interface

### Web Application (`/parking`)

**Search Section:**
- Current location (GPS or manual)
- Price and distance filters
- Preference checkboxes (covered, EV charging, 24/7, security)

**Map Display:**
- Interactive AMap integration
- Parking lot markers with availability
- User location indicator
- Route visualization

**Recommendations:**
- Card-based layout with scores
- Detailed parking information
- One-click reservation buttons
- Navigation integration

## Implementation Details

### Scoring Algorithm

The recommendation system uses a multi-factor scoring algorithm:

```go
Score = (DistanceWeight * DistanceScore) +
        (PriceWeight * PriceScore) +
        (AvailabilityWeight * AvailabilityScore) +
        (FeatureWeight * FeatureMatchScore)
```

**Factors:**
- **Distance**: Closer is better (inverse relationship)
- **Price**: Lower is better (inverse relationship)
- **Availability**: More spaces is better
- **Features**: Match with user preferences

### Mock Data System

Currently uses mock data for demonstration:

```go
// Mock Parking Lots
- CBD Central Parking: 200 spaces, ¥15/hr, covered + EV
- Shopping Mall Parking: 150 spaces, ¥10/hr, covered + security  
- Airport Parking: 300 spaces, ¥8/hr, 24/7 + shuttle
```

### Future Enhancements

1. **Real Data Integration**
   - Connect to actual parking APIs
   - Live availability updates
   - Dynamic pricing

2. **Payment Integration**
   - WeChat Pay / Alipay
   - Credit card processing
   - Automatic billing

3. **Advanced Features**
   - Predictive availability
   - Dynamic pricing
   - Loyalty programs
   - Fleet management

## Technical Implementation

### Frontend Technologies
- **HTML5/CSS3**: Modern responsive design
- **Tailwind CSS**: Utility-first styling
- **JavaScript ES6+**: Modern web standards
- **AMap API**: Chinese mapping service

### Backend Technologies
- **Go**: High-performance server
- **Gin Framework**: HTTP routing and middleware
- **Domain-Driven Design**: Clean architecture
- **JSON APIs**: RESTful endpoints

### Integration Points
- **AMap JavaScript API**: Map rendering and geocoding
- **GPS/Geolocation**: User location detection
- **WebSocket**: Real-time updates (future)
- **Payment Gateways**: Financial processing (future)

## Deployment

### Local Development
```bash
# Build and run
go build -o main cmd/server/main.go
./main

# Access parking feature
open http://localhost:8080/parking
```

### Production Deployment
```bash
# Using Docker Compose
docker-compose -f deploy/docker-compose.prod.yml up -d

# Or direct deployment
./deploy/deploy.sh prod
```

### Environment Variables
```bash
# AMap Configuration
AMAP_JS_KEY=45109d104b3c8d03a2c84175a7749241
AMAP_SECURITY_CODE=c552677838e5f5e71de92ce532c936bc

# Server Configuration
PORT=8080
ENVIRONMENT=production
```

## API Usage Examples

### Find Parking Spots
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

### Reserve Parking Space
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

### Start Parking Session
```bash
curl -X POST http://localhost:8080/api/parking/session/start \
  -H "Content-Type: application/json" \
  -H "x-user-id: demo_user" \
  -d '{
    "parking_lot_id": "lot_1",
    "space_id": "space_1"
  }'
```

## Testing

### API Testing
```bash
# Run all API tests
./scripts/test_parking_apis.sh

# Run specific endpoint tests
curl -s http://localhost:8080/api/parking/lots | jq .
```

### UI Testing
- Access `http://localhost:8080/parking`
- Test GPS location detection
- Verify map rendering
- Test reservation flow

## Performance Considerations

### Caching Strategy
- Parking lot data: Cache for 5 minutes
- Availability data: Cache for 30 seconds
- User preferences: Cache for 1 hour

### Scalability
- Horizontal scaling with load balancers
- Database read replicas for parking data
- CDN for static assets and map tiles

### Monitoring
- API response times
- Map loading performance
- User interaction metrics
- Error rates and patterns

## Security Considerations

### Data Protection
- User location privacy
- Payment information security
- Personal data encryption

### API Security
- Rate limiting
- Authentication tokens
- Input validation
- CORS configuration

## Troubleshooting

### Common Issues

1. **Map Not Loading**
   - Check AMap API keys
   - Verify network connectivity
   - Check browser console errors

2. **GPS Not Working**
   - Enable location services
   - Check HTTPS requirement
   - Verify browser permissions

3. **API Errors**
   - Check server logs
   - Verify JSON format
   - Check authentication headers

### Debug Mode
```bash
# Enable debug logging
export LOG_LEVEL=debug
./main

# Check API responses
curl -v http://localhost:8080/api/parking/lots
```

## Future Roadmap

### Phase 1: Production Ready
- [ ] Real parking data integration
- [ ] Payment system integration
- [ ] User authentication
- [ ] Mobile app development

### Phase 2: Advanced Features
- [ ] Predictive analytics
- [ ] Dynamic pricing
- [ ] Fleet management
- [ ] Corporate partnerships

### Phase 3: Ecosystem Integration
- [ ] Smart city integration
- [ ] IoT sensor connectivity
- [ ] AI-powered optimization
- [ ] Multi-city expansion

## Conclusion

The AI Car Parking feature successfully addresses urban parking challenges through intelligent recommendations, seamless booking, and intuitive navigation. The modular architecture allows for easy expansion and integration with real-world parking systems, providing a solid foundation for a comprehensive parking solution platform.
