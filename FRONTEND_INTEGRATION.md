# Frontend Integration Guide — Web Streaming Backend

Panduan lengkap integrasi frontend dengan backend, termasuk contoh kode React/Vue, error handling, dan best practices.

---

## 📦 Setup Environment

### 1. Backend Configuration
Pastikan backend sudah running dengan env yang tepat:

```bash
cd backend
cp .env.example .env

# .env harus berisi:
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=web_streaming
DB_SSLMODE=disable
JWT_SECRET=your-very-long-secret-key-min-32-chars
ALLOWED_ORIGINS=http://localhost:3000,http://localhost:5173,https://yourdomain.com
ACCESS_TOKEN_DURATION=15m
REFRESH_TOKEN_DURATION=168h

go run main.go
```

### 2. Frontend Configuration
```bash
# Frontend project (.env or .env.local)
VITE_API_URL=http://localhost:1010/api/v1  # Development
# VITE_API_URL=https://api.yourdomain.com/api/v1  # Production
```

---

## 🔐 Authentication Pattern

### Option A: Bearer Token (RECOMMENDED)

**Backend Response Format:**
```json
{
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "550e8400-e29b-41d4-a716-446655440000",
    "expires_in": 900
  },
  "error": null
}
```

**Frontend: Store & Use Tokens**
```typescript
// store.ts (Pinia/Vuex example)
import { defineStore } from 'pinia';

export const useAuthStore = defineStore('auth', () => {
  const accessToken = ref<string | null>(null);
  const refreshToken = ref<string | null>(null);
  const user = ref<any>(null);

  const setTokens = (access: string, refresh: string) => {
    accessToken.value = access;
    refreshToken.value = refresh;
    localStorage.setItem('access_token', access);
    localStorage.setItem('refresh_token', refresh);
  };

  const clearTokens = () => {
    accessToken.value = null;
    refreshToken.value = null;
    localStorage.removeItem('access_token');
    localStorage.removeItem('refresh_token');
  };

  const loadTokensFromStorage = () => {
    const stored_access = localStorage.getItem('access_token');
    const stored_refresh = localStorage.getItem('refresh_token');
    if (stored_access) {
      accessToken.value = stored_access;
      refreshToken.value = stored_refresh;
    }
  };

  return { accessToken, refreshToken, user, setTokens, clearTokens, loadTokensFromStorage };
});
```

### Option B: HttpOnly Cookie (SECURE but CORS-Complex)

**Backend Sets Cookie:**
```go
c.SetCookie("access_token", accessToken, 900, "/", "", true, true)  // HttpOnly+Secure
```

**Frontend Usage:**
```typescript
// Cookies sent automatically with credentials:true
fetch(apiUrl, {
  method: 'POST',
  credentials: 'include',  // CRITICAL for cookies
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({ email, password })
})
```

**⚠️ Important:** CORS must include `Access-Control-Allow-Credentials: true` (backend already does this).

---

## 🚀 API Client Implementation

### React Example (with Axios)
```typescript
// api/client.ts
import axios from 'axios';
import { useAuthStore } from '@/stores/auth';

const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:1010/api/v1';

export const apiClient = axios.create({
  baseURL: API_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Request interceptor: add auth token
apiClient.interceptors.request.use((config) => {
  const authStore = useAuthStore();
  if (authStore.accessToken) {
    config.headers.Authorization = `Bearer ${authStore.accessToken}`;
  }
  return config;
});

// Response interceptor: handle 401 & refresh token
apiClient.interceptors.response.use(
  (response) => response,
  async (error) => {
    const originalRequest = error.config;

    if (error.response?.status === 401 && !originalRequest._retry) {
      originalRequest._retry = true;
      const authStore = useAuthStore();

      try {
        const response = await axios.post(`${API_URL}/refresh-token`, {
          refresh_token: authStore.refreshToken,
        });

        const { access_token, refresh_token } = response.data.data;
        authStore.setTokens(access_token, refresh_token);

        // Retry original request
        originalRequest.headers.Authorization = `Bearer ${access_token}`;
        return apiClient(originalRequest);
      } catch (refreshError) {
        authStore.clearTokens();
        window.location.href = '/login';
        return Promise.reject(refreshError);
      }
    }

    return Promise.reject(error);
  }
);

export default apiClient;
```

### Vue Composable Example
```typescript
// composables/useApi.ts
import { ref } from 'vue';
import apiClient from '@/api/client';
import type { AxiosError } from 'axios';

interface ApiResponse<T> {
  data: T | null;
  error: string | null;
  isLoading: boolean;
}

export function useApi<T = any>(url: string, method: 'get' | 'post' | 'put' | 'delete' = 'get') {
  const response = ref<ApiResponse<T>>({
    data: null,
    error: null,
    isLoading: false,
  });

  const execute = async (body?: any) => {
    response.value.isLoading = true;
    try {
      const result = await apiClient[method](url, method !== 'get' ? body : undefined);
      response.value.data = result.data.data;
      response.value.error = result.data.error;
    } catch (err) {
      const axiosError = err as AxiosError<any>;
      response.value.error = axiosError.response?.data?.error || 'An error occurred';
    } finally {
      response.value.isLoading = false;
    }
  };

  return { ...response, execute };
}
```

---

## 📱 Component Examples

### 1. Login Component (Vue 3)
```vue
<template>
  <div class="login-container">
    <form @submit.prevent="handleLogin">
      <div class="form-group">
        <label for="email">Email:</label>
        <input
          id="email"
          v-model="email"
          type="email"
          placeholder="Enter your email"
          required
        />
      </div>

      <div class="form-group">
        <label for="password">Password:</label>
        <input
          id="password"
          v-model="password"
          type="password"
          placeholder="Enter your password"
          required
        />
      </div>

      <button type="submit" :disabled="isLoading">
        {{ isLoading ? 'Logging in...' : 'Login' }}
      </button>

      <div v-if="error" class="error-message">
        {{ error }}
      </div>
    </form>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import { useAuthStore } from '@/stores/auth';
import apiClient from '@/api/client';

const email = ref('');
const password = ref('');
const isLoading = ref(false);
const error = ref('');
const router = useRouter();
const authStore = useAuthStore();

const handleLogin = async () => {
  isLoading.value = true;
  error.value = '';

  try {
    const response = await apiClient.post('/login', { email: email.value, password: password.value });

    if (response.data.error) {
      error.value = response.data.error;
    } else {
      const { access_token, refresh_token } = response.data.data;
      authStore.setTokens(access_token, refresh_token);
      router.push('/home');
    }
  } catch (err: any) {
    error.value = err.response?.data?.error || 'Login failed. Please try again.';
  } finally {
    isLoading.value = false;
  }
};
</script>

<style scoped>
.login-container {
  max-width: 400px;
  margin: 50px auto;
  padding: 20px;
  border: 1px solid #ccc;
  border-radius: 8px;
}

.form-group {
  margin-bottom: 15px;
}

input {
  width: 100%;
  padding: 10px;
  border: 1px solid #ccc;
  border-radius: 4px;
}

button {
  width: 100%;
  padding: 10px;
  background-color: #007bff;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}

button:disabled {
  background-color: #ccc;
  cursor: not-allowed;
}

.error-message {
  color: red;
  margin-top: 10px;
}
</style>
```

### 2. Film List Component (React)
```typescript
// components/FilmList.tsx
import React, { useEffect, useState } from 'react';
import apiClient from '@/api/client';

interface Film {
  id: number;
  title: string;
  description: string;
  genre: string;
  year: number;
  poster_url: string;
  rating: number;
}

interface PaginatedResponse {
  items: Film[];
  page: number;
  limit: number;
  total: number;
}

export function FilmList() {
  const [films, setFilms] = useState<Film[]>([]);
  const [loading, setLoading] = useState(true);
  const [page, setPage] = useState(1);
  const [totalPages, setTotalPages] = useState(1);
  const [error, setError] = useState('');

  useEffect(() => {
    fetchFilms(page);
  }, [page]);

  const fetchFilms = async (pageNum: number) => {
    setLoading(true);
    setError('');
    try {
      const response = await apiClient.get(`/films?page=${pageNum}&limit=10`);
      if (response.data.error) {
        setError(response.data.error);
      } else {
        const data = response.data.data as PaginatedResponse;
        setFilms(data.items);
        setTotalPages(Math.ceil(data.total / data.limit));
      }
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to load films');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="film-list">
      <h1>Films</h1>

      {error && <div className="error-banner">{error}</div>}

      {loading ? (
        <div className="loading">Loading films...</div>
      ) : (
        <>
          <div className="films-grid">
            {films.map((film) => (
              <div key={film.id} className="film-card">
                <img src={film.poster_url} alt={film.title} />
                <h3>{film.title}</h3>
                <p className="year">{film.year}</p>
                <p className="genre">{film.genre}</p>
                <p className="rating">⭐ {film.rating}</p>
                <p className="description">{film.description}</p>
                <button>Add to Watchlist</button>
              </div>
            ))}
          </div>

          <div className="pagination">
            <button onClick={() => setPage(page - 1)} disabled={page === 1}>
              Previous
            </button>
            <span>
              Page {page} of {totalPages}
            </span>
            <button onClick={() => setPage(page + 1)} disabled={page === totalPages}>
              Next
            </button>
          </div>
        </>
      )}
    </div>
  );
}
```

### 3. Search Component
```typescript
// components/FilmSearch.tsx
import React, { useState, useRef } from 'react';
import apiClient from '@/api/client';

interface Film {
  id: number;
  title: string;
  rating: number;
  poster_url: string;
}

export function FilmSearch() {
  const [searchQuery, setSearchQuery] = useState('');
  const [results, setResults] = useState<Film[]>([]);
  const [loading, setLoading] = useState(false);
  const debounceTimer = useRef<NodeJS.Timeout>();

  const handleSearch = (query: string) => {
    setSearchQuery(query);

    // Debounce API call (500ms)
    if (debounceTimer.current) {
      clearTimeout(debounceTimer.current);
    }

    if (!query.trim()) {
      setResults([]);
      return;
    }

    setLoading(true);
    debounceTimer.current = setTimeout(async () => {
      try {
        const response = await apiClient.get(`/films/search?title=${encodeURIComponent(query)}`);
        if (response.data.data?.films) {
          setResults(response.data.data.films);
        }
      } catch (err) {
        console.error('Search failed:', err);
        setResults([]);
      } finally {
        setLoading(false);
      }
    }, 500);
  };

  return (
    <div className="search-container">
      <input
        type="text"
        placeholder="Search films..."
        value={searchQuery}
        onChange={(e) => handleSearch(e.target.value)}
      />

      {loading && <div className="loading">Searching...</div>}

      {results.length > 0 && (
        <div className="search-results">
          {results.map((film) => (
            <div key={film.id} className="search-result-item">
              <img src={film.poster_url} alt={film.title} width="50" />
              <div>
                <h4>{film.title}</h4>
                <p>Rating: {film.rating}</p>
              </div>
              <button>Add to Watchlist</button>
            </div>
          ))}
        </div>
      )}

      {searchQuery && results.length === 0 && !loading && (
        <div className="no-results">No films found</div>
      )}
    </div>
  );
}
```

### 4. Watchlist Component
```typescript
// components/Watchlist.tsx
import React, { useEffect, useState } from 'react';
import apiClient from '@/api/client';

interface WatchlistItem {
  id: number;
  film_id: number;
  user_id: number;
  added_at: string;
}

export function Watchlist() {
  const [watchlist, setWatchlist] = useState<WatchlistItem[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    fetchWatchlist();
  }, []);

  const fetchWatchlist = async () => {
    try {
      const response = await apiClient.get('/watchlist');
      if (response.data.data) {
        setWatchlist(response.data.data);
      }
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to load watchlist');
    } finally {
      setLoading(false);
    }
  };

  const removeFromWatchlist = async (id: number) => {
    try {
      await apiClient.delete(`/watchlist/${id}`);
      setWatchlist(watchlist.filter((item) => item.id !== id));
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to remove from watchlist');
    }
  };

  return (
    <div className="watchlist">
      <h2>My Watchlist</h2>

      {error && <div className="error-banner">{error}</div>}

      {loading ? (
        <div className="loading">Loading watchlist...</div>
      ) : watchlist.length === 0 ? (
        <p>Your watchlist is empty</p>
      ) : (
        <ul>
          {watchlist.map((item) => (
            <li key={item.id}>
              <span>Film ID: {item.film_id}</span>
              <span>Added: {new Date(item.added_at).toLocaleDateString()}</span>
              <button onClick={() => removeFromWatchlist(item.id)}>Remove</button>
            </li>
          ))}
        </ul>
      )}
    </div>
  );
}
```

---

## 🛡️ Error Handling Best Practices

### Global Error Handler (React)
```typescript
// hooks/useAsyncError.tsx
import React, { useEffect } from 'react';

export function ErrorBoundary({ children }: { children: React.ReactNode }) {
  const [error, setError] = React.useState<string | null>(null);

  useEffect(() => {
    const handleError = (event: ErrorEvent) => {
      setError(event.message);
      setTimeout(() => setError(null), 5000);
    };

    window.addEventListener('error', handleError);
    return () => window.removeEventListener('error', handleError);
  }, []);

  return (
    <>
      {error && <div className="global-error">{error}</div>}
      {children}
    </>
  );
}
```

### API Error Status Mapping
```typescript
export const getErrorMessage = (status: number, defaultMsg: string) => {
  const messages: Record<number, string> = {
    400: 'Invalid request. Please check your input.',
    401: 'Session expired. Please login again.',
    403: 'You do not have permission to perform this action.',
    404: 'Resource not found.',
    409: 'This resource already exists.',
    429: 'Too many requests. Please try again later.',
    500: 'Server error. Please try again later.',
  };

  return messages[status] || defaultMsg;
};
```

---

## 🔄 Token Refresh Strategy

### Automatic Refresh on Expiry
```typescript
// middleware/authInterceptor.ts
const setupTokenRefresh = (authStore) => {
  // Check token expiry every minute
  setInterval(() => {
    const token = authStore.accessToken;
    if (!token) return;

    const decoded = jwtDecode(token);
    const expiresIn = decoded.exp * 1000 - Date.now();

    // Refresh if less than 5 minutes left
    if (expiresIn < 5 * 60 * 1000) {
      refreshAccessToken(authStore);
    }
  }, 60 * 1000); // Check every minute
};

const refreshAccessToken = async (authStore) => {
  try {
    const response = await axios.post(`${API_URL}/refresh-token`, {
      refresh_token: authStore.refreshToken,
    });
    authStore.setTokens(response.data.data.access_token, response.data.data.refresh_token);
  } catch (err) {
    authStore.clearTokens();
    window.location.href = '/login';
  }
};
```

---

## 📋 Testing API Endpoints

### Postman Collection JSON
```json
{
  "info": {
    "name": "Web Streaming API",
    "schema": "https://schema.getpostman.com/json/collection/v2.1.0/collection.json"
  },
  "item": [
    {
      "name": "Auth",
      "item": [
        {
          "name": "Register",
          "request": {
            "method": "POST",
            "header": [{"key": "Content-Type", "value": "application/json"}],
            "body": {
              "mode": "raw",
              "raw": "{\"username\":\"testuser\",\"email\":\"test@example.com\",\"password\":\"Pass123\"}"
            },
            "url": {"raw": "{{base_url}}/register", "host": ["{{base_url}}"], "path": ["register"]}
          }
        },
        {
          "name": "Login",
          "request": {
            "method": "POST",
            "header": [{"key": "Content-Type", "value": "application/json"}],
            "body": {
              "mode": "raw",
              "raw": "{\"email\":\"test@example.com\",\"password\":\"Pass123\"}"
            },
            "url": {"raw": "{{base_url}}/login", "host": ["{{base_url}}"], "path": ["login"]}
          }
        }
      ]
    },
    {
      "name": "Films",
      "item": [
        {
          "name": "Get All Films",
          "request": {
            "method": "GET",
            "url": {"raw": "{{base_url}}/films?page=1&limit=10", "host": ["{{base_url}}"], "path": ["films"], "query": [{"key": "page", "value": "1"}, {"key": "limit", "value": "10"}]}
          }
        }
      ]
    }
  ],
  "variable": [
    {"key": "base_url", "value": "http://localhost:1010/api/v1"},
    {"key": "access_token", "value": ""}
  ]
}
```

---

## ✅ Frontend Checklist

- [ ] Environment variables configured
- [ ] API client with interceptors set up
- [ ] Auth store/state management implemented
- [ ] Login & register pages created
- [ ] Protected routes implemented
- [ ] Film list & search pages working
- [ ] Watchlist functionality working
- [ ] Error handling & loading states
- [ ] Token refresh logic working
- [ ] Tested on both dev & staging environments
- [ ] CORS working (no origin errors)
- [ ] API responses parsed correctly
- [ ] Error messages displayed properly

