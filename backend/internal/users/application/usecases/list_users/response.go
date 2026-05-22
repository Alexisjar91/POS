package listusers

type UsuarioDTO struct {
	ID        string
	Email     string
	FullName  string
	Active    bool
	CreatedAt string
}

type RespuestaListarUsuarios struct {
	Usuarios   []UsuarioDTO
	Pagina     int
	TamanoPagina int
	Total      int64
	TotalPages int
}
