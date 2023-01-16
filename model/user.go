package model

type UserID int64

type User struct {
	ID       UserID `gorm:"primaryKey,autoIncrement" json:"-"`
	Email    string `gorm:"VARCHAR(255)" json:"email"`
	Password string `gorm:"VARCHAR(255)" json:"-"`
	Role     string `gorm:"VARCHAR(255)" json:"role"`
	Enabled  bool   `json:"-"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserRepository interface {
	GetByEmail(email string) (*User, error)
}
