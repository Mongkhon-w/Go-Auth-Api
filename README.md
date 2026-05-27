# 🚀 Golang Fiber & GORM Authentication API Standard Template Guide

## 🛠️ Required Tools

* Database: MySQL (XAMPP, Laragon, etc.)
* Language: Go (Golang)
* Framework: Fiber
* ORM: GORM
* Security: JWT, crypto/bcrypt (Password), และ crypto/sha256 (Refresh Token)
* Editor: VS Code
    * Extension: Thunder Client (For API Testing)

## 🏗️ Development Setup
### Initialize Project (First time)
```bash
go mod init go-auth-api
go get -u github.com/gofiber/fiber/v2
go get -u gorm.io/gorm
go get -u gorm.io/driver/mysql
go get -u golang.org/x/crypto/bcrypt
go get -u github.com/golang-jwt/jwt/v5
go get -u github.com/joho/godotenv
```
สร้างไฟล์ .env
*(อย่าลืมตั้งค่า `DATABASE_URL` ในไฟล์ `.env` ให้เรียบร้อย)*

### Environment Variables (.env)
```bash
สร้างไฟล์ .env
(อย่าลืมตั้งค่า `DATABASE_URL` ในไฟล์ `.env` ให้เรียบร้อย)
```

### Create Folder Structure
```bash
mkdir controllers database middlewares models routes
touch main.go controllers/authController.go database/db.go middlewares/authMiddleware.go models/user.go routes/routes.go
```
### Database Migration
```bash
มีถัง Database เปล่าๆ เตรียมไว้ใน MySQL เช่น my_db
```

## 🏃‍♂️ Running the Server
```bash
# 1. คอมไพล์โค้ดเป็นไฟล์ .exe
go build -o api.exe

# 2. สั่งรันเซิร์ฟเวอร์
.\api.exe

(ถ้าไม่มีปัญหาเรื่องโดนบล็อก สามารถใช้ go run main.go ได้ตามปกติ)
```

## 📡 API Endpoints Testing (Thunder Client)

**1. Register** (`POST http://localhost:3000/api/register`)
```json
{
  "username": "testuser",
  "password": "password123"
}
```

**2. Login** (`POST http://localhost:3000/api/login`)
```json
{
  "username": "testuser",
  "password": "password123"
}
(ระบบจะคืนค่า accessToken และทำการฝัง refreshToken ลงในแท็บ Cookies แบบ HttpOnly อัตโนมัติ)
```

**3. Protected** (`GET http://localhost:3000/api/protected`)
* **Header:** `Authorization -> Bearer <วาง_accessToken_ตรงนี้>`

**4. Refresh Token** (`POST http://localhost:3000/api/refresh`)
```json
        { } Body: (ปล่อยว่างเปล่า ไม่ต้องใส่อะไรเลย!)
(ระบบจะดึง HttpOnly Cookie ที่ซ่อนอยู่อัตโนมัติ ไปตรวจสอบด้วย SHA-256 และคืนค่า accessToken ใบใหม่กลับมา)
```

**5. Logout** (`POST http://localhost:3000/api/logout`)
* **Header:** `Authorization -> Bearer <วาง_accessToken_ตรงนี้>`
(ทดสอบความปลอดภัย: หลังจาก Logout สำเร็จ คุกกี้จะถูกลบ หากพยายามยิง API ในข้อ 4 อีกครั้ง ระบบจะแจ้ง 403 Forbidden ทันที)