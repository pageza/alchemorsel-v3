# Frontend Architecture & Component Design

## 1. Directory Structure

```
alchemorsel-frontend/
├── src/
│   ├── assets/
│   │   ├── images/              # Static images
│   │   ├── styles/              # Global styles
│   │   │   ├── main.css        # Main stylesheet
│   │   │   └── variables.css   # CSS variables
│   │   └── fonts/              # Custom fonts
│   ├── components/
│   │   ├── common/             # Shared components
│   │   │   ├── AppHeader.vue
│   │   │   ├── AppFooter.vue
│   │   │   ├── LoadingSpinner.vue
│   │   │   ├── ErrorAlert.vue
│   │   │   └── ConfirmDialog.vue
│   │   ├── auth/               # Authentication components
│   │   │   ├── LoginForm.vue
│   │   │   ├── RegisterForm.vue
│   │   │   └── PasswordReset.vue
│   │   ├── recipe/             # Recipe components
│   │   │   ├── RecipeCard.vue
│   │   │   ├── RecipeList.vue
│   │   │   ├── RecipeDetail.vue
│   │   │   ├── RecipeForm.vue
│   │   │   ├── RecipeSearch.vue
│   │   │   └── RecipeGenerator.vue
│   │   └── user/               # User components
│   │       ├── ProfileForm.vue
│   │       ├── ProfilePicture.vue
│   │       └── PreferencesEditor.vue
│   ├── composables/            # Vue 3 composables
│   │   ├── useAuth.ts
│   │   ├── useRecipes.ts
│   │   ├── useUser.ts
│   │   ├── useNotification.ts
│   │   └── useDebounce.ts
│   ├── layouts/                # Layout components
│   │   ├── DefaultLayout.vue
│   │   ├── AuthLayout.vue
│   │   └── DashboardLayout.vue
│   ├── pages/                  # Page components
│   │   ├── Home.vue
│   │   ├── Login.vue
│   │   ├── Register.vue
│   │   ├── Dashboard.vue
│   │   ├── RecipeSearch.vue
│   │   ├── RecipeDetail.vue
│   │   ├── RecipeGenerate.vue
│   │   ├── Profile.vue
│   │   └── NotFound.vue
│   ├── router/                 # Vue Router
│   │   ├── index.ts
│   │   ├── routes.ts
│   │   └── guards.ts
│   ├── services/               # API services
│   │   ├── api.ts             # Axios instance
│   │   ├── auth.service.ts
│   │   ├── recipe.service.ts
│   │   ├── user.service.ts
│   │   └── storage.service.ts
│   ├── stores/                 # Pinia stores
│   │   ├── auth.store.ts
│   │   ├── recipe.store.ts
│   │   ├── user.store.ts
│   │   └── notification.store.ts
│   ├── types/                  # TypeScript types
│   │   ├── api.types.ts
│   │   ├── auth.types.ts
│   │   ├── recipe.types.ts
│   │   ├── user.types.ts
│   │   └── common.types.ts
│   ├── utils/                  # Utility functions
│   │   ├── validators.ts
│   │   ├── formatters.ts
│   │   ├── constants.ts
│   │   └── helpers.ts
│   ├── App.vue                # Root component
│   ├── main.ts                # Application entry
│   └── env.d.ts              # Environment types
├── public/
│   ├── index.html
│   ├── favicon.ico
│   └── robots.txt
├── tests/
│   ├── unit/                  # Unit tests
│   ├── component/             # Component tests
│   └── e2e/                   # E2E tests
├── .env.example               # Environment template
├── .eslintrc.js              # ESLint config
├── .prettierrc               # Prettier config
├── index.html                # Entry HTML
├── package.json              # Dependencies
├── tailwind.config.js        # Tailwind config
├── tsconfig.json             # TypeScript config
├── vite.config.ts            # Vite config
└── vitest.config.ts          # Vitest config
```

## 2. Core Components Design

### 2.1 Type Definitions

```typescript
// src/types/user.types.ts
export interface User {
  id: string;
  email: string;
  username: string;
  name: string;
  profilePictureUrl?: string;
  dietaryPreferences: string[];
  allergies: string[];
  createdAt: string;
  updatedAt: string;
}

export interface RegisterRequest {
  email: string;
  username: string;
  password: string;
  name: string;
  dietaryPreferences: string[];
  allergies: string[];
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface AuthTokens {
  accessToken: string;
  refreshToken: string;
}
```

```typescript
// src/types/recipe.types.ts
export interface Recipe {
  id: string;
  userId: string;
  title: string;
  description: string;
  ingredients: Ingredient[];
  instructions: string[];
  prepTime: number;
  cookTime: number;
  servings: number;
  category: string;
  dietaryCategories: string[];
  allergens: string[];
  nutritionalInfo: NutritionalInfo;
  imageUrl?: string;
  isPublic: boolean;
  isFavorited?: boolean;
  createdAt: string;
  updatedAt: string;
}

export interface Ingredient {
  name: string;
  amount: number;
  unit: string;
  optional: boolean;
}

export interface NutritionalInfo {
  calories: number;
  protein: number;
  carbohydrates: number;
  fat: number;
  fiber: number;
  sugar: number;
  sodium: number;
}

export interface RecipeSearchParams {
  q?: string;
  category?: string;
  dietary?: string[];
  exclude?: string[];
  page?: number;
  limit?: number;
}

export interface RecipeGenerateRequest {
  style?: string;
  ingredients?: string[];
  cookingTime?: number;
  servings?: number;
  customPrompt?: string;
}
```

### 2.2 API Service Layer

```typescript
// src/services/api.ts
import axios, { AxiosInstance, AxiosError } from 'axios';
import { useAuthStore } from '@/stores/auth.store';
import { useNotificationStore } from '@/stores/notification.store';

class ApiService {
  private instance: AxiosInstance;
  
  constructor() {
    this.instance = axios.create({
      baseURL: import.meta.env.VITE_API_URL || 'http://localhost:8080/api/v1',
      timeout: 30000,
      headers: {
        'Content-Type': 'application/json',
      },
    });
    
    this.setupInterceptors();
  }
  
  private setupInterceptors(): void {
    // Request interceptor
    this.instance.interceptors.request.use(
      (config) => {
        const authStore = useAuthStore();
        if (authStore.accessToken) {
          config.headers.Authorization = `Bearer ${authStore.accessToken}`;
        }
        return config;
      },
      (error) => Promise.reject(error)
    );
    
    // Response interceptor
    this.instance.interceptors.response.use(
      (response) => response,
      async (error: AxiosError) => {
        const authStore = useAuthStore();
        const notificationStore = useNotificationStore();
        
        if (error.response?.status === 401) {
          // Try to refresh token
          try {
            await authStore.refreshToken();
            // Retry original request
            return this.instance.request(error.config!);
          } catch {
            authStore.logout();
            notificationStore.error('Session expired. Please login again.');
          }
        }
        
        return Promise.reject(error);
      }
    );
  }
  
  get api(): AxiosInstance {
    return this.instance;
  }
}

export default new ApiService().api;
```

```typescript
// src/services/recipe.service.ts
import api from './api';
import type { Recipe, RecipeSearchParams, RecipeGenerateRequest } from '@/types/recipe.types';

export class RecipeService {
  static async search(params: RecipeSearchParams): Promise<{
    recipes: Recipe[];
    total: number;
    page: number;
    totalPages: number;
  }> {
    const response = await api.get('/recipes', { params });
    return response.data;
  }
  
  static async getById(id: string): Promise<Recipe> {
    const response = await api.get(`/recipes/${id}`);
    return response.data;
  }
  
  static async generate(request: RecipeGenerateRequest): Promise<Recipe> {
    const response = await api.post('/llm/generate', request);
    return response.data;
  }
  
  static async addFavorite(recipeId: string): Promise<void> {
    await api.post(`/recipes/${recipeId}/favorite`);
  }
  
  static async removeFavorite(recipeId: string): Promise<void> {
    await api.delete(`/recipes/${recipeId}/favorite`);
  }
  
  static async getFavorites(page = 1, limit = 20): Promise<Recipe[]> {
    const response = await api.get('/recipes/favorites', {
      params: { page, limit }
    });
    return response.data.recipes;
  }
}
```

### 2.3 Pinia Store Layer

```typescript
// src/stores/auth.store.ts
import { defineStore } from 'pinia';
import { ref, computed } from 'vue';
import type { User, LoginRequest, RegisterRequest, AuthTokens } from '@/types';
import { AuthService } from '@/services/auth.service';
import { StorageService } from '@/services/storage.service';

export const useAuthStore = defineStore('auth', () => {
  // State
  const user = ref<User | null>(null);
  const accessToken = ref<string | null>(null);
  const refreshToken = ref<string | null>(null);
  const isLoading = ref(false);
  
  // Getters
  const isAuthenticated = computed(() => !!accessToken.value);
  const currentUser = computed(() => user.value);
  
  // Actions
  async function login(credentials: LoginRequest): Promise<void> {
    isLoading.value = true;
    try {
      const response = await AuthService.login(credentials);
      setTokens(response.tokens);
      user.value = response.user;
    } finally {
      isLoading.value = false;
    }
  }
  
  async function register(data: RegisterRequest): Promise<void> {
    isLoading.value = true;
    try {
      const response = await AuthService.register(data);
      setTokens(response.tokens);
      user.value = response.user;
    } finally {
      isLoading.value = false;
    }
  }
  
  async function logout(): Promise<void> {
    try {
      await AuthService.logout(refreshToken.value!);
    } finally {
      clearAuth();
    }
  }
  
  async function refreshTokens(): Promise<void> {
    if (!refreshToken.value) throw new Error('No refresh token');
    
    const response = await AuthService.refreshToken(refreshToken.value);
    setTokens(response.tokens);
  }
  
  function setTokens(tokens: AuthTokens): void {
    accessToken.value = tokens.accessToken;
    refreshToken.value = tokens.refreshToken;
    StorageService.setTokens(tokens);
  }
  
  function clearAuth(): void {
    user.value = null;
    accessToken.value = null;
    refreshToken.value = null;
    StorageService.clearTokens();
  }
  
  // Initialize from storage
  function initialize(): void {
    const tokens = StorageService.getTokens();
    if (tokens) {
      accessToken.value = tokens.accessToken;
      refreshToken.value = tokens.refreshToken;
      // Fetch user profile
      AuthService.getProfile().then(profile => {
        user.value = profile;
      }).catch(() => {
        clearAuth();
      });
    }
  }
  
  return {
    // State
    user,
    isLoading,
    // Getters
    isAuthenticated,
    currentUser,
    // Actions
    login,
    register,
    logout,
    refreshTokens,
    initialize,
  };
});
```

### 2.4 Composables

```typescript
// src/composables/useRecipes.ts
import { ref, computed, watch } from 'vue';
import { useRecipeStore } from '@/stores/recipe.store';
import { useNotificationStore } from '@/stores/notification.store';
import type { RecipeSearchParams } from '@/types';

export function useRecipes() {
  const recipeStore = useRecipeStore();
  const notificationStore = useNotificationStore();
  
  const searchParams = ref<RecipeSearchParams>({
    q: '',
    category: undefined,
    dietary: [],
    exclude: [],
    page: 1,
    limit: 20,
  });
  
  const isLoading = ref(false);
  const error = ref<string | null>(null);
  
  const recipes = computed(() => recipeStore.recipes);
  const totalRecipes = computed(() => recipeStore.total);
  const currentPage = computed(() => recipeStore.currentPage);
  const totalPages = computed(() => recipeStore.totalPages);
  
  async function searchRecipes(): Promise<void> {
    isLoading.value = true;
    error.value = null;
    
    try {
      await recipeStore.searchRecipes(searchParams.value);
    } catch (err) {
      error.value = 'Failed to search recipes';
      notificationStore.error(error.value);
    } finally {
      isLoading.value = false;
    }
  }
  
  async function loadMore(): Promise<void> {
    if (currentPage.value >= totalPages.value) return;
    
    searchParams.value.page = currentPage.value + 1;
    await searchRecipes();
  }
  
  async function toggleFavorite(recipeId: string): Promise<void> {
    try {
      await recipeStore.toggleFavorite(recipeId);
      notificationStore.success('Favorite updated');
    } catch {
      notificationStore.error('Failed to update favorite');
    }
  }
  
  // Auto-search when params change
  watch(searchParams, () => {
    searchRecipes();
  }, { deep: true });
  
  return {
    searchParams,
    recipes,
    totalRecipes,
    currentPage,
    totalPages,
    isLoading,
    error,
    searchRecipes,
    loadMore,
    toggleFavorite,
  };
}
```

```typescript
// src/composables/useDebounce.ts
import { ref, watch, type Ref } from 'vue';

export function useDebounce<T>(value: Ref<T>, delay = 300) {
  const debouncedValue = ref<T>(value.value) as Ref<T>;
  let timeout: NodeJS.Timeout;
  
  watch(value, (newValue) => {
    clearTimeout(timeout);
    timeout = setTimeout(() => {
      debouncedValue.value = newValue;
    }, delay);
  });
  
  return debouncedValue;
}
```

### 2.5 Component Examples

```vue
<!-- src/components/recipe/RecipeSearch.vue -->
<template>
  <div class="recipe-search">
    <div class="search-header">
      <h2 class="text-2xl font-bold mb-4">Search Recipes</h2>
      
      <div class="search-controls">
        <input
          v-model="searchQuery"
          type="text"
          placeholder="Search recipes..."
          class="search-input"
          @input="handleSearch"
        />
        
        <select v-model="selectedCategory" class="category-select">
          <option value="">All Categories</option>
          <option v-for="cat in categories" :key="cat" :value="cat">
            {{ cat }}
          </option>
        </select>
        
        <div class="dietary-filters">
          <label v-for="diet in dietaryOptions" :key="diet" class="checkbox-label">
            <input
              type="checkbox"
              :value="diet"
              v-model="selectedDietary"
            />
            {{ diet }}
          </label>
        </div>
      </div>
    </div>
    
    <div v-if="isLoading" class="loading-state">
      <LoadingSpinner />
    </div>
    
    <div v-else-if="error" class="error-state">
      <ErrorAlert :message="error" />
    </div>
    
    <div v-else class="search-results">
      <div class="recipe-grid">
        <RecipeCard
          v-for="recipe in recipes"
          :key="recipe.id"
          :recipe="recipe"
          @toggle-favorite="toggleFavorite"
        />
      </div>
      
      <div v-if="hasMore" class="load-more">
        <button
          @click="loadMore"
          :disabled="isLoadingMore"
          class="btn btn-secondary"
        >
          {{ isLoadingMore ? 'Loading...' : 'Load More' }}
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';
import { useDebounce } from '@/composables/useDebounce';
import { useRecipes } from '@/composables/useRecipes';
import RecipeCard from './RecipeCard.vue';
import LoadingSpinner from '@/components/common/LoadingSpinner.vue';
import ErrorAlert from '@/components/common/ErrorAlert.vue';

const {
  searchParams,
  recipes,
  totalPages,
  currentPage,
  isLoading,
  error,
  searchRecipes,
  loadMore: loadMoreRecipes,
  toggleFavorite,
} = useRecipes();

const searchQuery = ref('');
const selectedCategory = ref('');
const selectedDietary = ref<string[]>([]);
const isLoadingMore = ref(false);

const debouncedSearchQuery = useDebounce(searchQuery, 500);

const categories = [
  'Breakfast',
  'Lunch',
  'Dinner',
  'Dessert',
  'Snack',
  'Appetizer',
];

const dietaryOptions = [
  'Vegan',
  'Vegetarian',
  'Gluten-Free',
  'Dairy-Free',
  'Keto',
  'Paleo',
];

const hasMore = computed(() => currentPage.value < totalPages.value);

function handleSearch(): void {
  searchParams.value.q = debouncedSearchQuery.value;
  searchParams.value.category = selectedCategory.value;
  searchParams.value.dietary = selectedDietary.value;
  searchParams.value.page = 1;
}

async function loadMore(): Promise<void> {
  isLoadingMore.value = true;
  try {
    await loadMoreRecipes();
  } finally {
    isLoadingMore.value = false;
  }
}
</script>

<style scoped>
.recipe-search {
  @apply container mx-auto px-4 py-8;
}

.search-controls {
  @apply space-y-4 mb-8;
}

.search-input {
  @apply w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent;
}

.category-select {
  @apply w-full md:w-auto px-4 py-2 border border-gray-300 rounded-lg;
}

.dietary-filters {
  @apply flex flex-wrap gap-4;
}

.checkbox-label {
  @apply flex items-center space-x-2;
}

.recipe-grid {
  @apply grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6;
}

.load-more {
  @apply text-center mt-8;
}

.loading-state,
.error-state {
  @apply flex justify-center items-center min-h-[400px];
}
</style>
```

### 2.6 Router Configuration

```typescript
// src/router/routes.ts
import type { RouteRecordRaw } from 'vue-router';

export const routes: RouteRecordRaw[] = [
  {
    path: '/',
    component: () => import('@/layouts/DefaultLayout.vue'),
    children: [
      {
        path: '',
        name: 'home',
        component: () => import('@/pages/Home.vue'),
      },
      {
        path: 'recipes',
        name: 'recipes',
        component: () => import('@/pages/RecipeSearch.vue'),
      },
      {
        path: 'recipes/:id',
        name: 'recipe-detail',
        component: () => import('@/pages/RecipeDetail.vue'),
      },
    ],
  },
  {
    path: '/auth',
    component: () => import('@/layouts/AuthLayout.vue'),
    children: [
      {
        path: 'login',
        name: 'login',
        component: () => import('@/pages/Login.vue'),
        meta: { requiresGuest: true },
      },
      {
        path: 'register',
        name: 'register',
        component: () => import('@/pages/Register.vue'),
        meta: { requiresGuest: true },
      },
    ],
  },
  {
    path: '/dashboard',
    component: () => import('@/layouts/DashboardLayout.vue'),
    meta: { requiresAuth: true },
    children: [
      {
        path: '',
        name: 'dashboard',
        component: () => import('@/pages/Dashboard.vue'),
      },
      {
        path: 'profile',
        name: 'profile',
        component: () => import('@/pages/Profile.vue'),
      },
      {
        path: 'generate',
        name: 'generate-recipe',
        component: () => import('@/pages/RecipeGenerate.vue'),
      },
    ],
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'not-found',
    component: () => import('@/pages/NotFound.vue'),
  },
];
```

```typescript
// src/router/guards.ts
import type { NavigationGuardNext, RouteLocationNormalized } from 'vue-router';
import { useAuthStore } from '@/stores/auth.store';

export function authGuard(
  to: RouteLocationNormalized,
  from: RouteLocationNormalized,
  next: NavigationGuardNext
): void {
  const authStore = useAuthStore();
  
  if (to.meta.requiresAuth && !authStore.isAuthenticated) {
    next({ name: 'login', query: { redirect: to.fullPath } });
  } else if (to.meta.requiresGuest && authStore.isAuthenticated) {
    next({ name: 'dashboard' });
  } else {
    next();
  }
}
```

## 3. Testing Strategy

### 3.1 Unit Tests

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
  });
  
  it('should login user successfully', async () => {
    const mockUser = { id: '1', email: 'test@example.com' };
    const mockTokens = { accessToken: 'access', refreshToken: 'refresh' };
    
    vi.mocked(AuthService.login).mockResolvedValue({
      user: mockUser,
      tokens: mockTokens,
    });
    
    const store = useAuthStore();
    await store.login({ email: 'test@example.com', password: 'password' });
    
    expect(store.isAuthenticated).toBe(true);
    expect(store.currentUser).toEqual(mockUser);
  });
});
```

### 3.2 Component Tests

```typescript
// tests/component/RecipeCard.spec.ts
import { mount } from '@vue/test-utils';
import { describe, it, expect } from 'vitest';
import RecipeCard from '@/components/recipe/RecipeCard.vue';

describe('RecipeCard', () => {
  const mockRecipe = {
    id: '1',
    title: 'Test Recipe',
    description: 'A test recipe',
    prepTime: 10,
    cookTime: 20,
    servings: 4,
    imageUrl: 'test.jpg',
    isFavorited: false,
  };
  
  it('renders recipe information correctly', () => {
    const wrapper = mount(RecipeCard, {
      props: { recipe: mockRecipe },
    });
    
    expect(wrapper.text()).toContain('Test Recipe');
    expect(wrapper.text()).toContain('A test recipe');
    expect(wrapper.text()).toContain('30 min');
    expect(wrapper.text()).toContain('4 servings');
  });
  
  it('emits toggle-favorite event when heart icon clicked', async () => {
    const wrapper = mount(RecipeCard, {
      props: { recipe: mockRecipe },
    });
    
    await wrapper.find('[data-test="favorite-button"]').trigger('click');
    
    expect(wrapper.emitted('toggle-favorite')).toBeTruthy();
    expect(wrapper.emitted('toggle-favorite')[0]).toEqual(['1']);
  });
});
```

## 4. Performance Optimization

### 4.1 Route-based Code Splitting
- Lazy load routes for optimal initial bundle size
- Prefetch critical routes
- Use Suspense for loading states

### 4.2 Image Optimization
- Lazy load images with Intersection Observer
- Use WebP format with fallbacks
- Implement responsive images
- CDN integration for static assets

### 4.3 State Management
- Normalize data structure in stores
- Implement proper caching strategies
- Use computed properties for derived state
- Avoid unnecessary reactivity

### 4.4 API Optimization
- Implement request debouncing
- Use pagination for large datasets
- Cache responses with proper invalidation
- Implement optimistic updates

## 5. Accessibility

### 5.1 ARIA Implementation
- Proper semantic HTML
- ARIA labels and descriptions
- Keyboard navigation support
- Screen reader compatibility

### 5.2 Focus Management
- Visible focus indicators
- Focus trap for modals
- Skip navigation links
- Proper tab order

### 5.3 Color & Contrast
- WCAG AA compliance
- High contrast mode support
- Color blind friendly palette
- Proper text sizing