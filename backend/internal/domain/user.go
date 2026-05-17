package domain

import "time"

type User struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

type RefreshToken struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      `gorm:"not null;index"`
	Token     string    `gorm:"not null ;uniqueIndex"`
	ExpiresAt time.Time `gorm:"not null"`
	CreatedAt time.Time
}

type RefreshTokenRepository interface {
	Create(token *RefreshToken) error
	FindByToken(token string) (*RefreshToken, error)
	DeleteByUserID(userID uint) error
	DeleteByToken(token string) error
}

type UserRepository interface {
	FindByID(id uint) (*User, error)
	FindByEmail(email string) (*User, error)
	FindByUser(username string) (*User, error)
	Create(user *User) error
}

type UserService interface {
	UserRegister(user *User) error
	UserLogin(email, password string) (accessToken string, refreshToken string, err error)
	RefreshAccessToken(refreshToken string) (string, string, error)
	UserLogout(userID uint) error
}
