# Database Migrations

Sistem database migrations untuk Go API yang memungkinkan Anda mengelola perubahan database secara manual dan terkontrol.

## Fitur

- ✅ Migrations manual (tidak otomatis)
- ✅ Tracking migrations yang sudah dijalankan
- ✅ Rollback migrations
- ✅ Status migrations
- ✅ Template otomatis untuk migration files
- ✅ Command berbasis Go (tanpa shell scripts)

## Struktur Migration File

Setiap migration file menggunakan format berikut:

```sql
-- +migrate Up
-- SQL statements for applying the migration
CREATE TABLE example (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL
);

-- +migrate Down
-- SQL statements for rolling back the migration
DROP TABLE IF EXISTS example;
```

## Commands

### 1. Membuat Migration Baru

```bash
# Menggunakan go run
go run cmd/migrate/main.go create "create_users_table"

# Menggunakan make (lebih mudah)
make migrate-create name="create_users_table"
```

Ini akan membuat file migration dengan timestamp di folder `migrations/`:

```
migrations/20250527120345_create_users_table.sql
```

### 2. Menjalankan Migrations

```bash
# Menggunakan go run
go run cmd/migrate/main.go migrate

# Menggunakan make
make migrate-up
```

Command ini akan:

- Menjalankan semua migrations yang belum dieksekusi
- Mencatat migrations yang sudah dijalankan dalam tabel `migrations`
- Mengelompokkan migrations dalam batch untuk rollback

### 3. Rollback Migrations

```bash
# Menggunakan go run
go run cmd/migrate/main.go rollback

# Menggunakan make
make migrate-down
```

Command ini akan:

- Rollback batch terakhir dari migrations
- Menjalankan SQL dari section `-- +migrate Down`
- Menghapus record migration dari database

### 4. Status Migrations

```bash
# Menggunakan go run
go run cmd/migrate/main.go status

# Menggunakan make
make migrate-status
```

Menampilkan:

- ✅ Migrations yang sudah dijalankan
- ⏳ Migrations yang menunggu untuk dijalankan

### 5. Help

```bash
go run cmd/migrate/main.go help
```

## Contoh Workflow

### 1. Membuat migration untuk tabel users

```bash
make migrate-create name="create_users_table"
```

Edit file yang dibuat di `migrations/` dengan SQL yang sesuai:

```sql
-- +migrate Up
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_email ON users(email);

-- +migrate Down
DROP INDEX IF EXISTS idx_users_email;
DROP TABLE IF EXISTS users;
```

### 2. Jalankan migration

```bash
make migrate-up
```

### 3. Cek status

```bash
make migrate-status
```

Output:

```
Migration Status:
================

Executed migrations:
  ✓ 20250527120345_create_users_table

No pending migrations
```

### 4. Jika perlu rollback

```bash
make migrate-down
```

## Best Practices

### 1. Naming Convention

- Gunakan nama yang deskriptif: `create_users_table`, `add_index_to_users`, `alter_users_add_phone`
- Gunakan snake_case
- Sertakan action dan target: `create_`, `add_`, `alter_`, `drop_`

### 2. Migration Content

- Selalu sertakan section `-- +migrate Up` dan `-- +migrate Down`
- Pastikan rollback bisa mengembalikan state sebelumnya
- Gunakan `IF EXISTS` atau `IF NOT EXISTS` untuk safety
- Test migration di development environment terlebih dahulu

### 3. Index Management

```sql
-- +migrate Up
CREATE INDEX CONCURRENTLY idx_users_email ON users(email);

-- +migrate Down
DROP INDEX IF EXISTS idx_users_email;
```

### 4. Data Migration

```sql
-- +migrate Up
ALTER TABLE users ADD COLUMN phone VARCHAR(20);
UPDATE users SET phone = '' WHERE phone IS NULL;

-- +migrate Down
ALTER TABLE users DROP COLUMN phone;
```

## Struktur Database

Sistem migration menggunakan tabel `migrations` untuk tracking:

```sql
CREATE TABLE migrations (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL,
    batch INTEGER NOT NULL,
    executed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

- `name`: Nama file migration (tanpa .sql)
- `batch`: Batch number untuk rollback
- `executed_at`: Waktu eksekusi

## Environment Variables

Pastikan `.env` file sudah dikonfigurasi dengan:

```env
DATABASE_URL=postgres://username:password@localhost/database_name?sslmode=disable
```

## Troubleshooting

### Error: "DATABASE_URL environment variable is required"

- Pastikan file `.env` ada dan berisi `DATABASE_URL`
- Jalankan dari root directory project

### Migration file tidak terdeteksi

- Pastikan file ada di folder `migrations/`
- Pastikan format file sesuai dengan template
- Pastikan ada section `-- +migrate Up` dan `-- +migrate Down`

### Rollback gagal

- Pastikan SQL di section `-- +migrate Down` valid
- Check apakah ada dependencies yang menghalangi (foreign keys, dll)

## Development

### Testing Migrations

1. Test di development database terlebih dahulu
2. Jalankan migration: `make migrate-up`
3. Cek hasil di database
4. Test rollback: `make migrate-down`
5. Pastikan rollback mengembalikan state awal

### Production Deployment

1. Backup database sebelum migration
2. Jalankan migration di staging environment dulu
3. Monitor log output saat migration
4. Verifikasi hasil setelah migration selesai
