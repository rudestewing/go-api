# Access Token Based Authentication

Sistem autentikasi telah diubah dari JWT-based menjadi Access Token based yang disimpan di database. Ini memberikan kontrol penuh terhadap autentikasi dan memudahkan revocation token.

## Fitur Utama

1. **Token tersimpan di database** - Memberikan kontrol penuh
2. **Soft delete untuk revocation** - Menggunakan `deleted_at` dari GORM
3. **Expiration handling** - Token memiliki masa kadaluarsa
4. **Secure token generation** - Menggunakan crypto/rand untuk keamanan
5. **User relationship** - Token terhubung dengan user

## API Endpoints

### 1. Login

```http
POST /api/v1/auth/login
Content-Type: application/json

{
    "email": "user@example.com",
    "password": "password"
}
```

**Response:**

```json
{
  "success": true,
  "data": {
    "access_token": "64-char-hex-token",
    "expires_at": "2025-06-02T20:30:00Z",
    "user": {
      "id": 1,
      "name": "User Name",
      "email": "user@example.com",
      "role_id": 2
    }
  }
}
```

### 2. Register

```http
POST /api/v1/auth/register
Content-Type: application/json

{
    "name": "New User",
    "email": "newuser@example.com",
    "password": "securepassword"
}
```

### 3. Logout (Revoke current token)

```http
POST /api/v1/auth/logout
Authorization: Bearer your-access-token
```

### 4. Logout All (Revoke all user tokens)

```http
POST /api/v1/auth/logout-all
Authorization: Bearer your-access-token
```

### 5. Protected Endpoints

```http
GET /api/v1/user/profile
Authorization: Bearer your-access-token
```

## Model Structure

### AccessToken

```go
type AccessToken struct {
    gorm.Model
    Token     string    `gorm:"uniqueIndex;not null" json:"token"`
    UserID    uint      `gorm:"not null" json:"user_id"`
    ExpiresAt time.Time `gorm:"not null" json:"expires_at"`

    User User `gorm:"foreignKey:UserID" json:"user"`
}
```

**Fitur:**

- `Token`: 64-character hex string (256-bit random)
- `UserID`: Foreign key ke tabel users
- `ExpiresAt`: Timestamp kadaluarsa token
- `DeletedAt`: Untuk soft delete (revocation)

## Repository Methods

### AccessTokenRepository Interface

```go
type AccessTokenRepository interface {
    Create(userID uint, expiresIn time.Duration) (*model.AccessToken, error)
    FindByToken(token string) (*model.AccessToken, error)
    RevokeToken(token string) error
    RevokeAllUserTokens(userID uint) error
    DeleteExpiredTokens() error
    CleanupExpiredTokens() error
}
```

### Key Methods

1. **Create**: Membuat token baru dengan expiration
2. **FindByToken**: Mencari token yang valid (tidak deleted dan belum expired)
3. **RevokeToken**: Soft delete token (set deleted_at)
4. **RevokeAllUserTokens**: Revoke semua token milik user
5. **CleanupExpiredTokens**: Hard delete token yang expired

## Service Layer

### AuthService Methods

```go
func (s *AuthService) Login(ctx context.Context, email, password string) (*model.AccessToken, error)
func (s *AuthService) ValidateToken(token string) (*model.AccessToken, error)
func (s *AuthService) Logout(token string) error
func (s *AuthService) LogoutAll(userID uint) error
```

## Middleware

### AuthMiddleware

Middleware baru menggunakan access token validation:

```go
func AuthMiddleware(authService *service.AuthService) fiber.Handler
```

**Context yang di-set:**

- `user_id`: ID user
- `user`: Object user lengkap
- `access_token`: Object access token

## Database Migration

Migration untuk tabel `access_tokens`:

```sql
CREATE TABLE access_tokens (
    id SERIAL PRIMARY KEY,
    token VARCHAR(255) UNIQUE NOT NULL,
    user_id INTEGER NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
```

## Security Benefits

1. **Immediate Revocation**: Token bisa di-revoke instantly
2. **Session Management**: Kontrol penuh atas session user
3. **Audit Trail**: Semua token tersimpan dengan timestamp
4. **Force Logout**: Bisa logout user dari semua device
5. **Secure Generation**: Token menggunakan crypto/rand

## Maintenance

### Cleanup Expired Tokens

Jalankan secara berkala untuk membersihkan token expired:

```go
accessTokenRepo.CleanupExpiredTokens()
```

Ini bisa dijalankan sebagai cron job atau background task.

## Migration dari JWT

Perubahan utama dari sistem sebelumnya:

1. ✅ Tidak lagi menggunakan JWT
2. ✅ Token disimpan di database
3. ✅ Middleware diupdate untuk validation dari database
4. ✅ Response login berubah format
5. ✅ Tambahan endpoint logout dan logout-all

## Testing

Setelah menjalankan migration dan aplikasi:

1. Test register user baru
2. Test login dan dapatkan access_token
3. Test access protected endpoint dengan token
4. Test logout (token tidak bisa digunakan lagi)
5. Test logout-all
