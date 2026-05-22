package domain

// User representa un usuario del sistema.
// Los campos no exportados garantizan inmutabilidad desde fuera del paquete.
// No contiene colecciones de roles; la relación many-to-many se maneja
// a través de la entidad UserRole.
type User struct {
	id           string
	email        string
	passwordHash string
	fullName     string
	active       bool
	createdBy    *string // nil solo para OWNER
	createdAt    string  // timestamp como string para evitar dependencias externas
}

// NuevoUsuario crea un nuevo User normal (no OWNER).
// Valida que ningún campo requerido esté vacío.
// active se inicializa a true, createdBy se guarda como puntero.
func NuevoUsuario(email, fullName, passwordHash, createdBy string) (*User, error) {
	if email == "" {
		return nil, ErrEmailRequerido
	}
	if fullName == "" {
		return nil, ErrNombreRequerido
	}
	if passwordHash == "" {
		return nil, ErrPasswordHashRequerido
	}
	if createdBy == "" {
		return nil, ErrCreatedByRequerido
	}
	return &User{
		email:        email,
		fullName:     fullName,
		passwordHash: passwordHash,
		active:       true,
		createdBy:    &createdBy,
	}, nil
}

// NuevoUsuarioOwner crea un nuevo User con rol OWNER.
// Es el único caso donde createdBy puede ser nil.
func NuevoUsuarioOwner(email, fullName, passwordHash string) (*User, error) {
	if email == "" {
		return nil, ErrEmailRequerido
	}
	if fullName == "" {
		return nil, ErrNombreRequerido
	}
	if passwordHash == "" {
		return nil, ErrPasswordHashRequerido
	}
	return &User{
		email:        email,
		fullName:     fullName,
		passwordHash: passwordHash,
		active:       true,
		createdBy:    nil,
	}, nil
}

// NuevoUsuarioDesdeBD reconstruye un User desde datos persistentes.
// No realiza validación: asume datos consistentes provenientes de la base de datos.
func NuevoUsuarioDesdeBD(id, email, fullName, passwordHash string, active bool, createdBy *string, createdAt string) *User {
	return &User{
		id:           id,
		email:        email,
		fullName:     fullName,
		passwordHash: passwordHash,
		active:       active,
		createdBy:    createdBy,
		createdAt:    createdAt,
	}
}

// ID retorna el identificador único del usuario.
func (u *User) ID() string {
	return u.id
}

// Email retorna el email del usuario.
func (u *User) Email() string {
	return u.email
}

// FullName retorna el nombre completo del usuario.
func (u *User) FullName() string {
	return u.fullName
}

// PasswordHash retorna el hash de la contraseña del usuario.
func (u *User) PasswordHash() string {
	return u.passwordHash
}

// IsActive retorna true si el usuario está activo.
func (u *User) IsActive() bool {
	return u.active
}

// Disable desactiva el usuario.
// Retorna ErrUsuarioYaInactivo si el usuario ya está inactivo.
func (u *User) Disable() error {
	if !u.active {
		return ErrUsuarioYaInactivo
	}
	u.active = false
	return nil
}

// Enable activa el usuario.
// Retorna ErrUsuarioYaActivo si el usuario ya está activo.
func (u *User) Enable() error {
	if u.active {
		return ErrUsuarioYaActivo
	}
	u.active = true
	return nil
}

// SetPasswordHash asigna un nuevo hash de contraseña.
// No realiza validación: el hash es opaco para el dominio.
func (u *User) SetPasswordHash(hash string) {
	u.passwordHash = hash
}

// CreatedBy retorna el identificador del usuario que creó este usuario.
// Puede retornar nil para el usuario OWNER.
func (u *User) CreatedBy() *string {
	return u.createdBy
}

// CreatedAt retorna el timestamp de creación del usuario.
func (u *User) CreatedAt() string {
	return u.createdAt
}
