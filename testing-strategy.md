# Comprehensive Testing Strategy

## 1. Testing Philosophy

### Test-Driven Development (TDD) Principles
1. **Red-Green-Refactor**: Write failing test → Make it pass → Improve code
2. **Test First**: Tests define behavior before implementation
3. **Fast Feedback**: Tests run quickly for rapid iteration
4. **Living Documentation**: Tests document system behavior
5. **Confidence in Changes**: Comprehensive tests enable refactoring

### Testing Pyramid
```
         E2E Tests (5%)
        /           \
    Integration (15%) 
   /                 \
  Component Tests (30%)
 /                     \
Unit Tests (50%)
```

## 2. Backend Testing Strategy

### 2.1 Unit Tests

#### Repository Tests
```go
// internal/infrastructure/database/postgres/repository/user_test.go
package repository

import (
    "context"
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestUserRepository_Create(t *testing.T) {
    db := setupTestDB(t)
    defer cleanupTestDB(db)
    
    repo := NewUserRepository(db)
    
    tests := []struct {
        name    string
        user    *domain.User
        wantErr bool
        errType error
    }{
        {
            name: "valid user creation",
            user: &domain.User{
                Email:    "test@example.com",
                Username: "testuser",
                Name:     "Test User",
            },
            wantErr: false,
        },
        {
            name: "duplicate email",
            user: &domain.User{
                Email:    "duplicate@example.com",
                Username: "unique",
            },
            wantErr: true,
            errType: ErrDuplicateEmail,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := repo.Create(context.Background(), tt.user)
            
            if tt.wantErr {
                assert.Error(t, err)
                assert.ErrorIs(t, err, tt.errType)
            } else {
                assert.NoError(t, err)
                assert.NotEmpty(t, tt.user.ID)
            }
        })
    }
}
```

#### Service Tests
```go
// internal/domain/recipe/service_test.go
package recipe

import (
    "context"
    "testing"
    "github.com/stretchr/testify/mock"
)

type MockRepository struct {
    mock.Mock
}

func (m *MockRepository) Create(ctx context.Context, recipe *Recipe) error {
    args := m.Called(ctx, recipe)
    return args.Error(0)
}

func TestRecipeService_GenerateRecipe(t *testing.T) {
    mockRepo := new(MockRepository)
    mockLLM := new(MockLLMClient)
    service := NewRecipeService(mockRepo, mockLLM)
    
    mockLLM.On("Generate", mock.Anything, mock.Anything).Return(&LLMResponse{
        Title: "Test Recipe",
        Ingredients: []string{"ingredient1", "ingredient2"},
    }, nil)
    
    mockRepo.On("Create", mock.Anything, mock.Anything).Return(nil)
    
    recipe, err := service.GenerateRecipe(context.Background(), GenerateRequest{
        Style: "italian",
    })
    
    assert.NoError(t, err)
    assert.Equal(t, "Test Recipe", recipe.Title)
    mockLLM.AssertExpectations(t)
    mockRepo.AssertExpectations(t)
}
```

#### Utility Tests
```go
// internal/pkg/utils/password_test.go
package utils

import "testing"

func TestHashPassword(t *testing.T) {
    password := "SecurePassword123!"
    
    hash, err := HashPassword(password)
    assert.NoError(t, err)
    assert.NotEmpty(t, hash)
    assert.NotEqual(t, password, hash)
    
    // Verify password
    valid := CheckPassword(password, hash)
    assert.True(t, valid)
    
    // Wrong password
    invalid := CheckPassword("WrongPassword", hash)
    assert.False(t, invalid)
}
```

### 2.2 Integration Tests

#### Database Integration
```go
// tests/integration/database_test.go
package integration

import (
    "testing"
    "github.com/testcontainers/testcontainers-go"
    "github.com/testcontainers/testcontainers-go/modules/postgres"
)

func TestRecipeSearch(t *testing.T) {
    ctx := context.Background()
    
    // Start PostgreSQL container
    pgContainer, err := postgres.RunContainer(ctx,
        testcontainers.WithImage("pgvector/pgvector:pg15"),
        postgres.WithDatabase("testdb"),
        postgres.WithUsername("test"),
        postgres.WithPassword("test"),
    )
    require.NoError(t, err)
    defer pgContainer.Terminate(ctx)
    
    // Get connection string
    connStr, err := pgContainer.ConnectionString(ctx)
    require.NoError(t, err)
    
    // Setup database
    db := setupDatabase(connStr)
    runMigrations(db)
    
    // Test search functionality
    repo := repository.NewRecipeRepository(db)
    
    // Insert test recipes
    recipes := createTestRecipes()
    for _, recipe := range recipes {
        err := repo.Create(ctx, recipe)
        require.NoError(t, err)
    }
    
    // Test vector search
    results, err := repo.Search(ctx, SearchParams{
        Query: "pasta italian",
        Limit: 10,
    })
    
    assert.NoError(t, err)
    assert.NotEmpty(t, results)
    assert.Contains(t, results[0].Title, "Italian")
}
```

#### API Integration Tests
```go
// tests/integration/api_test.go
package integration

func TestRecipeAPI(t *testing.T) {
    app := setupTestApp()
    
    t.Run("Create Recipe", func(t *testing.T) {
        token := authenticateTestUser(t, app)
        
        body := `{
            "title": "Test Recipe",
            "ingredients": [{"name": "flour", "amount": 2, "unit": "cups"}],
            "instructions": ["Mix", "Bake"]
        }`
        
        req := httptest.NewRequest("POST", "/api/v1/recipes", strings.NewReader(body))
        req.Header.Set("Authorization", "Bearer "+token)
        req.Header.Set("Content-Type", "application/json")
        
        resp := httptest.NewRecorder()
        app.ServeHTTP(resp, req)
        
        assert.Equal(t, http.StatusCreated, resp.Code)
        
        var result map[string]interface{}
        json.Unmarshal(resp.Body.Bytes(), &result)
        assert.Equal(t, "Test Recipe", result["data"].(map[string]interface{})["recipe"].(map[string]interface{})["title"])
    })
}
```

### 2.3 End-to-End Tests

```go
// tests/e2e/recipe_flow_test.go
package e2e

func TestCompleteRecipeFlow(t *testing.T) {
    // Start all services
    compose := startDockerCompose(t)
    defer compose.Down()
    
    // Wait for services to be ready
    waitForServices(t)
    
    client := newAPIClient()
    
    // 1. Register user
    user := registerUser(t, client, "test@example.com", "password")
    
    // 2. Login
    tokens := login(t, client, "test@example.com", "password")
    
    // 3. Generate recipe
    recipe := generateRecipe(t, client, tokens.AccessToken, GenerateRequest{
        Style: "italian",
        Ingredients: []string{"tomatoes", "basil"},
    })
    
    // 4. Search for recipe
    results := searchRecipes(t, client, tokens.AccessToken, "italian tomato")
    assert.Contains(t, results, recipe.ID)
    
    // 5. Add to favorites
    addFavorite(t, client, tokens.AccessToken, recipe.ID)
    
    // 6. Verify in favorites
    favorites := getFavorites(t, client, tokens.AccessToken)
    assert.Contains(t, favorites, recipe.ID)
}
```

## 3. Frontend Testing Strategy

### 3.1 Unit Tests

#### Store Tests
```typescript
// tests/unit/stores/auth.store.spec.ts
import { setActivePinia, createPinia } from 'pinia';
import { describe, it, expect, beforeEach, vi } from 'vitest';
import { useAuthStore } from '@/stores/auth.store';
import { AuthService } from '@/services/auth.service';

vi.mock('@/services/auth.service');

describe('Auth Store', () => {
  beforeEach(() => {
    setActivePinia(createPinia());
    vi.clearAllMocks();
  });
  
  describe('login', () => {
    it('should login successfully', async () => {
      const mockResponse = {
        user: { id: '1', email: 'test@example.com' },
        tokens: { accessToken: 'access', refreshToken: 'refresh' }
      };
      
      vi.mocked(AuthService.login).mockResolvedValue(mockResponse);
      
      const store = useAuthStore();
      await store.login({ email: 'test@example.com', password: 'password' });
      
      expect(store.isAuthenticated).toBe(true);
      expect(store.currentUser).toEqual(mockResponse.user);
    });
    
    it('should handle login error', async () => {
      vi.mocked(AuthService.login).mockRejectedValue(new Error('Invalid credentials'));
      
      const store = useAuthStore();
      await expect(store.login({ email: 'test@example.com', password: 'wrong' }))
        .rejects.toThrow('Invalid credentials');
      
      expect(store.isAuthenticated).toBe(false);
    });
  });
});
```

#### Composable Tests
```typescript
// tests/unit/composables/useDebounce.spec.ts
import { ref } from 'vue';
import { describe, it, expect, vi } from 'vitest';
import { useDebounce } from '@/composables/useDebounce';

describe('useDebounce', () => {
  beforeEach(() => {
    vi.useFakeTimers();
  });
  
  afterEach(() => {
    vi.restoreAllMocks();
  });
  
  it('should debounce value changes', async () => {
    const source = ref('initial');
    const debounced = useDebounce(source, 300);
    
    expect(debounced.value).toBe('initial');
    
    source.value = 'changed';
    expect(debounced.value).toBe('initial');
    
    vi.advanceTimersByTime(300);
    expect(debounced.value).toBe('changed');
  });
});
```

### 3.2 Component Tests

```typescript
// tests/component/RecipeCard.spec.ts
import { mount } from '@vue/test-utils';
import { describe, it, expect, vi } from 'vitest';
import RecipeCard from '@/components/recipe/RecipeCard.vue';
import { createTestingPinia } from '@pinia/testing';

describe('RecipeCard', () => {
  const mockRecipe = {
    id: '1',
    title: 'Delicious Pasta',
    description: 'A wonderful pasta dish',
    prepTime: 15,
    cookTime: 30,
    servings: 4,
    imageUrl: '/test.jpg',
    category: 'dinner',
    isFavorited: false
  };
  
  it('renders recipe information', () => {
    const wrapper = mount(RecipeCard, {
      props: { recipe: mockRecipe },
      global: {
        plugins: [createTestingPinia()]
      }
    });
    
    expect(wrapper.text()).toContain('Delicious Pasta');
    expect(wrapper.text()).toContain('45 min total');
    expect(wrapper.text()).toContain('4 servings');
  });
  
  it('emits toggle-favorite event', async () => {
    const wrapper = mount(RecipeCard, {
      props: { recipe: mockRecipe },
      global: {
        plugins: [createTestingPinia()]
      }
    });
    
    await wrapper.find('[data-testid="favorite-button"]').trigger('click');
    
    expect(wrapper.emitted('toggle-favorite')).toBeTruthy();
    expect(wrapper.emitted('toggle-favorite')[0]).toEqual([mockRecipe.id]);
  });
  
  it('shows favorited state', () => {
    const wrapper = mount(RecipeCard, {
      props: { 
        recipe: { ...mockRecipe, isFavorited: true }
      },
      global: {
        plugins: [createTestingPinia()]
      }
    });
    
    const favoriteIcon = wrapper.find('[data-testid="favorite-icon"]');
    expect(favoriteIcon.classes()).toContain('text-red-500');
  });
});
```

### 3.3 E2E Tests with Playwright

```typescript
// tests/e2e/auth.spec.ts
import { test, expect } from '@playwright/test';

test.describe('Authentication Flow', () => {
  test('should register, login, and access protected routes', async ({ page }) => {
    // Navigate to register page
    await page.goto('/auth/register');
    
    // Fill registration form
    await page.fill('[data-testid="email-input"]', 'newuser@example.com');
    await page.fill('[data-testid="username-input"]', 'newuser');
    await page.fill('[data-testid="password-input"]', 'SecurePassword123!');
    await page.fill('[data-testid="name-input"]', 'New User');
    
    // Select dietary preferences
    await page.check('[data-testid="dietary-vegan"]');
    await page.check('[data-testid="dietary-gluten-free"]');
    
    // Submit form
    await page.click('[data-testid="register-button"]');
    
    // Should redirect to dashboard
    await expect(page).toHaveURL('/dashboard');
    await expect(page.locator('h1')).toContainText('Welcome, New User');
    
    // Logout
    await page.click('[data-testid="logout-button"]');
    
    // Login with created account
    await page.goto('/auth/login');
    await page.fill('[data-testid="email-input"]', 'newuser@example.com');
    await page.fill('[data-testid="password-input"]', 'SecurePassword123!');
    await page.click('[data-testid="login-button"]');
    
    // Should be back on dashboard
    await expect(page).toHaveURL('/dashboard');
  });
});
```

```typescript
// tests/e2e/recipe-generation.spec.ts
import { test, expect } from '@playwright/test';

test.describe('Recipe Generation', () => {
  test.beforeEach(async ({ page }) => {
    // Login before each test
    await loginTestUser(page);
  });
  
  test('should generate and save a recipe', async ({ page }) => {
    // Navigate to generate page
    await page.goto('/dashboard/generate');
    
    // Fill generation form
    await page.selectOption('[data-testid="style-select"]', 'italian');
    await page.fill('[data-testid="ingredients-input"]', 'tomatoes, basil, mozzarella');
    await page.fill('[data-testid="servings-input"]', '4');
    await page.selectOption('[data-testid="time-select"]', '30');
    
    // Generate recipe
    await page.click('[data-testid="generate-button"]');
    
    // Wait for generation to complete
    await expect(page.locator('[data-testid="recipe-result"]')).toBeVisible({
      timeout: 30000
    });
    
    // Verify recipe details
    const recipeTitle = await page.locator('[data-testid="recipe-title"]').textContent();
    expect(recipeTitle).toBeTruthy();
    
    // Save recipe
    await page.click('[data-testid="save-recipe-button"]');
    
    // Should redirect to recipe detail page
    await expect(page).toHaveURL(/\/recipes\/[\w-]+/);
    
    // Verify saved
    await expect(page.locator('[data-testid="success-message"]')).toContainText('Recipe saved');
  });
});
```

## 4. Performance Testing

### 4.1 Load Testing with k6

```javascript
// scripts/load-tests/recipe-search.js
import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
  stages: [
    { duration: '2m', target: 100 }, // Ramp up
    { duration: '5m', target: 100 }, // Stay at 100 users
    { duration: '2m', target: 200 }, // Ramp up more
    { duration: '5m', target: 200 }, // Stay at 200 users
    { duration: '2m', target: 0 },   // Ramp down
  ],
  thresholds: {
    http_req_duration: ['p(95)<500'], // 95% of requests under 500ms
    http_req_failed: ['rate<0.1'],    // Error rate under 10%
  },
};

export default function() {
  const token = getAuthToken();
  
  const params = {
    headers: {
      'Authorization': `Bearer ${token}`,
      'Content-Type': 'application/json',
    },
  };
  
  // Search recipes
  const searchRes = http.get(
    'http://localhost:8080/api/v1/recipes?q=pasta&dietary=vegan',
    params
  );
  
  check(searchRes, {
    'search status is 200': (r) => r.status === 200,
    'search response time < 500ms': (r) => r.timings.duration < 500,
    'search returns results': (r) => JSON.parse(r.body).data.recipes.length > 0,
  });
  
  sleep(1);
}
```

### 4.2 Database Performance Tests

```go
// tests/performance/database_bench_test.go
package performance

func BenchmarkRecipeSearch(b *testing.B) {
    db := setupBenchDB()
    repo := repository.NewRecipeRepository(db)
    
    // Seed with 10,000 recipes
    seedRecipes(db, 10000)
    
    b.ResetTimer()
    
    for i := 0; i < b.N; i++ {
        _, err := repo.Search(context.Background(), SearchParams{
            Query: "pasta italian garlic",
            Limit: 20,
        })
        if err != nil {
            b.Fatal(err)
        }
    }
}

func BenchmarkRecipeInsert(b *testing.B) {
    db := setupBenchDB()
    repo := repository.NewRecipeRepository(db)
    
    recipes := generateBenchmarkRecipes(b.N)
    
    b.ResetTimer()
    
    for i := 0; i < b.N; i++ {
        err := repo.Create(context.Background(), recipes[i])
        if err != nil {
            b.Fatal(err)
        }
    }
}
```

## 5. Security Testing

### 5.1 SQL Injection Tests

```go
// tests/security/sql_injection_test.go
package security

func TestSQLInjection(t *testing.T) {
    app := setupTestApp()
    
    maliciousInputs := []string{
        "'; DROP TABLE users; --",
        "1' OR '1'='1",
        "admin'--",
        "1' UNION SELECT * FROM users--",
    }
    
    for _, input := range maliciousInputs {
        t.Run(input, func(t *testing.T) {
            resp := searchRecipes(app, input)
            
            // Should not return error or expose data
            assert.Equal(t, http.StatusOK, resp.StatusCode)
            
            var result map[string]interface{}
            json.NewDecoder(resp.Body).Decode(&result)
            
            // Verify no SQL error exposed
            assert.NotContains(t, result, "syntax error")
            assert.NotContains(t, result, "SQL")
        })
    }
}
```

### 5.2 XSS Prevention Tests

```typescript
// tests/security/xss.spec.ts
import { test, expect } from '@playwright/test';

test.describe('XSS Prevention', () => {
  const xssPayloads = [
    '<script>alert("XSS")</script>',
    '<img src=x onerror=alert("XSS")>',
    'javascript:alert("XSS")',
    '<svg onload=alert("XSS")>',
  ];
  
  test('should sanitize user input in recipes', async ({ page }) => {
    await loginTestUser(page);
    
    for (const payload of xssPayloads) {
      await page.goto('/dashboard/recipes/create');
      
      // Try to inject XSS in title
      await page.fill('[data-testid="title-input"]', payload);
      await page.fill('[data-testid="description-input"]', 'Normal description');
      
      // Save recipe
      await page.click('[data-testid="save-button"]');
      
      // Navigate to recipe
      await page.waitForURL(/\/recipes\/[\w-]+/);
      
      // Verify script didn't execute
      const alertFired = await page.evaluate(() => {
        return window.xssAlertFired || false;
      });
      
      expect(alertFired).toBe(false);
      
      // Verify content is escaped
      const titleText = await page.locator('h1').textContent();
      expect(titleText).not.toContain('<script>');
    }
  });
});
```

## 6. Test Data Management

### 6.1 Fixtures

```go
// tests/fixtures/recipes.go
package fixtures

var TestRecipes = []Recipe{
    {
        Title: "Classic Margherita Pizza",
        Category: "dinner",
        DietaryCategories: []string{"vegetarian"},
        Ingredients: []Ingredient{
            {Name: "Pizza dough", Amount: 1, Unit: "lb"},
            {Name: "Tomato sauce", Amount: 1, Unit: "cup"},
            {Name: "Mozzarella", Amount: 8, Unit: "oz"},
            {Name: "Basil", Amount: 10, Unit: "leaves"},
        },
    },
    // More test recipes...
}
```

### 6.2 Test Factories

```typescript
// tests/factories/recipe.factory.ts
import { Factory } from 'fishery';
import type { Recipe } from '@/types/recipe.types';

export const recipeFactory = Factory.define<Recipe>(() => ({
  id: faker.datatype.uuid(),
  title: faker.lorem.words(3),
  description: faker.lorem.paragraph(),
  prepTime: faker.datatype.number({ min: 5, max: 60 }),
  cookTime: faker.datatype.number({ min: 10, max: 120 }),
  servings: faker.datatype.number({ min: 1, max: 8 }),
  ingredients: ingredientFactory.buildList(5),
  instructions: faker.lorem.paragraphs(5).split('\n'),
  category: faker.helpers.arrayElement(['breakfast', 'lunch', 'dinner']),
  createdAt: faker.date.past().toISOString(),
}));
```

## 7. Continuous Integration Tests

### 7.1 GitHub Actions Workflow

```yaml
# .github/workflows/test.yml
name: Test Suite

on: [push, pull_request]

jobs:
  backend-tests:
    runs-on: ubuntu-latest
    services:
      postgres:
        image: pgvector/pgvector:pg15
        env:
          POSTGRES_PASSWORD: test
        options: >-
          --health-cmd pg_isready
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5
      redis:
        image: redis:7
        options: >-
          --health-cmd "redis-cli ping"
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5
    
    steps:
    - uses: actions/checkout@v3
    
    - name: Set up Go
      uses: actions/setup-go@v4
      with:
        go-version: '1.21'
    
    - name: Install dependencies
      run: go mod download
    
    - name: Run unit tests
      run: go test -v -cover -race ./internal/...
    
    - name: Run integration tests
      run: go test -v -tags=integration ./tests/integration/...
      env:
        DATABASE_URL: postgres://postgres:test@localhost:5432/test?sslmode=disable
        REDIS_URL: redis://localhost:6379
    
    - name: Generate coverage report
      run: go test -coverprofile=coverage.out ./...
    
    - name: Upload coverage
      uses: codecov/codecov-action@v3
      with:
        file: ./coverage.out

  frontend-tests:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v3
    
    - name: Setup Node.js
      uses: actions/setup-node@v3
      with:
        node-version: '18'
        cache: 'npm'
    
    - name: Install dependencies
      run: npm ci
    
    - name: Run linting
      run: npm run lint
    
    - name: Run unit tests
      run: npm run test:unit -- --coverage
    
    - name: Run component tests
      run: npm run test:component
    
    - name: Build application
      run: npm run build
    
    - name: Upload coverage
      uses: codecov/codecov-action@v3
      with:
        file: ./coverage/lcov.info

  e2e-tests:
    runs-on: ubuntu-latest
    needs: [backend-tests, frontend-tests]
    steps:
    - uses: actions/checkout@v3
    
    - name: Start services
      run: docker-compose -f docker-compose.test.yml up -d
    
    - name: Wait for services
      run: |
        timeout 30 bash -c 'until curl -f http://localhost:8080/api/v1/health; do sleep 1; done'
    
    - name: Run E2E tests
      run: npm run test:e2e
    
    - name: Upload test artifacts
      if: failure()
      uses: actions/upload-artifact@v3
      with:
        name: e2e-artifacts
        path: |
          tests/e2e/screenshots/
          tests/e2e/videos/
```

## 8. Test Utilities

### 8.1 Backend Test Helpers

```go
// tests/helpers/auth.go
package helpers

import (
    "testing"
    "net/http/httptest"
)

func AuthenticateTestUser(t *testing.T, email string) string {
    user := CreateTestUser(t, email)
    token, err := GenerateTestToken(user.ID)
    require.NoError(t, err)
    return token
}

func AuthorizedRequest(t *testing.T, method, path, token string) *httptest.ResponseRecorder {
    req := httptest.NewRequest(method, path, nil)
    req.Header.Set("Authorization", "Bearer "+token)
    rec := httptest.NewRecorder()
    return rec
}
```

```go
// tests/helpers/database.go
package helpers

import (
    "database/sql"
    "testing"
    _ "github.com/golang-migrate/migrate/v4/source/file"
    _ "github.com/golang-migrate/migrate/v4/database/postgres"
)

func SetupTestDB(t *testing.T) *sql.DB {
    db, err := sql.Open("postgres", GetTestDatabaseURL())
    require.NoError(t, err)
    
    // Run migrations
    RunMigrations(t, db)
    
    // Clean database before tests
    CleanDatabase(t, db)
    
    return db
}

func CleanDatabase(t *testing.T, db *sql.DB) {
    tables := []string{"recipe_favorites", "recipes", "user_sessions", "users"}
    
    for _, table := range tables {
        _, err := db.Exec("TRUNCATE TABLE " + table + " CASCADE")
        require.NoError(t, err)
    }
}
```

### 8.2 Frontend Test Helpers

```typescript
// tests/helpers/auth.ts
export async function loginTestUser(page: Page, email = 'test@example.com', password = 'password') {
  await page.goto('/auth/login');
  await page.fill('[data-testid="email-input"]', email);
  await page.fill('[data-testid="password-input"]', password);
  await page.click('[data-testid="login-button"]');
  await page.waitForURL('/dashboard');
}

export function createAuthenticatedContext() {
  const tokens = {
    accessToken: 'test-access-token',
    refreshToken: 'test-refresh-token',
  };
  
  localStorage.setItem('auth_tokens', JSON.stringify(tokens));
  return tokens;
}
```

```typescript
// tests/helpers/mock-server.ts
import { rest } from 'msw';
import { setupServer } from 'msw/node';

export const mockServer = setupServer(
  rest.post('/api/v1/auth/login', (req, res, ctx) => {
    return res(
      ctx.json({
        success: true,
        data: {
          user: { id: '1', email: 'test@example.com' },
          tokens: {
            accessToken: 'mock-access-token',
            refreshToken: 'mock-refresh-token',
          },
        },
      })
    );
  }),
  
  rest.get('/api/v1/recipes', (req, res, ctx) => {
    return res(
      ctx.json({
        success: true,
        data: {
          recipes: recipeFactory.buildList(10),
        },
        pagination: {
          page: 1,
          total: 10,
          totalPages: 1,
        },
      })
    );
  })
);
```

## 9. Test Monitoring & Reporting

### 9.1 Test Coverage Requirements

```yaml
# .codecov.yml
coverage:
  status:
    project:
      default:
        target: 80%
        threshold: 2%
    patch:
      default:
        target: 90%
        threshold: 5%

comment:
  layout: "reach,diff,flags,files,footer"
  behavior: default
  require_changes: false
```

### 9.2 Test Report Generation

```javascript
// vitest.config.ts
import { defineConfig } from 'vitest/config';

export default defineConfig({
  test: {
    coverage: {
      provider: 'v8',
      reporter: ['text', 'json', 'html', 'lcov'],
      exclude: [
        'node_modules/**',
        'tests/**',
        '**/*.d.ts',
        '**/*.config.*',
        '**/mockData/**',
      ],
      thresholds: {
        lines: 80,
        functions: 80,
        branches: 80,
        statements: 80,
      },
    },
    reporters: ['default', 'html', 'junit'],
    outputFile: {
      junit: './test-results/junit.xml',
      html: './test-results/index.html',
    },
  },
});
```

## 10. Test Best Practices

### 10.1 Test Naming Conventions

```go
// Good test names
func TestUserService_Register_WithValidData_CreatesUser(t *testing.T) {}
func TestUserService_Register_WithDuplicateEmail_ReturnsError(t *testing.T) {}
func TestRecipeRepository_Search_WithVectorSimilarity_ReturnsRelevantResults(t *testing.T) {}

// Component: TestXxx_Method_Condition_Result
```

### 10.2 Test Organization

```go
func TestRecipeService(t *testing.T) {
    t.Run("GenerateRecipe", func(t *testing.T) {
        t.Run("with valid prompt", func(t *testing.T) {
            // Test implementation
        })
        
        t.Run("with invalid prompt", func(t *testing.T) {
            // Test implementation
        })
    })
    
    t.Run("UpdateRecipe", func(t *testing.T) {
        t.Run("as owner", func(t *testing.T) {
            // Test implementation
        })
        
        t.Run("as non-owner", func(t *testing.T) {
            // Test implementation
        })
    })
}
```

### 10.3 Test Isolation

```typescript
// Each test should be independent
describe('RecipeStore', () => {
  let store: ReturnType<typeof useRecipeStore>;
  
  beforeEach(() => {
    // Fresh store for each test
    setActivePinia(createPinia());
    store = useRecipeStore();
    // Reset mocks
    vi.clearAllMocks();
  });
  
  afterEach(() => {
    // Clean up
    cleanup();
  });
  
  it('should not depend on other tests', () => {
    // Test implementation
  });
});
```

### 10.4 Assertion Guidelines

```go
// Be specific with assertions
assert.Equal(t, expected, actual, "Recipe title should match")
assert.Len(t, recipes, 10, "Should return exactly 10 recipes")
assert.Contains(t, recipe.DietaryCategories, "vegan", "Recipe should be marked as vegan")

// Use appropriate matchers
assert.Eventually(t, func() bool {
    return cache.Get(key) != nil
}, 5*time.Second, 100*time.Millisecond, "Cache should be populated")

// Group related assertions
assert.NoError(t, err)
if assert.NotNil(t, recipe) {
    assert.Equal(t, "Test Recipe", recipe.Title)
    assert.Len(t, recipe.Ingredients, 5)
}
```

## 11. Performance Testing Benchmarks

### Target Metrics

| Operation | Target | Acceptable | Critical |
|-----------|--------|------------|----------|
| API Response (p95) | < 200ms | < 500ms | > 1000ms |
| Recipe Search | < 300ms | < 700ms | > 1500ms |
| Recipe Generation | < 5s | < 10s | > 20s |
| Image Upload | < 2s | < 5s | > 10s |
| Page Load (LCP) | < 2.5s | < 4s | > 4s |
| Time to Interactive | < 3.8s | < 7.3s | > 7.3s |
| Database Query | < 50ms | < 100ms | > 200ms |

## 12. Testing Checklist

### Before Committing
- [ ] All unit tests pass
- [ ] Test coverage meets threshold (80%)
- [ ] No skipped or commented tests
- [ ] Integration tests pass locally
- [ ] No console.log or debug statements
- [ ] Error cases are tested
- [ ] Edge cases are covered

### Before Merging
- [ ] CI/CD pipeline passes
- [ ] E2E tests pass
- [ ] Performance benchmarks met
- [ ] Security tests pass
- [ ] No flaky tests
- [ ] Documentation updated
- [ ] Breaking changes noted

### Before Release
- [ ] Full regression test suite
- [ ] Load testing completed
- [ ] Security audit performed
- [ ] Accessibility tests pass
- [ ] Cross-browser testing done
- [ ] Mobile testing completed
- [ ] Rollback plan tested

## Conclusion

This comprehensive testing strategy ensures:
1. **High Quality**: Bugs caught early through TDD
2. **Confidence**: Comprehensive coverage at all levels
3. **Performance**: Continuous monitoring of metrics
4. **Security**: Regular testing for vulnerabilities
5. **Maintainability**: Well-organized, readable tests

Remember: Tests are not just about catching bugs—they're about enabling confident refactoring, documenting behavior, and ensuring a great user experience.