# Simple Auth Service

A minimal authentication service using Go, Gin, Gorm, PostgreSQL, and JWT.

## Features

- User registration (/register)

- Login (/login)

- Token refresh (/refresh)

- Logout (/logout)

- Get authenticated user (/me)

## Requirements

- Docker & Docker Compose

- Go (if running locally)

## Setup

```
cp .env.sample .env  # Set environment variables
docker-compose up --build
```

## Usage

### Register

```
curl -X POST http://localhost:8000/register \
    -H "Content-Type: application/json" \
    -d '{"email": "test@example.com", "password": "password123"}'
```

### Login

```
curl -X POST http://localhost:8000/login \
    -H "Content-Type: application/json" \
    -d '{"email": "test@example.com", "password": "password123"}' \
    -c cookies.txt
```

### Get User Info

```
curl -X GET http://localhost:8000/me -b cookies.txt
```
