# Configuration Setup

## Setup untuk Development

1. **Copy template konfigurasi**:

```bash
cp config.example.yaml config.yaml
```

2. **Edit config.yaml dengan values yang sesuai**:

```bash
nano config.yaml
```

Pastikan untuk mengubah minimal:

- `database.url`: URL PostgreSQL database Anda
- `security.jwt_secret`: Secret key untuk JWT (minimal 32 karakter)

3. **Jalankan aplikasi**:

```bash
go run cmd/api/main.go
```

## Konfigurasi yang Diperlukan

### Database

- URL: `postgres://username:password@localhost:5432/database_name`
- Connection pooling: max idle/open connections, connection lifetime

### Security

- JWT Secret: Minimal 32 karakter untuk production
- JWT Expiry: Default 24 jam
- Trusted proxies: Untuk production behind load balancer

### Application

- Port: Default 8000
- Environment: development/production
- Timeouts: read, write, idle, shutdown

### CORS

- Allowed origins: Frontend URL
- Allowed methods: HTTP methods yang diizinkan
- Allowed headers: Headers yang diizinkan

### Rate Limiting

- Max requests per window
- Time window duration
- Enable/disable rate limiting

### Logging

- Log directory: Default storage/logs
- Max file size: Default 10MB
- Log retention: Default 30 hari
- Daily rotation: Default enabled

## File yang Ignored dari Git

- `config.yaml` - Berisi secrets dan konfigurasi lokal
- File ini tidak akan di-commit ke repository untuk keamanan

## Troubleshooting

### Config file not found

```
Config file 'config.yaml' not found. Please copy from config.example.yaml
```

**Solusi**: Copy `config.example.yaml` ke `config.yaml` dan edit sesuai kebutuhan

### Required configurations missing

```
Required configurations need to be updated in config.yaml: [database.url security.jwt_secret]
```

**Solusi**: Edit `config.yaml` dan ubah placeholder values dengan values yang sebenarnya

### Database connection failed

```
failed to connect to database
```

**Solusi**:

1. Pastikan PostgreSQL berjalan
2. Check database URL di `config.yaml`
3. Pastikan database dan user sudah dibuat

## File yang Di-commit ke Git

- `config.example.yaml` - Template untuk team
- `config/config.go` - Kode konfigurasi
