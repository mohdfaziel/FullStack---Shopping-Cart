// API Base URL - uses environment variable or fallback to localhost
const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8081';

// Helper function to get token from localStorage
const getToken = () => localStorage.getItem('token');

// Helper function to make API requests
const apiRequest = async (endpoint, options = {}) => {
  const url = `${API_BASE_URL}${endpoint}`;
  const token = getToken();
  
  const config = {
    headers: {
      'Content-Type': 'application/json',
      ...(token && { 'Authorization': `Bearer ${token}` }),
      ...options.headers,
    },
    ...options,
  };

  try {
    const response = await fetch(url, config);
    
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }
    
    return await response.json();
  } catch (error) {
    console.error('API request failed:', error);
    throw error;
  }
};

// Auth API functions
export const authAPI = {
  // Create user (signup)
  signup: (userData) => apiRequest('/users', {
    method: 'POST',
    body: JSON.stringify(userData),
  }),

  // Login user
  login: (credentials) => apiRequest('/users/login', {
    method: 'POST',
    body: JSON.stringify(credentials),
  }),
};

// Items API functions
export const itemsAPI = {
  // Get all items
  getAll: () => apiRequest('/items'),

  // Create item (admin function)
  create: (itemData) => apiRequest('/items', {
    method: 'POST',
    body: JSON.stringify(itemData),
  }),
};

// Cart API functions
export const cartAPI = {
  // Get user's cart
  get: () => apiRequest('/carts'),

  // Add item to cart
  addItem: (itemId) => apiRequest('/carts', {
    method: 'POST',
    body: JSON.stringify({ item_id: itemId }),
  }),

  // Remove item from cart
  removeItem: (itemId) => apiRequest(`/carts/${itemId}`, {
    method: 'DELETE',
  }),
};

// Orders API functions
export const ordersAPI = {
  // Create order (checkout)
  create: () => apiRequest('/orders', {
    method: 'POST',
  }),

  // Get user's orders
  getAll: () => apiRequest('/orders'),
};

// Authentication helpers
export const auth = {
  setToken: (token) => localStorage.setItem('token', token),
  getToken,
  removeToken: () => localStorage.removeItem('token'),
  isAuthenticated: () => !!getToken(),
};
