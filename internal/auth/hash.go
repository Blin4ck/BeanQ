package auth

import "golang.org/x/crypto/bcrypt"

// HashPassword хеширует пароль с помощью bcrypt.
func HashPassword(pasword string) (string, error) {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(pasword), bcrypt.DefaultCost)
	return string(hashPassword), err

}
