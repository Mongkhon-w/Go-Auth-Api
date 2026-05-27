package models

type User struct {
	ID           uint   `gorm:"primaryKey" json:"id"`
	Username     string `gorm:"unique;not null" json:"username"`
	Password     string `gorm:"not null" json:"-"` // json:"-" คือไม่ให้เผลอส่งรหัสผ่านกลับไปหน้าบ้าน
	RefreshToken string `gorm:"type:text" json:"-"`
}