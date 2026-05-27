# Go-Auth-Api

## Required

* MySQL
* Go(Golang)
* Fiber
* GORM
* JWT & bcryptjs
* Vscode
    * Thunder Client (For API Testing)

## Develop

### First time

```bash
go mod init go-auth-api
go get -u github.com/gofiber/fiber/v2
go get -u gorm.io/gorm
go get -u gorm.io/driver/mysql
go get -u golang.org/x/crypto/bcrypt
go get -u github.com/golang-jwt/jwt/v5
go get -u github.com/joho/godotenv
```
*(อย่าลืมตั้งค่า `DATABASE_URL` ในไฟล์ `.env` ให้เรียบร้อย)*

### Update schema

1. Run this command `npx prisma migrate dev --name init`
2. ล้างฐานข้อมูลเก่า `npx prisma migrate reset `

### Create Folder
```bash
touch server.js config.js routes.js
mkdir controllers models middlewares
touch controllers/authController.js middlewares/middleware.js models/userModel.js
```

## Running

### Schema

```bash
npx prisma generate
npx prisma db push (ใช้ db push จะสะดวกกว่า migrate ตอนที่เราแค่เพิ่มฟิลด์เล็กๆ)
```

### Server

```bash
node server.js
```

### API Endpoints

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
```

**3. Protected** (`GET http://localhost:3000/api/protected`)
* **Header:** `Authorization: <วาง_accessToken_ตรงนี้>`

**4. Refresh Token** (`POST http://localhost:3000/api/refresh`)
```json
{
  "token": "วาง_refreshToken_ตัวจริง_ตรงนี้"
}
(ระบบจะคืนค่า accessToken ใบใหม่กลับมาให้)
```

**5. Logout** (`POST http://localhost:3000/api/logout`)
* **Header:** `Authorization: <วาง_accessToken_ตรงนี้>`
(ทดสอบความปลอดภัย: หลังจาก Logout สำเร็จ หากพยายามยิง API ในข้อ 4 อีกครั้ง ระบบจะแจ้ง 403 Forbidden ทันที!)