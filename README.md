# Shopping Cart E-commerce Project

**Developer:** Mohd Faziel  
A full-stack e-commerce application built with **Golang (Gin)** backend and **React.js** frontend using **daisyUI** and **TailwindCSS**.

## 🏗️ Project Structure

```
Fullstack - Shopping Cart/
├── Backend/                 # Golang backend
│   ├── main.go             # Entry point
│   ├── go.mod              # Go modules
│   ├── models/             # Database models
│   ├── handlers/           # API handlers
│   ├── middleware/         # Authentication middleware
│   └── database/           # Database connection
└── Frontend/               # React frontend
    ├── src/
    │   ├── components/     # React components
    │   ├── pages/          # Page components
    │   ├── utils/          # API utilities
    │   └── App.jsx         # Main app component
    ├── package.json
    └── vite.config.js
```

## 🛠️ Tech Stack

### Backend
- **Framework**: Gin (Golang web framework)
- **ORM**: GORM (Go Object Relational Mapping)
- **Database**: SQLite (for development)
- **Authentication**: JWT tokens
- **Testing**: Ginkgo (can be added)

### Frontend
- **Framework**: React.js with Vite
- **Styling**: TailwindCSS v4 + daisyUI
- **Routing**: React Router DOM
- **HTTP Client**: Fetch API

## 📊 Database Schema

The application uses the following entities:

- **users** (id, username, password, token, cart_id, created_at)
- **items** (id, name, status, created_at)
- **carts** (id, user_id, name, status, created_at)
- **cart_items** (cart_id, item_id) - Many-to-many relationship
- **orders** (id, cart_id, user_id, created_at)

## 🚀 Getting Started

### Prerequisites

1. **Go** (version 1.21 or higher)
   - Download from: https://golang.org/dl/
   
2. **Node.js** (version 16 or higher)
   - Download from: https://nodejs.org/

### Backend Setup

1. **Navigate to Backend directory:**
   ```bash
   cd "Backend"
   ```

2. **Install Go dependencies:**
   ```bash
   go mod tidy
   ```

3. **Run the server:**
   ```bash
   go run main.go
   ```

The backend will start on `http://localhost:8080`

### Frontend Setup

1. **Navigate to Frontend directory:**
   ```bash
   cd "Frontend"
   ```

2. **Install dependencies:**
   ```bash
   npm install
   ```

3. **Start the development server:**
   ```bash
   npm run dev
   ```

The frontend will start on `http://localhost:5174` (or another port if 5174 is busy)

## 📡 API Endpoints

| Method | URL            | Description                                | Auth Required |
| ------ | -------------- | ------------------------------------------ | ------------- |
| POST   | `/users`       | Create a user                              | No            |
| GET    | `/users`       | List all users                             | No            |
| POST   | `/users/login` | Login user                                 | No            |
| POST   | `/items`       | Create item                                | No            |
| GET    | `/items`       | List items                                 | No            |
| POST   | `/carts`       | Add items to cart                          | Yes           |
| GET    | `/carts`       | Get user's cart                            | Yes           |
| POST   | `/orders`      | Convert cart to order (checkout)           | Yes           |
| GET    | `/orders`      | List user's orders                         | Yes           |

## 🔐 Authentication

- Users sign up and receive a unique token upon login
- Only one active session per user (single token)
- Token required for cart and order operations
- Tokens are stored in localStorage on the frontend

## 🎯 Features

### Frontend Features
- **Login Screen**: Username/password authentication
- **Items Listing**: Display all available items
- **Shopping Cart**: Add items to cart
- **Checkout**: Convert cart to order
- **Order History**: View past orders
- **Responsive Design**: Mobile-friendly with daisyUI

### Backend Features
- **User Management**: Registration and authentication
- **Item Management**: CRUD operations for items
- **Cart System**: One active cart per user
- **Order Processing**: Convert carts to orders
- **JWT Authentication**: Secure API access

## 🧪 Testing the API

### Using curl

1. **Create a user:**
   ```bash
   curl -X POST http://localhost:8080/users \
     -H "Content-Type: application/json" \
     -d '{"username": "testuser", "password": "password123"}'
   ```

2. **Login:**
   ```bash
   curl -X POST http://localhost:8080/users/login \
     -H "Content-Type: application/json" \
     -d '{"username": "testuser", "password": "password123"}'
   ```

3. **Create an item:**
   ```bash
   curl -X POST http://localhost:8080/items \
     -H "Content-Type: application/json" \
     -d '{"name": "Sample Item", "status": "available"}'
   ```

4. **Add item to cart (requires token):**
   ```bash
   curl -X POST http://localhost:8080/carts \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer YOUR_JWT_TOKEN" \
     -d '{"item_id": 1}'
   ```

5. **Checkout (create order):**
   ```bash
   curl -X POST http://localhost:8080/orders \
     -H "Authorization: Bearer YOUR_JWT_TOKEN"
   ```

### Using the Frontend

1. **Access the application**: http://localhost:5174
2. **Create an account** using the signup button
3. **Login** with your credentials
4. **Browse items** and add them to your cart
5. **Use the action buttons**: Checkout, Cart, Order History

## 🛡️ Business Logic

- Users must be authenticated to access cart/order routes
- Each user can only have one active cart at a time
- When checkout occurs, the cart is converted into an order
- Cart status changes from "active" to "ordered"
- Users can create a new cart after checkout

## 🔧 Development Notes

- **Database**: SQLite file (`shopping_cart.db`) is created automatically
- **CORS**: Enabled for frontend-backend communication
- **Hot Reload**: Both frontend (Vite) and backend support hot reload
- **Styling**: All UI components use daisyUI classes exclusively

## 📝 Future Enhancements

- Add item quantities to cart
- Implement user profiles
- Add item categories and search
- Payment integration
- Admin panel for item management
- Email notifications
- Unit tests with Ginkgo

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test thoroughly
5. Submit a pull request

## 📄 License

This project is licensed under the MIT License.
