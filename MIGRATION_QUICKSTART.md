# Quick Start - Database Migrations

Sistem database migrations untuk Go API sudah berhasil disetup dan siap digunakan!

## ✅ Yang Sudah Dikonfigurasi

1. **Migration Manager** - `internal/migration/migration.go`
2. **CLI Tools** - `cmd/migrate/main.go`
3. **Makefile Commands** - untuk kemudahan penggunaan
4. **Sample Migrations** - contoh migrations untuk tables users dan posts

## 🚀 Commands Utama

### Membuat Migration Baru

```bash
make migrate-create name="nama_migration"
```

### Menjalankan Migrations

```bash
make migrate-up
```

### Rollback Migration Terakhir

```bash
make migrate-down
```

### Cek Status Migrations

```bash
make migrate-status
```

## 📁 Struktur Files yang Dibuat

```
migrations/
├── 20250527000001_create_users_table.sql
├── 20250527231056_create_posts_table.sql
├── 20250527231209_add_email_verification_to_users.sql
└── 20250527231615_add_user_profile_fields.sql

cmd/
└── migrate/
    └── main.go

internal/
└── migration/
    └── migration.go

Makefile (updated)
MIGRATIONS.md (dokumentasi lengkap)
```

## 📊 Database Tables

### 1. Users Table

- ✅ Basic fields: `id`, `email`, `password`, `created_at`, `updated_at`
- ✅ Profile fields: `first_name`, `last_name`, `is_active`
- ✅ Email verification: `email_verified_at`, `email_verification_token`
- ✅ Extended profile: `phone`, `date_of_birth`, `gender`, `profile_image_url`
- ✅ Indexes pada `email`, `is_active`, `email_verification_token`, `phone`

### 2. Posts Table

- ✅ Basic fields: `id`, `title`, `content`, `user_id`, `status`
- ✅ Timestamps: `published_at`, `created_at`, `updated_at`
- ✅ Foreign key ke users table
- ✅ Indexes pada `user_id`, `status`, `published_at`

### 3. Migrations Table (auto-created)

- ✅ Tracking migrations yang sudah dijalankan
- ✅ Batch system untuk rollback

## 🎯 Next Steps

### 1. Membuat Migration Baru

```bash
# Contoh: Menambah table categories
make migrate-create name="create_categories_table"

# Edit file yang dibuat di migrations/
# Tambahkan SQL untuk UP dan DOWN
```

### 2. Development Workflow

```bash
# 1. Buat migration
make migrate-create name="add_category_to_posts"

# 2. Edit migration file
# 3. Test di development
make migrate-up

# 4. Jika ada masalah, rollback
make migrate-down

# 5. Fix migration, lalu jalankan lagi
make migrate-up
```

### 3. Production Deployment

```bash
# Backup database dulu
# Lalu jalankan migrations
make migrate-up

# Verify hasil
make migrate-status
```

## 🔧 Advanced Usage

### Direct Go Commands (tanpa Makefile)

```bash
# Membuat migration
go run cmd/migrate/main.go create "nama_migration"

# Menjalankan migrations
go run cmd/migrate/main.go migrate

# Rollback
go run cmd/migrate/main.go rollback

# Status
go run cmd/migrate/main.go status

# Help
go run cmd/migrate/main.go help
```

## ⚠️ Important Notes

1. **Manual Migrations** - Sistem ini dirancang untuk migrations manual, bukan otomatis
2. **Backup First** - Selalu backup database sebelum menjalankan migrations di production
3. **Test Rollback** - Pastikan rollback SQL sudah ditest
4. **Environment** - Pastikan `DATABASE_URL` sudah dikonfigurasi di `.env`

## 📖 Dokumentasi Lengkap

Lihat `MIGRATIONS.md` untuk dokumentasi lengkap termasuk:

- Best practices
- Troubleshooting
- Advanced features
- Production deployment guide

---

🎉 **Sistem migrations sudah siap digunakan!**

Untuk bantuan lebih lanjut, jalankan: `make migrate-help` atau `go run cmd/migrate/main.go help`
