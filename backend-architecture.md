# Backend Architecture & Component Design

## 1. Directory Structure

```
alchemorsel-backend/
├── cmd/
│   └── api/
│       └── main.go                 # Application entry point
├── internal/
│   ├── config/
│   │   ├── config.go              # Configuration management
│   │   └── validation.go          # Config validation
│   ├── domain/
│   │   ├── user/
│   │   │   ├── entity.go          # User domain entity
│   │   │   ├── repository.go      # User repository interface
│   │   │   └── service.go         # User business logic
│   │   ├── recipe/
│   │   │   ├── entity.go          # Recipe domain entity
│   │   │   ├── repository.go      # Recipe repository interface
│   │   │   ├── service.go         # Recipe business logic
│   │   │   └── generator.go       # LLM recipe generation
│   │   └── auth/
│   │       ├── entity.go          # Auth-related entities
│   │       ├── service.go         # Authentication logic
│   │       └── jwt.go             # JWT token management
│   ├── infrastructure/
│   │   ├── database/
│   │   │   ├── postgres/
│   │   │   │   ├── connection.go  # Database connection
│   │   │   │   ├── migrations.go  # Migration runner
│   │   │   │   └── repository/
│   │   │   │       ├── user.go    # User repository impl
│   │   │   │       └── recipe.go  # Recipe repository impl
│   │   │   └── redis/
│   │   │       ├── connection.go  # Redis connection
│   │   │       └── cache.go       # Cache implementation
│   │   ├── storage/
│   │   │   └── s3/
│   │   │       └── client.go      # S3 client for images
│   │   └── external/
│   │       └── deepseek/
│   │           └── client.go      # DeepSeek API client
│   ├── interfaces/
│   │   ├── http/
│   │   │   ├── server.go          # HTTP server setup
│   │   │   ├── router.go          # Route definitions
│   │   │   ├── middleware/
│   │   │   │   ├── auth.go        # JWT middleware
│   │   │   │   ├── cors.go        # CORS middleware
│   │   │   │   ├── logging.go    # Request logging
│   │   │   │   ├── ratelimit.go  # Rate limiting
│   │   │   │   └── recovery.go   # Panic recovery
│   │   │   ├── handlers/
│   │   │   │   ├── auth.go        # Auth endpoints
│   │   │   │   ├── user.go        # User endpoints
│   │   │   │   ├── recipe.go      # Recipe endpoints
│   │   │   │   └── health.go      # Health check
│   │   │   └── dto/
│   │   │       ├── request/       # Request DTOs
│   │   │       └── response/      # Response DTOs
│   │   └── grpc/                  # Future gRPC support
│   └── pkg/
│       ├── errors/
│       │   └── errors.go          # Custom error types
│       ├── logger/
│       │   └── logger.go          # Logging wrapper
│       ├── validator/
│       │   └── validator.go       # Input validation
│       └── utils/
│           ├── password.go        # Password hashing
│           └── pagination.go      # Pagination helpers
├── migrations/
│   └── postgres/                  # SQL migration files
├── scripts/
│   ├── setup.sh                   # Development setup
│   └── seed.sh                    # Database seeding
├── tests/
│   ├── integration/               # Integration tests
│   ├── e2e/                       # End-to-end tests
│   └── fixtures/                  # Test data
├── api/
│   └── openapi.yaml              # OpenAPI specification
├── build/
│   └── Dockerfile                # Multi-stage Dockerfile
├── .env.example                  # Environment template
├── go.mod                        # Go modules
├── go.sum                        # Go modules checksum
└── Makefile                      # Build automation
```

## 2. Core Components

### 2.1 Configuration Component

```go
// internal/config/config.go
type Config struct {
    Server   ServerConfig
    Database DatabaseConfig
    Redis    RedisConfig
    Auth     AuthConfig
    External ExternalConfig
    Storage  StorageConfig
}

type ServerConfig struct {
    Host         string
    Port         int
    ReadTimeout  time.Duration
    WriteTimeout time.Duration
    IdleTimeout  time.Duration
}

type DatabaseConfig struct {
    Host         string
    Port         int
    User         string
    Password     string
    Database     string
    SSLMode      string
    MaxOpenConns int
    MaxIdleConns int
    MaxLifetime  time.Duration
}
```

### 2.2 Domain Layer

#### User Domain
```go
// internal/domain/user/entity.go
type User struct {
    ID                 uuid.UUID
    Email              string
    Username           string
    PasswordHash       string
    Name               string
    ProfilePictureURL  *string
    DietaryPreferences []string
    Allergies          []string
    CreatedAt          time.Time
    UpdatedAt          time.Time
    DeletedAt          *time.Time
}

// internal/domain/user/repository.go
type Repository interface {
    Create(ctx context.Context, user *User) error
    GetByID(ctx context.Context, id uuid.UUID) (*User, error)
    GetByEmail(ctx context.Context, email string) (*User, error)
    GetByUsername(ctx context.Context, username string) (*User, error)
    Update(ctx context.Context, user *User) error
    Delete(ctx context.Context, id uuid.UUID) error
}

// internal/domain/user/service.go
type Service interface {
    Register(ctx context.Context, req RegisterRequest) (*User, error)
    GetProfile(ctx context.Context, userID uuid.UUID) (*User, error)
    UpdateProfile(ctx context.Context, userID uuid.UUID, req UpdateRequest) error
    UploadProfilePicture(ctx context.Context, userID uuid.UUID, file io.Reader) error
}
```

#### Recipe Domain
```go
// internal/domain/recipe/entity.go
type Recipe struct {
    ID                uuid.UUID
    UserID            uuid.UUID
    Title             string
    Description       string
    Ingredients       []Ingredient
    Instructions      []string
    PrepTime          int // minutes
    CookTime          int // minutes
    Servings          int
    Category          string
    DietaryCategories []string
    Allergens         []string
    NutritionalInfo   NutritionalInfo
    ImageURL          *string
    Embedding         pgvector.Vector
    IsPublic          bool
    CreatedAt         time.Time
    UpdatedAt         time.Time
}

type Ingredient struct {
    Name     string
    Amount   float64
    Unit     string
    Optional bool
}

// internal/domain/recipe/repository.go
type Repository interface {
    Create(ctx context.Context, recipe *Recipe) error
    GetByID(ctx context.Context, id uuid.UUID) (*Recipe, error)
    Search(ctx context.Context, params SearchParams) (*SearchResult, error)
    GetUserFavorites(ctx context.Context, userID uuid.UUID, pagination Pagination) ([]*Recipe, error)
    AddFavorite(ctx context.Context, userID, recipeID uuid.UUID) error
    RemoveFavorite(ctx context.Context, userID, recipeID uuid.UUID) error
}
```

### 2.3 Infrastructure Layer

#### Database Connection
```go
// internal/infrastructure/database/postgres/connection.go
type DB struct {
    *sqlx.DB
    config DatabaseConfig
}

func NewConnection(config DatabaseConfig) (*DB, error) {
    // Connection string with pgvector extension
    // Connection pooling setup
    // Health check implementation
}
```

#### Repository Implementation
```go
// internal/infrastructure/database/postgres/repository/user.go
type userRepository struct {
    db *DB
}

func (r *userRepository) Create(ctx context.Context, user *domain.User) error {
    query := `
        INSERT INTO users (id, email, username, password_hash, name, 
                          dietary_preferences, allergies, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
    `
    // Implementation with proper error handling
}
```

### 2.4 HTTP Layer

#### Router Setup
```go
// internal/interfaces/http/router.go
func SetupRoutes(r *chi.Mux, services Services) {
    r.Route("/api/v1", func(r chi.Router) {
        // Public routes
        r.Group(func(r chi.Router) {
            r.Post("/auth/register", handlers.Register(services.User))
            r.Post("/auth/login", handlers.Login(services.Auth))
            r.Post("/auth/refresh", handlers.RefreshToken(services.Auth))
        })
        
        // Protected routes
        r.Group(func(r chi.Router) {
            r.Use(middleware.Authenticate)
            
            r.Route("/users", func(r chi.Router) {
                r.Get("/profile", handlers.GetProfile(services.User))
                r.Put("/profile", handlers.UpdateProfile(services.User))
                r.Post("/profile/picture", handlers.UploadProfilePicture(services.User))
            })
            
            r.Route("/recipes", func(r chi.Router) {
                r.Get("/", handlers.SearchRecipes(services.Recipe))
                r.Post("/", handlers.CreateRecipe(services.Recipe))
                r.Get("/{id}", handlers.GetRecipe(services.Recipe))
                r.Post("/{id}/favorite", handlers.AddFavorite(services.Recipe))
                r.Delete("/{id}/favorite", handlers.RemoveFavorite(services.Recipe))
            })
            
            r.Post("/llm/generate", handlers.GenerateRecipe(services.Recipe))
        })
    })
}
```

#### Middleware
```go
// internal/interfaces/http/middleware/auth.go
func Authenticate(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Extract token from Authorization header
        // Validate JWT token
        // Add user context
        // Call next handler
    })
}

// internal/interfaces/http/middleware/ratelimit.go
func RateLimit(rpm int) func(http.Handler) http.Handler {
    // Redis-based rate limiting
    // Per-user and per-IP limiting
}
```

## 3. Database Schema

### 3.1 Users Table
```sql
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    username VARCHAR(50) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    profile_picture_url TEXT,
    dietary_preferences TEXT[],
    allergies TEXT[],
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_username ON users(username);
```

### 3.2 Recipes Table
```sql
CREATE EXTENSION IF NOT EXISTS vector;

CREATE TABLE recipes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    ingredients JSONB NOT NULL,
    instructions TEXT[] NOT NULL,
    prep_time INTEGER,
    cook_time INTEGER,
    servings INTEGER,
    category VARCHAR(100),
    dietary_categories TEXT[],
    allergens TEXT[],
    nutritional_info JSONB,
    image_url TEXT,
    embedding vector(1536),
    is_public BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_recipes_user_id ON recipes(user_id);
CREATE INDEX idx_recipes_embedding ON recipes USING ivfflat (embedding vector_cosine_ops);
```

### 3.3 Recipe Favorites Table
```sql
CREATE TABLE recipe_favorites (
    user_id UUID NOT NULL REFERENCES users(id),
    recipe_id UUID NOT NULL REFERENCES recipes(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    PRIMARY KEY (user_id, recipe_id)
);
```

### 3.4 User Sessions Table
```sql
CREATE TABLE user_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id),
    refresh_token VARCHAR(255) UNIQUE NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_user_sessions_user_id ON user_sessions(user_id);
CREATE INDEX idx_user_sessions_refresh_token ON user_sessions(refresh_token);
```

## 4. External Service Integration

### 4.1 DeepSeek Client
```go
// internal/infrastructure/external/deepseek/client.go
type Client struct {
    apiKey     string
    apiURL     string
    httpClient *http.Client
}

type GenerateRecipeRequest struct {
    UserPreferences UserPreferences
    Constraints     RecipeConstraints
    Style          string
}

func (c *Client) GenerateRecipe(ctx context.Context, req GenerateRecipeRequest) (*Recipe, error) {
    // Build prompt based on user preferences
    // Call DeepSeek API
    // Parse response
    // Generate embedding
    // Return structured recipe
}
```

### 4.2 S3 Storage Client
```go
// internal/infrastructure/storage/s3/client.go
type Client struct {
    bucket string
    region string
    client *s3.Client
}

func (c *Client) UploadImage(ctx context.Context, key string, data io.Reader) (string, error) {
    // Upload to S3
    // Return public URL
}
```

## 5. Testing Strategy

### 5.1 Unit Tests
```go
// internal/domain/user/service_test.go
func TestUserService_Register(t *testing.T) {
    // Mock repository
    // Test various scenarios
    // Assert expectations
}
```

### 5.2 Integration Tests
```go
// tests/integration/api_test.go
func TestRecipeAPI(t *testing.T) {
    // Setup test database
    // Create test server
    // Run API tests
    // Clean up
}
```

### 5.3 Repository Tests
```go
// internal/infrastructure/database/postgres/repository/recipe_test.go
func TestRecipeRepository_Search(t *testing.T) {
    // Use test container
    // Insert test data
    // Test search functionality
    // Verify results
}
```

## 6. Error Handling

### 6.1 Custom Error Types
```go
// internal/pkg/errors/errors.go
type AppError struct {
    Code    string
    Message string
    Status  int
    Details map[string]interface{}
}

var (
    ErrUserNotFound     = AppError{Code: "USER_NOT_FOUND", Status: 404}
    ErrInvalidInput     = AppError{Code: "INVALID_INPUT", Status: 400}
    ErrUnauthorized     = AppError{Code: "UNAUTHORIZED", Status: 401}
    ErrDuplicateEmail   = AppError{Code: "DUPLICATE_EMAIL", Status: 409}
)
```

## 7. Logging & Monitoring

### 7.1 Structured Logging
```go
// internal/pkg/logger/logger.go
type Logger interface {
    Debug(msg string, fields ...Field)
    Info(msg string, fields ...Field)
    Warn(msg string, fields ...Field)
    Error(msg string, fields ...Field)
}

// Usage example
logger.Info("Recipe created", 
    Field("user_id", userID),
    Field("recipe_id", recipeID),
    Field("duration_ms", duration))
```

### 7.2 Metrics Collection
```go
// internal/interfaces/http/middleware/metrics.go
func Metrics() func(http.Handler) http.Handler {
    // Prometheus metrics
    // Request duration
    // Request count
    // Error rate
}
```