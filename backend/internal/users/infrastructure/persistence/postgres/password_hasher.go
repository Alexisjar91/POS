package postgres

import (
	"golang.org/x/crypto/bcrypt"
)

// passwordHasher implementa domain.PasswordHasher usando bcrypt.
type passwordHasher struct{}

// NewPasswordHasher crea una nueva instancia de passwordHasher.
// Retorna el tipo concreto *passwordHasher; la asignación a domain.PasswordHasher
// se hace desde el llamante gracias a la inferencia de interfaces de Go.
func NewPasswordHasher() *passwordHasher {
	return &passwordHasher{}
}

// Hash genera un hash bcrypt de la contraseña en texto plano.
func (p *passwordHasher) Hash(plain string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// Compare verifica si una contraseña en texto plano coincide con un hash bcrypt.
// Retorna nil si coinciden, error en caso contrario.
func (p *passwordHasher) Compare(plain, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain))
}
