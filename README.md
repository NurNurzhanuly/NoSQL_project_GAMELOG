# Online Game Marketplace

## Overview
This project is an online game marketplace that allows users to:
- Register and log in to their accounts.
- Browse a list of games available for purchase.
- Add games to their cart and proceed to checkout.
- View and manage their order history.

The application is built using:
- **Backend**: Golang (with Gin framework) and MongoDB for database management.
- **Frontend**: HTML, CSS, and JavaScript.
- **Deployment**: Docker for containerization.

---

## Features

### 1. User Authentication
- Register and log in using email and password.
- JWT-based authentication for secure access.

### 2. Game Listings
- View games with details like title, description, price, and images.
- Pagination support for large lists of games.

### 3. Cart Management
- Add or remove games from the cart.
- View cart contents before proceeding to checkout.

## Project Structure

```
marketplace/
├── backend/                # Backend code in Go
│   ├── main.go             # Entry point for the server
│   ├── routes/             # API route definitions
│   ├── controllers/        # Logic for handling API requests
│   ├── models/             # Database models
│   ├── services/           # Business logic and utilities
│   └── config/             # Configuration settings
│
├── frontend/               # Frontend code
│   ├── index.html          # Main page
│   ├── css/                # Stylesheets
│   │   └── styles.css
│   ├── js/                 # JavaScript logic
│   │   ├── app.js          # Main frontend logic
│   │   └── api.js          # API interaction logic
│   ├── images/             # Game images
│   └── assets/             # Additional assets (icons, fonts, etc.)
│
└── README.md               # Documentation
```

---

## Technologies Used

### Backend:
- **Golang**: High-performance backend development.
- **Gin Framework**: Lightweight web framework for Go.
- **MongoDB**: NoSQL database for efficient data storage and retrieval.

### Frontend:
- **HTML, CSS, JavaScript**: To create a responsive and interactive user interface.

### Deployment:
- **Docker**: Containerization for easy deployment and scalability.

---

## API Endpoints

### Authentication
- `POST /api/auth/register`: Register a new user.
- `POST /api/auth/login`: Log in with credentials.

### Games
- `GET /api/games`: Retrieve a list of available games.

### Cart
- `POST /api/cart`: Add a game to the cart.
- `DELETE /api/cart/:id`: Remove a game from the cart.

---

## License
This project is licensed under the [MIT License](LICENSE).
