# E-commerce Microservices Project

This project is an e-commerce platform built using a microservices architecture. The system is divided into several services, each responsible for a specific set of functionalities. It uses API Gateway to handle incoming requests and route them to the appropriate services.

## Architecture

The following services are part of this e-commerce platform:

- **API Gateway**: Handles all incoming requests and routes them to the correct microservices.
- **Admin Service**: Responsible for managing categories, products, user and seller accounts.
- **Seller Service**: Manages seller-specific functionalities such as product CRUD and offers.
- **User Service**: Handles user registration, login, and cart management.
- **Order Service**: Manages customer orders, order status, and details.
- **Payment Service**: Handles payment processing and integrations.

## Features

- **Admin**: 
  - Login/Logout
  - Category CRUD (Create, Read, Update, Delete)
  - Manage users and roles
- **Seller**:
  - Signup/Login
  - Product CRUD
  - Offers Management
- **User**:
  - Signup/Login/Logout
  - Cart Management (CRUD operations)
  - Order Creation
- **Order**:
  - Order Creation
  - Order Status Management
- **Payment**:
  - Payment Gateway Integration (e.g., Razorpay)

## Tech Stack

- **Backend**: Go (Gin Framework)
- **Microservices Communication**: gRPC
- **API Gateway**: Handles routing and authentication
- **Database**: PostgreSQL for persistent storage
- **Authentication**: JWT for secure API calls
- **Payment Gateway**: Razorpay or other integrations
