# Go API - Environment Configuration

## Required Environment Variables

The following environment variables **MUST** be set or the application will fail to start:

### 🔴 Required Variables

| Variable       | Description                      | Example                                                      |
| -------------- | -------------------------------- | ------------------------------------------------------------ |
| `DATABASE_URL` | PostgreSQL connection string     | `postgres://user:pass@localhost:5432/dbname?sslmode=disable` |
| `JWT_SECRET`   | Secret key for JWT token signing | `your-super-secret-jwt-key-here`                             |

### 🟡 Optional Variables

| Variable   | Description | Default | Example |
| ---------- | ----------- | ------- | ------- |
| `APP_PORT` | Server port | `8000`  | `3000`  |

## Setup Instructions

1. **Copy environment template:**

   ```bash
   cp .env.example .env
   ```

2. **Edit the .env file:**

   ```bash
   nano .env
   ```

3. **Set required variables:**

   ```env
   DATABASE_URL=postgres://username:password@localhost:5432/database_name?sslmode=disable
   JWT_SECRET=your-super-secret-jwt-key-here
   APP_PORT=8000
   ```

4. **Run the application:**
   ```bash
   go run main.go
   ```

## Error Messages

If required environment variables are missing, you'll see:

```
Required environment variables are missing: [DATABASE_URL JWT_SECRET]
```

Make sure to set all required variables before starting the application.

## Security Notes

- **Never commit `.env` files** to version control
- Use strong, randomly generated JWT secrets in production
- Use secure database credentials
- Consider using environment-specific configurations for different deployment stages
