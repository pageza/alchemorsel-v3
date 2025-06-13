# Alchemorsel System Architecture

## High-Level Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                        Frontend Layer                           │
├─────────────────────────────────────────────────────────────────┤
│  Vue.js 3 + Composition API + TypeScript + Pinia + Vue Router  │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐ ┌──────────────┐│
│  │ Auth Views  │ │Recipe Views │ │Profile Views│ │ Search Views ││
│  └─────────────┘ └─────────────┘ └─────────────┘ └──────────────┘│
└─────────────────────────────────────────────────────────────────┘
                                   │
                              HTTPS/JSON API
                                   │
┌─────────────────────────────────────────────────────────────────┐
│                        API Gateway Layer                        │
├─────────────────────────────────────────────────────────────────┤
│  Gin HTTP Server + Middleware (CORS, Auth, Logging, Recovery)  │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐ ┌──────────────┐│
│  │ Auth Routes │ │Recipe Routes│ │Profile Route│ │ LLM Routes   ││
│  └─────────────┘ └─────────────┘ └─────────────┘ └──────────────┘│
└─────────────────────────────────────────────────────────────────┘
                                   │
                           Dependency Injection
                                   │
┌─────────────────────────────────────────────────────────────────┐
│                       Service Layer                            │
├─────────────────────────────────────────────────────────────────┤
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐ ┌──────────────┐│
│  │ Auth Service│ │Recipe Service│ │Profile Svc  │ │ LLM Service  ││
│  └─────────────┘ └─────────────┘ └─────────────┘ └──────────────┘│
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐ ┌──────────────┐│
│  │Vector Service│ │Cache Service│ │Email Service│ │ File Service ││
│  └─────────────┘ └─────────────┘ └─────────────┘ └──────────────┘│
└─────────────────────────────────────────────────────────────────┘
                                   │
                          Repository Interface
                                   │
┌─────────────────────────────────────────────────────────────────┐
│                      Repository Layer                          │
├─────────────────────────────────────────────────────────────────┤
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐ ┌──────────────┐│
│  │ User Repo   │ │Recipe Repo  │ │Profile Repo │ │ Cache Repo   ││
│  └─────────────┘ └─────────────┘ └─────────────┘ └──────────────┘│
└─────────────────────────────────────────────────────────────────┘
                                   │
                            Database Access
                                   │
┌─────────────────────────────────────────────────────────────────┐
│                      Data Storage Layer                        │
├─────────────────────────────────────────────────────────────────┤
│  ┌─────────────────────┐ ┌─────────────────────┐ ┌──────────────┐│
│  │   PostgreSQL        │ │       Redis         │ │    S3/Minio  ││
│  │  (with pgvector)    │ │   (Cache/Sessions)  │ │ (File Storage││
│  └─────────────────────┘ └─────────────────────┘ └──────────────┘│
└─────────────────────────────────────────────────────────────────┘
                                   │
                           External Services
                                   │
┌─────────────────────────────────────────────────────────────────┐
│                    External Services                           │
├─────────────────────────────────────────────────────────────────┤
│  ┌─────────────────────┐ ┌─────────────────────┐ ┌──────────────┐│
│  │   OpenAI API        │ │    DeepSeek API     │ │  SendGrid    ││
│  │   (Embeddings)      │ │  (Recipe Generation)│ │   (Email)    ││
│  └─────────────────────┘ └─────────────────────┘ └──────────────┘│
└─────────────────────────────────────────────────────────────────┘
```

## Key Architecture Principles

### 1. Clean Architecture
- **Domain Layer**: Core business entities and rules
- **Application Layer**: Use cases and business logic
- **Infrastructure Layer**: External concerns (DB, APIs, etc.)
- **Presentation Layer**: HTTP handlers and DTOs

### 2. Dependency Injection
- Service container for managing dependencies
- Interface-based design for testability
- Dependency inversion principle

### 3. CQRS Pattern (Light)
- Separate read and write operations where beneficial
- Optimized queries for search/listing
- Command/Query separation in services

### 4. Event-Driven Components
- Domain events for cross-cutting concerns
- Async processing for expensive operations
- Event sourcing for audit trails

## Technology Stack

### Backend
- **Language**: Go 1.21+
- **Web Framework**: Gin
- **ORM**: GORM with custom repository layer
- **Database**: PostgreSQL 15+ with pgvector extension
- **Cache**: Redis 7+
- **Authentication**: JWT with refresh tokens
- **Testing**: Testify + Testcontainers
- **Documentation**: Swagger/OpenAPI
- **Monitoring**: Prometheus + Grafana
- **Logging**: Structured logging with Zap

### Frontend
- **Framework**: Vue.js 3 with Composition API
- **Language**: TypeScript
- **State Management**: Pinia
- **Routing**: Vue Router 4
- **HTTP Client**: Axios with interceptors
- **UI Framework**: Vuetify 3 or Headless UI
- **Testing**: Vitest + Vue Test Utils
- **Build Tool**: Vite

### Infrastructure
- **Containerization**: Docker + Docker Compose
- **Orchestration**: Kubernetes (production)
- **CI/CD**: GitHub Actions
- **File Storage**: AWS S3 or MinIO
- **Email**: SendGrid or AWS SES
- **Monitoring**: Prometheus, Grafana, Jaeger

## Security Architecture

### Authentication & Authorization
```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   JWT Access    │    │  Refresh Token  │    │  Role-Based     │
│   Token         │    │  (HTTP-Only     │    │  Access Control │
│   (15 min TTL)  │    │   Cookie)       │    │  (RBAC)         │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

### Data Protection
- Password hashing with bcrypt
- API rate limiting
- Input validation and sanitization
- SQL injection prevention via GORM
- XSS protection
- CORS configuration

## Scalability Considerations

### Horizontal Scaling
- Stateless application design
- Database read replicas for queries
- Redis clustering for cache
- CDN for static assets
- Load balancing with session affinity

### Performance Optimization
- Database indexing strategy
- Query optimization
- Caching layers (Redis, application-level)
- Background job processing
- Image optimization and CDN

## Monitoring & Observability

### Metrics
- Application metrics (Prometheus)
- Database metrics
- Cache hit rates
- API response times
- Error rates

### Logging
- Structured logging
- Correlation IDs
- Request/response logging
- Error tracking

### Tracing
- Distributed tracing (Jaeger)
- Database query tracing
- External API call tracing