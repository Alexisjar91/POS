package getuser

type RespuestaObtenerUsuario struct {
	ID        string
	Email     string
	FullName  string
	Active    bool
	CreatedAt string
	// NOTA: No se incluye PasswordHash por seguridad
}
