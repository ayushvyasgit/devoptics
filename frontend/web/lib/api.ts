import axios from "axios";

const API_BASE_URL = "http://localhost:8091/api/v1";

console.log("🔗 API Base URL:", API_BASE_URL);

const api = axios.create({
  baseURL: API_BASE_URL,
  headers: { "Content-Type": "application/json" },
  timeout: 30000,
});

// Request interceptor
api.interceptors.request.use(
  (config) => {
    console.log(`📤 ${config.method?.toUpperCase()} ${config.url}`);
    return config;
  },
  (error) => {
    console.error("❌ Request Error:", error);
    return Promise.reject(error);
  }
);

// Response interceptor
api.interceptors.response.use(
  (response) => {
    console.log(`✅ ${response.status} ${response.config.url}`);
    return response;
  },
  (error) => {
    console.error("❌ Response Error:", error.response?.status, error.response?.data);
    return Promise.reject(error);
  }
);

// Types
export interface RegisterRequest {
  email: string;
  password: string;
  first_name: string;
  last_name: string;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface User {
  id: string;
  email: string;
  first_name: string;
  last_name: string;
  role: string;
}

export interface LoginResponse {
  token: string;
  refresh_token: string;
  user: User;
}

// Auth API
export const authAPI = {
  register: (data: RegisterRequest) => {
    console.log("🔥 Registering:", data.email);
    return api.post("/auth/register", data);
  },
  
  login: (data: LoginRequest) => {
    console.log("🔑 Logging in:", data.email);
    return api.post<LoginResponse>("/auth/login", data);
  },
};

export default api;