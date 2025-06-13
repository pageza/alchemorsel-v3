# API Design & Endpoint Specifications

## 1. API Design Principles

### RESTful Standards
- **Resource-based URLs**: Use nouns, not verbs
- **HTTP Methods**: GET (read), POST (create), PUT (update), PATCH (partial update), DELETE (remove)
- **Status Codes**: Use appropriate HTTP status codes
- **Versioning**: URL path versioning (`/api/v1/`)
- **Consistency**: Uniform response structure across all endpoints

### Response Format
```json
{
  "success": true,
  "data": {},
  "meta": {
    "timestamp": "2024-01-01T00:00:00Z",
    "version": "1.0.0"
  }
}
```

### Error Response Format
```json
{
  "success": false,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Invalid input data",
    "details": {
      "field": "email",
      "reason": "Invalid email format"
    }
  },
  "meta": {
    "timestamp": "2024-01-01T00:00:00Z",
    "request_id": "uuid"
  }
}
```

### Pagination Format
```json
{
  "success": true,
  "data": [],
  "pagination": {
    "page": 1,
    "per_page": 20,
    "total": 100,
    "total_pages": 5,
    "has_next": true,
    "has_prev": false
  }
}
```

## 2. Authentication Endpoints

### POST /api/v1/auth/register
Register a new user account.

**Request Body:**
```json
{
  "email": "user@example.com",
  "username": "johndoe",
  "password": "SecurePassword123!",
  "name": "John Doe",
  "dietary_preferences": ["vegetarian", "gluten-free"],
  "allergies": ["peanuts", "shellfish"]
}
```

**Validation Rules:**
- Email: Valid email format, unique
- Username: 3-30 characters, alphanumeric + underscore, unique
- Password: Min 8 chars, 1 uppercase, 1 lowercase, 1 number, 1 special char
- Name: 2-100 characters
- Dietary preferences: At least one required
- Allergies: Optional array

**Success Response (201 Created):**
```json
{
  "success": true,
  "data": {
    "user": {
      "id": "uuid",
      "email": "user@example.com",
      "username": "johndoe",
      "name": "John Doe",
      "dietary_preferences": ["vegetarian", "gluten-free"],
      "allergies": ["peanuts", "shellfish"],
      "created_at": "2024-01-01T00:00:00Z"
    },
    "tokens": {
      "access_token": "jwt_access_token",
      "refresh_token": "jwt_refresh_token",
      "expires_in": 900
    }
  }
}
```

### POST /api/v1/auth/login
Authenticate user and receive tokens.

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "SecurePassword123!"
}
```

**Success Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "user": {
      "id": "uuid",
      "email": "user@example.com",
      "username": "johndoe",
      "name": "John Doe"
    },
    "tokens": {
      "access_token": "jwt_access_token",
      "refresh_token": "jwt_refresh_token",
      "expires_in": 900
    }
  }
}
```

### POST /api/v1/auth/refresh
Refresh access token using refresh token.

**Request Body:**
```json
{
  "refresh_token": "jwt_refresh_token"
}
```

**Success Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "tokens": {
      "access_token": "new_jwt_access_token",
      "refresh_token": "new_jwt_refresh_token",
      "expires_in": 900
    }
  }
}
```

### POST /api/v1/auth/logout
Invalidate refresh token.

**Headers:**
```
Authorization: Bearer <access_token>
```

**Request Body:**
```json
{
  "refresh_token": "jwt_refresh_token"
}
```

**Success Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "message": "Successfully logged out"
  }
}
```

### POST /api/v1/auth/forgot-password
Initiate password reset process.

**Request Body:**
```json
{
  "email": "user@example.com"
}
```

**Success Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "message": "Password reset instructions sent to email"
  }
}
```

### POST /api/v1/auth/reset-password
Reset password with token.

**Request Body:**
```json
{
  "token": "reset_token_from_email",
  "password": "NewSecurePassword123!"
}
```

## 3. User Endpoints

### GET /api/v1/users/profile
Get current user profile.

**Headers:**
```
Authorization: Bearer <access_token>
```

**Success Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "user": {
      "id": "uuid",
      "email": "user@example.com",
      "username": "johndoe",
      "name": "John Doe",
      "profile_picture_url": "https://s3.amazonaws.com/...",
      "dietary_preferences": ["vegetarian", "gluten-free"],
      "allergies": ["peanuts", "shellfish"],
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    }
  }
}
```

### PUT /api/v1/users/profile
Update user profile.

**Headers:**
```
Authorization: Bearer <access_token>
```

**Request Body:**
```json
{
  "name": "John Updated",
  "dietary_preferences": ["vegan"],
  "allergies": ["dairy"]
}
```

### POST /api/v1/users/profile/picture
Upload profile picture.

**Headers:**
```
Authorization: Bearer <access_token>
Content-Type: multipart/form-data
```

**Form Data:**
- `picture`: Image file (JPEG, PNG, WebP)
- Max size: 5MB
- Dimensions: Min 100x100, Max 2000x2000

**Success Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "profile_picture_url": "https://s3.amazonaws.com/..."
  }
}
```

### DELETE /api/v1/users/profile
Delete user account (soft delete).

**Headers:**
```
Authorization: Bearer <access_token>
```

**Request Body:**
```json
{
  "password": "CurrentPassword123!"
}
```

## 4. Recipe Endpoints

### GET /api/v1/recipes
Search and list recipes.

**Headers:**
```
Authorization: Bearer <access_token> (optional)
```

**Query Parameters:**
- `q`: Search query (searches title, description, ingredients)
- `category`: Filter by category (breakfast, lunch, dinner, etc.)
- `dietary`: Comma-separated dietary preferences
- `exclude`: Comma-separated allergens to exclude
- `user_id`: Filter by recipe creator
- `favorites`: Boolean, show only favorites (requires auth)
- `page`: Page number (default: 1)
- `per_page`: Items per page (default: 20, max: 100)
- `sort`: Sort field (created_at, title, prep_time)
- `order`: Sort order (asc, desc)

**Success Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "recipes": [
      {
        "id": "uuid",
        "title": "Vegan Buddha Bowl",
        "description": "A nutritious and colorful bowl",
        "category": "lunch",
        "dietary_categories": ["vegan", "gluten-free"],
        "prep_time": 15,
        "cook_time": 20,
        "servings": 2,
        "image_url": "https://...",
        "is_favorited": false,
        "created_at": "2024-01-01T00:00:00Z"
      }
    ]
  },
  "pagination": {
    "page": 1,
    "per_page": 20,
    "total": 100,
    "total_pages": 5
  }
}
```

### GET /api/v1/recipes/:id
Get single recipe details.

**Headers:**
```
Authorization: Bearer <access_token> (optional)
```

**Success Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "recipe": {
      "id": "uuid",
      "user_id": "uuid",
      "title": "Vegan Buddha Bowl",
      "description": "A nutritious and colorful bowl packed with vegetables",
      "ingredients": [
        {
          "name": "Quinoa",
          "amount": 1,
          "unit": "cup",
          "optional": false
        },
        {
          "name": "Chickpeas",
          "amount": 1,
          "unit": "can",
          "optional": false
        }
      ],
      "instructions": [
        "Cook quinoa according to package instructions",
        "Drain and rinse chickpeas",
        "Prepare vegetables"
      ],
      "prep_time": 15,
      "cook_time": 20,
      "servings": 2,
      "category": "lunch",
      "dietary_categories": ["vegan", "gluten-free"],
      "allergens": [],
      "nutritional_info": {
        "calories": 450,
        "protein": 15,
        "carbohydrates": 65,
        "fat": 12,
        "fiber": 10,
        "sugar": 8,
        "sodium": 300
      },
      "image_url": "https://...",
      "is_public": true,
      "is_favorited": false,
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    }
  }
}
```

### POST /api/v1/recipes
Create a new recipe.

**Headers:**
```
Authorization: Bearer <access_token>
```

**Request Body:**
```json
{
  "title": "My Special Recipe",
  "description": "A family favorite",
  "ingredients": [
    {
      "name": "Flour",
      "amount": 2,
      "unit": "cups",
      "optional": false
    }
  ],
  "instructions": [
    "Mix ingredients",
    "Bake at 350°F"
  ],
  "prep_time": 20,
  "cook_time": 30,
  "servings": 4,
  "category": "dessert",
  "dietary_categories": ["vegetarian"],
  "is_public": true
}
```

### PUT /api/v1/recipes/:id
Update a recipe (must be owner).

**Headers:**
```
Authorization: Bearer <access_token>
```

### DELETE /api/v1/recipes/:id
Delete a recipe (must be owner).

**Headers:**
```
Authorization: Bearer <access_token>
```

### POST /api/v1/recipes/:id/favorite
Add recipe to favorites.

**Headers:**
```
Authorization: Bearer <access_token>
```

**Success Response (201 Created):**
```json
{
  "success": true,
  "data": {
    "message": "Recipe added to favorites"
  }
}
```

### DELETE /api/v1/recipes/:id/favorite
Remove recipe from favorites.

**Headers:**
```
Authorization: Bearer <access_token>
```

## 5. LLM Generation Endpoints

### POST /api/v1/llm/generate
Generate a recipe using AI.

**Headers:**
```
Authorization: Bearer <access_token>
```

**Request Body:**
```json
{
  "style": "italian",
  "ingredients": ["tomatoes", "basil", "mozzarella"],
  "cooking_time": 30,
  "servings": 4,
  "difficulty": "easy",
  "custom_prompt": "Make it suitable for a romantic dinner"
}
```

**Success Response (201 Created):**
```json
{
  "success": true,
  "data": {
    "recipe": {
      "id": "uuid",
      "title": "Romantic Caprese Pasta",
      "description": "A simple yet elegant Italian dish",
      "ingredients": [...],
      "instructions": [...],
      "prep_time": 10,
      "cook_time": 20,
      "servings": 4,
      "category": "dinner",
      "dietary_categories": ["vegetarian"],
      "nutritional_info": {...}
    }
  }
}
```

### POST /api/v1/llm/suggest
Get recipe suggestions based on preferences.

**Headers:**
```
Authorization: Bearer <access_token>
```

**Request Body:**
```json
{
  "meal_type": "lunch",
  "max_time": 30,
  "avoid_ingredients": ["dairy"]
}
```

**Success Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "suggestions": [
      {
        "title": "Quick Vegan Stir-Fry",
        "description": "A fast and healthy lunch option",
        "estimated_time": 25,
        "difficulty": "easy"
      }
    ]
  }
}
```

## 6. Health & Monitoring Endpoints

### GET /api/v1/health
Basic health check.

**Success Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "status": "healthy",
    "timestamp": "2024-01-01T00:00:00Z"
  }
}
```

### GET /api/v1/health/detailed
Detailed health check (internal use).

**Success Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "status": "healthy",
    "services": {
      "database": "connected",
      "redis": "connected",
      "s3": "accessible",
      "deepseek": "responsive"
    },
    "version": "1.0.0",
    "uptime": 3600
  }
}
```

## 7. Rate Limiting

### Rate Limit Headers
All responses include rate limit information:
```
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 99
X-RateLimit-Reset: 1640995200
```

### Rate Limits by Endpoint Type
- **Authentication**: 5 requests per minute
- **Recipe Generation**: 10 requests per hour
- **General API**: 100 requests per minute
- **File Upload**: 10 requests per hour

### Rate Limit Exceeded Response (429)
```json
{
  "success": false,
  "error": {
    "code": "RATE_LIMIT_EXCEEDED",
    "message": "Too many requests",
    "retry_after": 60
  }
}
```

## 8. WebSocket Events (Future Enhancement)

### Connection
```
wss://api.alchemorsel.com/ws
```

### Events
- `recipe:created` - New public recipe created
- `recipe:updated` - Recipe updated
- `user:online` - User came online
- `notification` - User-specific notifications

## 9. API Security

### CORS Configuration
```
Access-Control-Allow-Origin: https://alchemorsel.com
Access-Control-Allow-Methods: GET, POST, PUT, DELETE, OPTIONS
Access-Control-Allow-Headers: Content-Type, Authorization
Access-Control-Max-Age: 86400
```

### Security Headers
```
X-Content-Type-Options: nosniff
X-Frame-Options: DENY
X-XSS-Protection: 1; mode=block
Content-Security-Policy: default-src 'self'
```

### Request Validation
- All inputs sanitized and validated
- SQL injection prevention via parameterized queries
- XSS prevention through proper encoding
- CSRF protection via double-submit cookies