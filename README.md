# Go API - RESTful API with JWT Authentication

A robust REST API built with Go, Fiber, GORM, and PostgreSQL. This project features JWT authentication, database migrations, and clean architecture patterns.

## 🚀 Features

### Authentication & Security

- **JWT Authentication** - Secure token-based authentication system
- **Password Hashing** - bcrypt password encryption
- **Protected Routes** - Middleware-based route protection
- **CORS Support** - Cross-origin resource sharing enabled

### Database Management

- **PostgreSQL Integration** - Production-ready database support
- **Database Migrations** - Version-controlled schema management
- **Connection Pooling** - Optimized database connections
- **GORM ORM** - Type-safe database operations

### API Features

- **User Registration** - Create new user accounts
- **User Login** - Authenticate existing users
- **User Profile** - Protected user profile endpoint
- **RESTful Design** - Standard HTTP methods and status codes

### Development Tools

- **Hot Reload** - Development server with automatic restart (Air)
- **Environment Configuration** - Flexible environment variable management
- **Makefile** - Simplified command execution
- **Error Handling** - Comprehensive error responses
- **Graceful Shutdown** - Clean server termination

## 📋 Tech Stack

- **Framework**: [Fiber](https://gofiber.io/) - Express-inspired web framework
- **Database**: PostgreSQL with [GORM](https://gorm.io/) ORM
- **Authentication**: JWT with [golang-jwt](https://github.com/golang-jwt/jwt)
- **Password**: bcrypt encryption
- **Hot Reload**: [Air](https://github.com/cosmtrek/air)
- **Environment**: [godotenv](https://github.com/joho/godotenv)

## 🏗️ Project Structure

```
go-api/
├── cmd/
│   └── migrate/           # Migration CLI tool
├── config/                # Configuration management
├── container/             # Dependency injection container
├── internal/
│   ├── handler/          # HTTP request handlers
│   ├── middleware/       # HTTP middlewares
│   ├── migration/        # Migration management
│   ├── model/           # Database models
│   ├── repository/      # Data access layer
│   └── service/         # Business logic layer
├── migrations/          # SQL migration files
├── router/             # Route definitions
├── tmp/               # Build artifacts
├── .env.example       # Environment variables template
├── air.toml          # Hot reload configuration
├── Makefile          # Build and development commands
└── main.go           # Application entry point
```

## 🛠️ Installation & Setup

### Prerequisites

- Go 1.24.3 or higher
- PostgreSQL 12+
- Git

### 1. Clone Repository

```bash
git clone <repository-url>
cd go-api
```

### 2. Install Dependencies

```bash
go mod download
```

### 3. Environment Configuration

```bash
# Copy environment template
cp .env.example .env

# Edit environment variables
nano .env
```

### Required Environment Variables

| Variable       | Description                  | Example                                                      |
| -------------- | ---------------------------- | ------------------------------------------------------------ |
| `DATABASE_URL` | PostgreSQL connection string | `postgres://user:pass@localhost:5432/dbname?sslmode=disable` |
| `JWT_SECRET`   | Secret key for JWT signing   | `your-super-secret-jwt-key-here`                             |

### Optional Environment Variables

| Variable            | Default       | Description             |
| ------------------- | ------------- | ----------------------- |
| `APP_PORT`          | `8000`        | Server port             |
| `JWT_EXPIRY`        | `24h`         | JWT token expiration    |
| `ENVIRONMENT`       | `development` | Environment mode        |
| `DB_MAX_IDLE_CONNS` | `10`          | Max idle DB connections |
| `DB_MAX_OPEN_CONNS` | `100`         | Max open DB connections |

### 4. Database Setup

```bash
# Create PostgreSQL database
createdb your_database_name

# Run migrations
make migrate-up
```

### 5. Install Development Tools (Optional)

```bash
# Install Air for hot reload
go install github.com/cosmtrek/air@latest
```

## 🚦 Running the Application

### Development Mode (with Hot Reload)

```bash
make dev
# or
air
```

### Production Mode

```bash
make run
# or
go run main.go
```

### Build Binary

```bash
make build
```

## 📊 Database Migrations

### Available Migration Commands

```bash
# Run all pending migrations
make migrate-up

# Rollback last migration batch
make migrate-down

# Check migration status
make migrate-status

# Create new migration
make migrate-create name="add_new_feature"
```

### Manual Migration Commands

```bash
# Using the migration CLI
go run cmd/migrate/main.go migrate     # Run migrations
go run cmd/migrate/main.go rollback    # Rollback migrations
go run cmd/migrate/main.go status      # Check status
go run cmd/migrate/main.go create "migration_name"  # Create new migration
```

## 🔗 API Endpoints

### Base URL

```
http://localhost:8000/api/v1
```

### Authentication Endpoints

| Method | Endpoint         | Description       | Auth Required |
| ------ | ---------------- | ----------------- | ------------- |
| `POST` | `/auth/register` | Register new user | ❌            |
| `POST` | `/auth/login`    | Login user        | ❌            |

### User Endpoints

| Method | Endpoint        | Description      | Auth Required |
| ------ | --------------- | ---------------- | ------------- |
| `GET`  | `/user/profile` | Get user profile | ✅            |

### Health Check

| Method | Endpoint | Description      | Auth Required |
| ------ | -------- | ---------------- | ------------- |
| `GET`  | `/`      | API health check | ❌            |

## 📝 API Usage Examples

### Register User

```bash
curl -X POST http://localhost:8000/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "password": "securepassword"
  }'
```

### Login User

```bash
curl -X POST http://localhost:8000/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "securepassword"
  }'
```

### Get User Profile (Protected)

```bash
curl -X GET http://localhost:8000/api/v1/user/profile \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## 🧪 Testing

```bash
# Run all tests
make test

# Run tests with coverage
go test -cover ./...

# Run specific test
go test ./internal/service/...
```

## 📦 Available Make Commands

```bash
make help              # Show all available commands
make run               # Run application
make build             # Build binary
make test              # Run tests
make clean             # Clean build artifacts
make dev               # Start development server with hot reload

# Migration commands
make migrate-up        # Run all pending migrations
make migrate-down      # Rollback last migration
make migrate-status    # Show migration status
make migrate-create    # Create new migration
```

## 🏛️ Architecture Patterns

### Clean Architecture

- **Handler Layer**: HTTP request/response handling
- **Service Layer**: Business logic implementation
- **Repository Layer**: Data access abstraction
- **Model Layer**: Data structures and entities

### Dependency Injection

- Container-based dependency management
- Loose coupling between components
- Easy testing and mocking

### Middleware Pipeline

- JWT authentication middleware
- CORS handling
- Request logging
- Error recovery

## 🔧 Configuration

The application uses a flexible configuration system with environment variables:

- **Required configs**: Database URL and JWT secret
- **Optional configs**: Server port, JWT expiry, connection limits
- **Validation**: Startup validation for required variables
- **Defaults**: Sensible defaults for optional settings

## 🛡️ Security Features

- **Password Hashing**: bcrypt with salt
- **JWT Tokens**: Secure token generation and validation
- **Environment Variables**: Sensitive data protection
- **CORS**: Cross-origin request handling
- **SQL Injection Protection**: GORM ORM protection

## 🚀 Deployment

### Docker (Coming Soon)

```bash
# Build image
docker build -t go-api .

# Run container
docker run -p 8000:8000 --env-file .env go-api
```

### Environment-Specific Configurations

- Development: Hot reload, debug logging
- Production: Optimized builds, connection pooling

## 🤝 Contributing

1. Fork the repository
2. Create feature branch (`git checkout -b feature/amazing-feature`)
3. Commit changes (`git commit -m 'Add amazing feature'`)
4. Push to branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📚 Additional Documentation

- [Environment Configuration](./ENV_CONFIG.md)
- [Migration Guide](./MIGRATIONS.md)
- [Migration Implementation](./MIGRATION_IMPLEMENTATION_SUMMARY.md)
- [Migration Quickstart](./MIGRATION_QUICKSTART.md)
- [Refactor Summary](./REFACTOR_SUMMARY.md)

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🆘 Troubleshooting

### Common Issues

1. **Database Connection Error**

   - Check PostgreSQL is running
   - Verify DATABASE_URL in .env file
   - Ensure database exists

2. **JWT Token Issues**

   - Verify JWT_SECRET is set
   - Check token format in Authorization header
   - Ensure token hasn't expired

3. **Migration Errors**

   - Check database permissions
   - Verify migration files syntax
   - Run `make migrate-status` to check state

4. **Hot Reload Not Working**
   - Install Air: `go install github.com/cosmtrek/air@latest`
   - Check air.toml configuration
   - Ensure tmp/ directory exists

### Getting Help

- Check existing documentation in the docs/ folder
- Review error logs in tmp/build-errors.log
- Ensure all environment variables are properly set
- Verify Go version compatibility (1.24.3+)
