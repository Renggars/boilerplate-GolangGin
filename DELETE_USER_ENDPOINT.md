# Delete User Endpoint

## Overview
Endpoint untuk menghapus user dengan persyaratan keamanan yang ketat. Menggunakan soft delete untuk memastikan keamanan data. Jika user menghapus akunnya sendiri, cookie akan dihapus untuk logout otomatis. Jika admin menghapus user lain, user yang dihapus akan otomatis logout dari semua sesi mereka.

## Endpoint
```
DELETE /api/user/{id}
```

## Authorization
- **User biasa**: Hanya bisa menghapus akun mereka sendiri
- **Admin**: Bisa menghapus user mana saja dengan ID yang diberikan
- **Auto Logout**: Jika user/admin menghapus akunnya sendiri, cookie akan dihapus untuk logout otomatis

## Request Headers
```
Authorization: Bearer <token>
```

## Path Parameters
- `id` (integer, required): ID user yang akan dihapus

## Response

### Success (200 OK)
```json
{
  "status_code": 200,
  "message": "success delete user"
}
```

### Error Responses

#### 400 Bad Request
```json
{
  "message": "Invalid user ID"
}
```

#### 401 Unauthorized
```json
{
  "message": "Unauthorized"
}
```

#### 403 Forbidden
```json
{
  "message": "Access denied: you can only delete your own account"
}
```

#### 404 Not Found
```json
{
  "message": "User not found"
}
```

#### 500 Internal Server Error
```json
{
  "message": "Failed to delete user"
}
```

## Examples

### User menghapus akun sendiri
```bash
curl -X DELETE \
  http://localhost:8080/api/user/1 \
  -H "Authorization: Bearer <token>"
```

### Admin menghapus user lain
```bash
curl -X DELETE \
  http://localhost:8080/api/user/2 \
  -H "Authorization: Bearer <admin_token>"
```

### Admin menghapus akunnya sendiri
```bash
curl -X DELETE \
  http://localhost:8080/api/user/1 \
  -H "Authorization: Bearer <admin_token>"
```

## Security Features
1. **Authentication**: Semua request harus menyertakan token yang valid
2. **Authorization**: 
   - User biasa hanya bisa menghapus akun mereka sendiri
   - Admin bisa menghapus user mana saja
3. **Validation**: ID user divalidasi sebelum diproses
4. **Error Handling**: Penanganan error yang komprehensif
5. **Auto Logout**: 
   - Jika user/admin menghapus akunnya sendiri, cookie akan dihapus untuk logout otomatis
   - Jika admin menghapus user lain, user yang dihapus akan otomatis logout dari semua sesi mereka
6. **Soft Delete**: Menggunakan soft delete untuk mempertahankan data dan memastikan user yang dihapus tidak bisa mengakses sistem

## Implementation Details
- **Repository Layer**: `DeleteUser(id int) error` dengan soft delete menggunakan `deleted_at`
- **Service Layer**: `DeleteUser(id int) error` dengan validasi user exists dan tidak sudah di-delete
- **Controller Layer**: `DeleteUser(ctx *gin.Context)` dengan authorization logic dan cookie management
- **Middleware**: `AuthUser` untuk autentikasi dengan pengecekan `deleted_at`
- **Model**: Field `DeletedAt` untuk soft delete

## Testing
Endpoint ini telah dilengkapi dengan unit test yang mencakup:
- User berhasil menghapus akun sendiri
- Admin berhasil menghapus user lain
- User dilarang menghapus user lain
- Invalid user ID
- User tidak ditemukan
- Service error
- Invalid user context 