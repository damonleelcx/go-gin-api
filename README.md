# Go Gin API

RESTful API project based on the Gin framework, providing user authentication functionality.

## Requirements

- Go 1.24.0 or higher

## Getting Started

### 1. Install Dependencies

```bash
go mod download
```

Or use `go mod tidy` to organize dependencies:

```bash
go mod tidy
```

### 2. Run the Project

Start the server using the following command:

```bash
go run server.go
```

Or:

```bash
go run .
```

### 3. Verify Service

The server will start at `http://localhost:8080`. You can:

- Access root path: `http://localhost:8080/` - should return `{"message":"Hello, World!"}`
- API route prefix: `/api`

## Available API Endpoints

The project provides the following authentication-related API endpoints:

- `POST /api/signup` - User registration
- `POST /api/login` - User login
- `POST /api/logout` - User logout
- `POST /api/refresh` - Refresh token
- `POST /api/forgot-password` - Forgot password
- `POST /api/reset-password` - Reset password

## Database

The project uses SQLite in-memory database (`:memory:`), which means data will be lost after server restart. The database will automatically migrate the following tables:

- User (User table)
- Session (Session table)
- PasswordResetToken (Password reset token table)

## Build Executable

If you want to compile to an executable file:

```bash
go build -o go-gin-api.exe server.go
```

Then run:

```bash
./go-gin-api.exe
```

## Reference Resources

GraphQL: https://gqlgen.com/recipes/gin/

Testing:
https://github.com/onsi/ginkgo

https://github.com/onsi/gomega

https://github.com/onsi/biloba
