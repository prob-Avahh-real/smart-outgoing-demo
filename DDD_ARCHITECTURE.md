# DDD Auto-Scaling Architecture Documentation

## Overview

This document describes the Domain-Driven Design (DDD) auto-scaling architecture implemented for the Smart Vehicle System. The architecture provides automatic storage scaling between in-memory and Redis storage based on system load thresholds.

## Architecture Components

### Domain Layer (`/internal/domain/`)

#### Aggregates
- **Vehicle**: Core vehicle entity with coordinates, status, and destination management
- **Route**: Vehicle route information with waypoints and status tracking
- **Metrics**: System performance metrics including memory, CPU, and connection counts
- **ScalingDecision**: Auto-scaling decisions with strategy changes and reasoning

#### Value Objects
- **Coordinates**: Geographic location with longitude, latitude, and altitude
- **MemoryMetrics**: Memory usage statistics
- **ScalingThreshold**: Configuration thresholds for scaling decisions
- **StorageStrategy**: Enum for memory, Redis, and hybrid storage strategies

#### Domain Events
- **StorageStrategyChangedEvent**: Fired when storage strategy changes
- **ThresholdBreachedEvent**: Fired when system thresholds are exceeded
- **VehicleCreatedEvent**: Fired when a new vehicle is created
- **RouteCompletedEvent**: Fired when a vehicle completes a route

#### Domain Services
- **ScalingDecisionService**: Evaluates metrics and makes scaling decisions
- **VehicleManagementService**: Manages vehicle lifecycle and route assignments
- **EventBus**: In-memory event publishing and subscription

### Application Layer (`/internal/application/services/`)

#### Services
- **DDDScalingOrchestrator**: Coordinates zero-downtime scaling migrations
- **DDDIntegrationService**: Provides external API for the DDD system
- **MigrationService**: Handles data migration between storage strategies
- **UnifiedRepositoryFactory**: Creates appropriate repositories based on strategy

#### Event Handlers
- **StorageStrategyChangedHandler**: Processes storage strategy changes
- **ThresholdBreachedHandler**: Handles threshold breach notifications
- **VehicleCreatedHandler**: Processes new vehicle creation
- **RouteCompletedHandler**: Handles route completion events

### Infrastructure Layer (`/internal/infrastructure/`)

#### Memory Storage (`/memory/`)
- **DDDMemoryVehicleRepository**: In-memory vehicle storage
- **DDDMemoryRouteRepository**: In-memory route storage
- **DDDMemoryMetricsRepository**: In-memory metrics storage
- **DDDMemoryScalingDecisionRepository**: In-memory scaling decision storage

#### Redis Storage (`/redis/`)
- **DDDRedisVehicleRepository**: Redis-based vehicle storage
- **DDDRedisRouteRepository**: Redis-based route storage
- **DDDRedisMetricsRepository**: Redis-based metrics storage
- **DDDRedisScalingDecisionRepository**: Redis-based scaling decision storage

## Scaling Strategies

### Memory Strategy
- Fast access for low-load scenarios
- Suitable for development and small deployments
- Limited by available RAM

### Redis Strategy
- Persistent storage for medium-load scenarios
- Supports data persistence across restarts
- Network latency considerations

### Hybrid Strategy
- Combines memory and Redis for optimal performance
- Uses Redis as fallback for memory constraints
- Provides best of both worlds

## Auto-Scaling Process

1. **Metrics Collection**: System metrics are collected periodically
2. **Threshold Evaluation**: ScalingDecisionService evaluates against configured thresholds
3. **Decision Making**: Scaling decisions are made based on multiple factors
4. **Zero-Downtime Migration**: Data is migrated without service interruption
5. **Strategy Switch**: Storage strategy is updated atomically
6. **Event Publishing**: Domain events are published for monitoring

## Configuration

### Default Thresholds
```go
ScalingThreshold{
    MemoryUsagePercent:   80.0,
    AgentCount:          100,
    ConnectionCount:     1000,
    CPUUsagePercent:     75.0,
}
```

### Scaling Factors
- **Scale-up**: Memory → Redis → Hybrid
- **Scale-down**: Hybrid → Redis → Memory

## API Usage

### Integration Service
```go
// Create integration service
integrationService := services.NewDDDIntegrationService()

// Start the system
ctx := context.Background()
err := integrationService.Start(ctx)

// Create a vehicle
vehicle, err := integrationService.CreateVehicle(ctx, "v1", "Car", "transport", 0.0, 0.0, 0.0)

// Assign a route
waypoints := []domain.Coordinates{{0,0,0}, {1,1,0}, {2,2,0}}
route, err := integrationService.AssignRoute(ctx, "v1", waypoints)

// Force scaling
err = integrationService.ForceScaling(ctx, domain.StorageStrategyRedis)
```

## Testing

### Unit Tests
- Domain aggregate behavior tests
- Repository implementation tests
- Service logic tests

### Integration Tests
- End-to-end scaling scenarios
- Event handling validation
- Performance benchmarks

### Test Coverage
- Domain layer: >90%
- Application layer: >85%
- Infrastructure layer: >80%

## Performance Characteristics

### Memory Storage
- **Read Latency**: <1ms
- **Write Latency**: <1ms
- **Capacity**: Limited by RAM
- **Persistence**: No

### Redis Storage
- **Read Latency**: 1-5ms
- **Write Latency**: 1-5ms
- **Capacity**: Limited by Redis memory
- **Persistence**: Configurable

### Migration Time
- **Small datasets** (<1000 entities): <100ms
- **Medium datasets** (1000-10000 entities): <1s
- **Large datasets** (>10000 entities): <10s

## Monitoring and Observability

### Metrics
- Current storage strategy
- Migration progress
- Threshold breaches
- Event processing rates

### Logging
- Structured logging with correlation IDs
- Event-driven audit trails
- Error tracking and alerting

### Health Checks
- Repository connectivity
- Event bus status
- Migration service availability

## Deployment Considerations

### Production Configuration
- Redis cluster setup for high availability
- Memory limits and monitoring
- Backup and recovery procedures

### Scaling Recommendations
- Start with Memory strategy for new deployments
- Monitor metrics and set appropriate thresholds
- Use Hybrid strategy for production workloads
- Implement automated alerts for threshold breaches

## Future Enhancements

### Planned Features
- Multi-region Redis support
- Advanced caching strategies
- Machine learning-based scaling predictions
- Grafana dashboard integration

### Extensibility
- Plugin architecture for custom repositories
- Configurable scaling algorithms
- Additional storage backends (PostgreSQL, MongoDB)
- Event sourcing with replay capabilities

## Troubleshooting

### Common Issues
1. **Migration Failures**: Check Redis connectivity and memory limits
2. **Event Loss**: Verify event bus configuration and handler registration
3. **Performance Degradation**: Monitor repository performance and optimize queries
4. **Inconsistent State**: Run consistency checks and data reconciliation

### Debugging Tools
- Structured logging with trace IDs
- Metrics dashboards
- Health check endpoints
- Event replay capabilities

## Conclusion

The DDD auto-scaling architecture provides a robust, scalable solution for dynamic storage management. The implementation follows domain-driven design principles, ensuring clean separation of concerns and maintainable code structure. The system can handle varying load conditions while maintaining data consistency and service availability.
