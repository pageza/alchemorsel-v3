# Alchemorsel v2 - System Overview & Architecture

## 1. Application Overview

Alchemorsel is an AI-powered recipe generation and management platform that creates personalized recipes based on user preferences, dietary restrictions, and allergies.

### Core Features
- **AI Recipe Generation**: Uses LLM (DeepSeek) to generate custom recipes
- **User Management**: Registration, authentication, profile management
- **Personalization**: Dietary preferences and allergy tracking
- **Recipe Management**: Search, save, favorite recipes
- **Semantic Search**: Vector-based recipe search using pgvector
- **Media Storage**: Profile picture management via S3

## 2. High-Level Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                         Load Balancer                           │
└─────────────────────────────────────────────────────────────────┘
                                  │
                                  ▼
┌─────────────────────────────────────────────────────────────────┐
│                        Docker Network                           │
├─────────────────────────────────────────────────────────────────┤
│  ┌─────────────┐     ┌─────────────┐     ┌─────────────┐     │
│  │   Vue 3     │     │   Go API    │     │   Nginx     │     │
│  │  Frontend   │────▶│   Backend   │◀────│   Reverse   │     │
│  │  (Port 80)  │     │ (Port 8080) │     │    Proxy    │     │
│  └─────────────┘     └─────────────┘     └─────────────┘     │
│                              │                                  │
│                              ▼                                  │
│  ┌─────────────────────────────────────────────────┐          │
│  │              Service Layer                       │          │
│  ├─────────────┬─────────────┬────────────────────┤          │
│  │ PostgreSQL  │    Redis    │   External APIs    │          │
│  │  + pgvector │   (Cache)   │  (DeepSeek, S3)   │          │
│  └─────────────┴─────────────┴────────────────────┘          │
└─────────────────────────────────────────────────────────────────┘
```

## 3. Technology Stack

### Backend (Go)
- **Framework**: Standard library with gorilla/mux or chi router
- **Database**: PostgreSQL 15+ with pgvector extension
- **Cache**: Redis 7+
- **ORM/Query Builder**: sqlx or pgx for type-safe queries
- **Authentication**: JWT tokens with refresh token rotation
- **API Documentation**: OpenAPI 3.0 with go-swagger
- **Testing**: Standard testing package with testify
- **Configuration**: Viper for environment management
- **Logging**: Zerolog or Zap
- **Validation**: go-playground/validator
- **Migration**: golang-migrate

### Frontend (Vue 3)
- **Framework**: Vue 3 with Composition API
- **Build Tool**: Vite
- **Language**: TypeScript with strict mode
- **State Management**: Pinia
- **Routing**: Vue Router 4
- **HTTP Client**: Axios with interceptors
- **UI Framework**: Tailwind CSS or Vuetify 3
- **Form Handling**: VeeValidate with Yup
- **Testing**: Vitest + Vue Test Utils
- **E2E Testing**: Playwright
- **Code Quality**: ESLint + Prettier

### Infrastructure
- **Containerization**: Docker with multi-stage builds
- **Orchestration**: Docker Compose for development
- **Reverse Proxy**: Nginx
- **Monitoring**: Prometheus + Grafana
- **Logging**: ELK stack or Loki
- **CI/CD**: GitHub Actions

## 4. Domain Model

### Core Entities

```
User
├── ID (UUID)
├── Email (unique)
├── Username (unique)
├── PasswordHash
├── Name
├── ProfilePictureURL
├── DietaryPreferences []
├── Allergies []
├── CreatedAt
├── UpdatedAt
└── DeletedAt (soft delete)

Recipe
├── ID (UUID)
├── UserID (creator)
├── Title
├── Description
├── Ingredients []
├── Instructions []
├── PrepTime
├── CookTime
├── Servings
├── Category
├── DietaryCategories []
├── Allergens []
├── NutritionalInfo
├── ImageURL
├── Embedding (vector)
├── IsPublic
├── CreatedAt
└── UpdatedAt

RecipeFavorite
├── UserID
├── RecipeID
└── CreatedAt

UserSession
├── ID
├── UserID
├── RefreshToken
├── ExpiresAt
└── CreatedAt
```

## 5. API Design Principles

### RESTful Endpoints
- **Versioned API**: `/api/v1/`
- **Resource-based URLs**: `/api/v1/recipes`, `/api/v1/users`
- **HTTP Methods**: GET, POST, PUT, PATCH, DELETE
- **Status Codes**: Proper HTTP status codes
- **Pagination**: Cursor-based for scalability
- **Filtering**: Query parameters for search/filter
- **Error Handling**: Consistent error response format

### Authentication Flow
1. **Registration**: Create user with validation
2. **Login**: Return access + refresh tokens
3. **Token Refresh**: Rotate refresh tokens
4. **Logout**: Invalidate refresh token
5. **Password Reset**: Email-based flow

## 6. Security Considerations

### Backend Security
- **Input Validation**: All inputs validated and sanitized
- **SQL Injection**: Parameterized queries only
- **Rate Limiting**: Per-endpoint and per-user limits
- **CORS**: Strict origin validation
- **Secrets Management**: Environment variables, never in code
- **Password Security**: bcrypt with appropriate cost factor
- **JWT Security**: Short-lived access tokens (15 min)
- **HTTPS**: TLS 1.3 only

### Frontend Security
- **XSS Prevention**: Content Security Policy
- **CSRF Protection**: Double-submit cookies
- **Secure Storage**: No sensitive data in localStorage
- **Input Sanitization**: Client-side validation
- **API Key Protection**: Proxy through backend

## 7. Performance Optimization

### Backend
- **Database Indexing**: Strategic indexes on search fields
- **Connection Pooling**: Optimal pool sizes
- **Query Optimization**: N+1 query prevention
- **Caching Strategy**: Redis for hot data
- **Async Processing**: Goroutines for I/O operations
- **Vector Search**: Proper pgvector indexing

### Frontend
- **Code Splitting**: Route-based lazy loading
- **Asset Optimization**: Image compression, WebP
- **Bundle Size**: Tree shaking, minification
- **Caching**: Service worker for offline capability
- **Virtual Scrolling**: For large lists
- **Debouncing**: Search and form inputs

## 8. Scalability Considerations

### Horizontal Scaling
- **Stateless Backend**: No server-side sessions
- **Database Replication**: Read replicas
- **Cache Distribution**: Redis cluster
- **Load Balancing**: Round-robin with health checks
- **Message Queue**: For async recipe generation

### Vertical Scaling
- **Database Optimization**: Query tuning
- **Connection Limits**: Proper configuration
- **Memory Management**: Efficient data structures
- **Goroutine Pools**: Controlled concurrency

## 9. Monitoring & Observability

### Metrics
- **Application Metrics**: Request rate, latency, errors
- **Business Metrics**: Recipe generation, user engagement
- **Infrastructure Metrics**: CPU, memory, disk usage

### Logging
- **Structured Logging**: JSON format
- **Log Levels**: Debug, Info, Warn, Error
- **Correlation IDs**: Request tracing
- **Log Aggregation**: Centralized logging

### Tracing
- **Distributed Tracing**: OpenTelemetry
- **Performance Profiling**: pprof for Go
- **Frontend Monitoring**: Web vitals

## 10. Development Workflow

### Git Strategy
- **Main Branch**: Production-ready code
- **Develop Branch**: Integration branch
- **Feature Branches**: `feature/description`
- **Hotfix Branches**: `hotfix/description`
- **Commit Convention**: Conventional commits

### Code Review Process
1. Feature branch created
2. Tests written (TDD)
3. Implementation completed
4. CI/CD passes
5. Code review by peers
6. Merge to develop
7. QA testing
8. Merge to main

### Testing Strategy
- **Unit Tests**: 80%+ coverage
- **Integration Tests**: API endpoints
- **E2E Tests**: Critical user flows
- **Performance Tests**: Load testing
- **Security Tests**: OWASP compliance