package createuser

// RespuestaCrearUsuario contiene los datos del usuario creado.
type RespuestaCrearUsuario struct {
	ID        string
	Email     string
	FullName  string
	Active    bool
	CreatedAt string
}
