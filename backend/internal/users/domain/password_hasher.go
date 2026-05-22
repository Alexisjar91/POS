package domain

// PasswordHasher define el contrato para el hashing y verificación de contraseñas.
type PasswordHasher interface {
	// Hash genera un hash seguro de la contraseña.
	Hash(plainPassword string) (string, error)

	// Compare verifica si una contraseña en texto plano coincide con un hash.
	// Retorna nil si coinciden, error si no.
	Compare(plainPassword string, hash string) error
}
