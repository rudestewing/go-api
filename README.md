# GO-API

Simple Go REST API boilerplate

## 🚀 Quick Start

1. **Clone and install dependencies**

   ```bash
   git clone <repository-url>
   cd go-api
   go mod download
   ```

2. **Setup configuration**

   ```bash
   cp config.example.yaml config.yaml
   ```

   Edit `config.yaml` with your database and settings.

3. **Run migrations and start server**
   ```bash
   make migrate-up
   make dev
   ```

## 🛠️ What's Included

- **Authentication**: JWT-based auth with role management
- **Database**: GORM ORM with PostgreSQL, migrations & seeders
- **Middleware**: Rate limiting, CORS, validation, security headers
- **Email**: Template-based email system
- **Health Checks**: Monitoring endpoints
- **Hot Reload**: Development server with Air
- **Logging**: Structured logging with file rotation

## 🔧 Common Commands

```bash
make dev              # Start development server
make migrate-up       # Run database migrations
make migrate-create   # Create new migration
make seed-run         # Run database seeders
make test            # Run tests
```

## 📁 Project Structure

```
cmd/app/              # Main application entry point
cmd/migrate/          # Database migration tool
cmd/seed/             # Database seeding tool
config/               # Configuration management
database/             # Database connection and migrations
infrastructure/       # External services (email, etc.)
internal/             # Private application code
├── domain/           # Business domain logic
├── middleware/       # HTTP middleware
├── model/            # Data models
├── repository/       # Data access layer
├── router/           # Route definitions
└── shared/           # Utilities and helpers
storage/              # File storage and logs
```

Perfect for building REST APIs, microservices, or backend services with Go.
