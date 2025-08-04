# Unit Testing Documentation

Folder ini berisi unit testing untuk semua endpoint dalam Go Gin boilerplate API.

## Struktur Folder

```
tests/
├── README.md                    # File ini
└── unit/
    └── auth_controller_test.go  # Unit tests untuk auth controller
```

## Menjalankan Tests

### Jalankan semua unit tests
```bash
go test ./tests/unit/... -v
```


## Jenis Unit Tests

### Auth Controller Tests
- `TestRegister_Success` - Menguji register berhasil
- `TestRegister_InvalidRequest` - Menguji register dengan request tidak valid
- `TestLogin_Success` - Menguji login berhasil
- `TestLogout_Success` - Menguji logout berhasil
- `TestRefreshToken_Success` - Menguji refresh token berhasil
- `TestRefreshToken_NoToken` - Menguji refresh token tanpa token
- `TestForgotPassword_Success` - Menguji forgot password berhasil
- `TestVerifyOTP_Success` - Menguji verify OTP berhasil
- `TestResetPassword_Success` - Menguji reset password berhasil

## Fitur Unit Testing

### Mock Service
- Menggunakan mock sederhana untuk AuthService
- Tidak memerlukan database atau dependencies eksternal
- Testing cepat dan terisolasi

### Test Cases
- Success cases - menguji fungsi berjalan normal
- Error cases - menguji handling error
- Edge cases - menguji input tidak valid

### Assertions
- Menggunakan testing standar Go
- Tidak memerlukan library eksternal
- Error messages dalam bahasa Indonesia

## Cara Menambah Test Baru

1. Buat fungsi test dengan prefix `Test`
2. Gunakan mock service untuk dependencies
3. Test success dan error scenarios
4. Berikan nama yang deskriptif

Contoh:
```go
func TestNewFunction_Success(t *testing.T) {
    // Setup
    mockService := &MockAuthService{
        someFunc: func() error {
            return nil
        },
    }
    controller := controllers.NewAuthController(mockService)
    
    // Test logic here
    // Assertions here
}
```

## Best Practices

1. Setiap test harus independen
2. Gunakan nama test yang deskriptif
3. Test success dan failure scenarios
4. Mock semua dependencies eksternal
5. Gunakan error messages yang jelas 