# Setup Konfigurasi

## Setup untuk Development

1. **Copy template konfigurasi**:
```bash
cp config.example.yaml config.yaml
```

2. **Set environment variables yang required**:
```bash
export DATABASE_URL="postgres://username:password@localhost:5432/database_name"
export JWT_SECRET="your-super-secret-jwt-key-here-minimum-32-characters"
```

3. **Edit config.yaml sesuai kebutuhan** (opsional):
```bash
nano config.yaml
```

4. **Jalankan aplikasi**:
```bash
go run cmd/api/main.go
```

## Setup untuk Production

Untuk production, gunakan environment variables saja (jangan menggunakan config.yaml):

```bash
export DATABASE_URL="postgres://prod_user:prod_pass@prod_host:5432/prod_db"
export JWT_SECRET="super-secure-production-secret"
export GO_API_APP_ENVIRONMENT="production"
export GO_API_APP_PORT="8080"
# ... environment variables lainnya sesuai kebutuhan
```

## Prioritas Konfigurasi

1. **Environment Variables** (tertinggi) - `DATABASE_URL`, `JWT_SECRET`, `GO_API_*`
2. **config.yaml file** 
3. **Default values** (terendah)

## File yang Ignored dari Git

- `config.yaml` - Untuk menghindari commit secrets
- `.env` files - Legacy environment files

## File yang Di-commit ke Git

- `config.example.yaml` - Template untuk team
- `config/config.go` - Kode konfigurasi
