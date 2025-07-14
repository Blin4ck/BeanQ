package entity

import (
	"coffe/internal/common"
	"errors"
	"strings"

	"github.com/google/uuid"
)

// User представляет пользователя системы.
type User struct {
	ID       uuid.UUID `json:"id" db:"id"`
	Name     string    `json:"name" db:"name"`
	Surname  string    `json:"surname" db:"surname"`
	Email    string    `json:"email" db:"email"`
	Password string    `json:"-" db:"password"`
	Role     *Role     `json:"role,omitempty" db:"role"`
}

// Role представляет роль пользователя.
type Role struct {
	ID   int    `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}

// ValidateUser проверяет корректность данных пользователя.
func ValidateUser(user *common.User) error {
	if strings.TrimSpace(user.Email) == "" {
		return errors.New("email не может быть пустым")
	}
	if len(user.Password) < 6 {
		return errors.New("пароль слишком короткий")
	}
	return nil
}
