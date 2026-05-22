package listusers

import (
	"context"
	"errors"
	"testing"

	"github.com/Alexisjar91/POS/internal/users/domain"
	"github.com/Alexisjar91/POS/pkg/especificacion"
	"github.com/Alexisjar91/POS/pkg/paginacion"
)

// --- Mocks ---

type mockUserRepo struct {
	listarFunc func(ctx context.Context, spec especificacion.Especificacion, pag paginacion.Paginacion) ([]*domain.User, error)
}

func (m *mockUserRepo) Crear(_ context.Context, _ *domain.User) (*domain.User, error) {
	panic("no debería llamarse")
}
func (m *mockUserRepo) ObtenerPorID(_ context.Context, _ string) (*domain.User, error) {
	panic("no debería llamarse")
}
func (m *mockUserRepo) ObtenerPorEmail(_ context.Context, _ string) (*domain.User, error) {
	panic("no debería llamarse")
}
func (m *mockUserRepo) Actualizar(_ context.Context, _ *domain.User) (*domain.User, error) {
	panic("no debería llamarse")
}
func (m *mockUserRepo) Listar(ctx context.Context, spec especificacion.Especificacion, pag paginacion.Paginacion) ([]*domain.User, error) {
	if m.listarFunc != nil {
		return m.listarFunc(ctx, spec, pag)
	}
	return nil, nil
}
func (m *mockUserRepo) ExistePorEmail(_ context.Context, _ string) (bool, error) {
	panic("no debería llamarse")
}

type mockAuthSvc struct {
	verificarFunc func(ctx context.Context, userID, code string) (bool, error)
	esOwnerFunc   func(ctx context.Context, userID string) (bool, error)
}

func (m *mockAuthSvc) VerificarPermiso(ctx context.Context, userID, code string) (bool, error) {
	if m.verificarFunc != nil {
		return m.verificarFunc(ctx, userID, code)
	}
	return false, nil
}
func (m *mockAuthSvc) EsUsuarioOWNER(ctx context.Context, userID string) (bool, error) {
	if m.esOwnerFunc != nil {
		return m.esOwnerFunc(ctx, userID)
	}
	return false, nil
}

// --- Helpers ---

func comandoValido() *ComandoListarUsuarios {
	return &ComandoListarUsuarios{
		EjecutorID:     "user-executor",
		Especificacion: especificacion.Especificacion{},
		Paginacion:     paginacion.Paginacion{Pagina: 1, TamanoPagina: 10},
	}
}

func strPtr(s string) *string { return &s }
func boolPtr(b bool) *bool    { return &b }

func defaultMockAuth(_ context.Context, _ string, _ string) (bool, error) {
	return true, nil
}

func crearUsuario(id, email, fullName string, active bool) *domain.User {
	return domain.NuevoUsuarioDesdeBD(
		id, email, fullName, "hash", active, strPtr("creator"), "2024-01-01",
	)
}

// --- Tests existentes (adaptados) ---

func TestListarUsuarios_Exito_ConCeroUsuarios(t *testing.T) {
	repo := &mockUserRepo{
		listarFunc: func(_ context.Context, _ especificacion.Especificacion, _ paginacion.Paginacion) ([]*domain.User, error) {
			return []*domain.User{}, nil
		},
	}
	auth := &mockAuthSvc{verificarFunc: defaultMockAuth}
	uc := NewListarUsuariosCasoDeUso(repo, auth)

	resp, err := uc.Ejecutar(context.Background(), comandoValido())

	if err != nil {
		t.Fatalf("Ejecutar returned unexpected error: %v", err)
	}
	if resp == nil {
		t.Fatal("Ejecutar returned nil response")
	}
	if len(resp.Usuarios) != 0 {
		t.Errorf("len(resp.Usuarios) = %d, want 0", len(resp.Usuarios))
	}
	if resp.Total != 0 {
		t.Errorf("Total = %d, want 0", resp.Total)
	}
}

func TestListarUsuarios_Exito_ConUnUsuario(t *testing.T) {
	users := []*domain.User{crearUsuario("user-1", "user1@test.com", "User One", true)}
	repo := &mockUserRepo{
		listarFunc: func(_ context.Context, _ especificacion.Especificacion, _ paginacion.Paginacion) ([]*domain.User, error) {
			return users, nil
		},
	}
	auth := &mockAuthSvc{verificarFunc: defaultMockAuth}
	uc := NewListarUsuariosCasoDeUso(repo, auth)

	resp, err := uc.Ejecutar(context.Background(), comandoValido())

	if err != nil {
		t.Fatalf("Ejecutar returned unexpected error: %v", err)
	}
	if len(resp.Usuarios) != 1 {
		t.Fatalf("len(resp.Usuarios) = %d, want 1", len(resp.Usuarios))
	}
	if resp.Usuarios[0].ID != "user-1" {
		t.Errorf("Usuarios[0].ID = %q, want %q", resp.Usuarios[0].ID, "user-1")
	}
	if resp.Total != 1 {
		t.Errorf("Total = %d, want 1", resp.Total)
	}
}

func TestListarUsuarios_SinPermiso(t *testing.T) {
	repo := &mockUserRepo{
		listarFunc: func(_ context.Context, _ especificacion.Especificacion, _ paginacion.Paginacion) ([]*domain.User, error) {
			t.Error("Listar should not be called when not authorized")
			return nil, nil
		},
	}
	auth := &mockAuthSvc{
		verificarFunc: func(_ context.Context, _, _ string) (bool, error) {
			return false, nil
		},
	}
	uc := NewListarUsuariosCasoDeUso(repo, auth)

	resp, err := uc.Ejecutar(context.Background(), comandoValido())

	if resp != nil {
		t.Error("expected nil response on access denied")
	}
	if !errors.Is(err, domain.ErrAccesoDenegado) {
		t.Errorf("err = %v, want %v", err, domain.ErrAccesoDenegado)
	}
}

func TestListarUsuarios_ValidacionFallida(t *testing.T) {
	repo := &mockUserRepo{
		listarFunc: func(_ context.Context, _ especificacion.Especificacion, _ paginacion.Paginacion) ([]*domain.User, error) {
			t.Error("Listar should not be called on validation failure")
			return nil, nil
		},
	}
	auth := &mockAuthSvc{verificarFunc: defaultMockAuth}
	uc := NewListarUsuariosCasoDeUso(repo, auth)

	cmd := &ComandoListarUsuarios{EjecutorID: ""}
	resp, err := uc.Ejecutar(context.Background(), cmd)

	if resp != nil {
		t.Error("expected nil response on validation error")
	}
	if !errors.Is(err, ErrEjecutorRequerido) {
		t.Errorf("err = %v, want %v", err, ErrEjecutorRequerido)
	}
}

func TestListarUsuarios_ErrorRepositorio(t *testing.T) {
	repo := &mockUserRepo{
		listarFunc: func(_ context.Context, _ especificacion.Especificacion, _ paginacion.Paginacion) ([]*domain.User, error) {
			return nil, domain.ErrRepositorio
		},
	}
	auth := &mockAuthSvc{verificarFunc: defaultMockAuth}
	uc := NewListarUsuariosCasoDeUso(repo, auth)

	resp, err := uc.Ejecutar(context.Background(), comandoValido())

	if resp != nil {
		t.Error("expected nil response on repository error")
	}
	if !errors.Is(err, domain.ErrRepositorio) {
		t.Errorf("err = %v, want %v", err, domain.ErrRepositorio)
	}
}

func TestListarUsuarios_NoLlamaRepoSiNoAutorizado(t *testing.T) {
	listarLlamado := false
	repo := &mockUserRepo{
		listarFunc: func(_ context.Context, _ especificacion.Especificacion, _ paginacion.Paginacion) ([]*domain.User, error) {
			listarLlamado = true
			return nil, nil
		},
	}
	auth := &mockAuthSvc{
		verificarFunc: func(_ context.Context, _, _ string) (bool, error) {
			return false, nil
		},
	}
	uc := NewListarUsuariosCasoDeUso(repo, auth)

	resp, _ := uc.Ejecutar(context.Background(), comandoValido())

	if resp != nil {
		t.Error("expected nil response on access denied")
	}
	if listarLlamado {
		t.Error("Listar should not be called when authorization fails")
	}
}

// --- Nuevos tests: especificación y paginación ---

func TestListarUsuarios_PasaEspecVaciaAlRepositorio(t *testing.T) {
	var specRecibida especificacion.Especificacion
	repo := &mockUserRepo{
		listarFunc: func(_ context.Context, spec especificacion.Especificacion, _ paginacion.Paginacion) ([]*domain.User, error) {
			specRecibida = spec
			return []*domain.User{}, nil
		},
	}
	auth := &mockAuthSvc{verificarFunc: defaultMockAuth}
	uc := NewListarUsuariosCasoDeUso(repo, auth)

	uc.Ejecutar(context.Background(), comandoValido())

	if len(specRecibida.Filtros) != 0 {
		t.Error("expected empty filtros for empty especificacion")
	}
}

func TestListarUsuarios_PasaFiltroAlRepositorio(t *testing.T) {
	var specRecibida especificacion.Especificacion
	repo := &mockUserRepo{
		listarFunc: func(_ context.Context, spec especificacion.Especificacion, _ paginacion.Paginacion) ([]*domain.User, error) {
			specRecibida = spec
			return []*domain.User{}, nil
		},
	}
	auth := &mockAuthSvc{verificarFunc: defaultMockAuth}
	uc := NewListarUsuariosCasoDeUso(repo, auth)

	cmd := &ComandoListarUsuarios{
		EjecutorID: "exec",
		Especificacion: especificacion.Especificacion{
			Filtros: []especificacion.CriterioFiltro{
				{Campo: "email", Operador: "=", Valor: "test@example.com"},
			},
		},
		Paginacion: paginacion.Paginacion{Pagina: 1, TamanoPagina: 10},
	}
	uc.Ejecutar(context.Background(), cmd)

	if len(specRecibida.Filtros) != 1 {
		t.Fatalf("expected 1 filtro, got %d", len(specRecibida.Filtros))
	}
	if specRecibida.Filtros[0].Campo != "email" {
		t.Errorf("expected field 'email', got %q", specRecibida.Filtros[0].Campo)
	}
}

func TestListarUsuarios_PasaPaginacion(t *testing.T) {
	var pagRecibida paginacion.Paginacion
	repo := &mockUserRepo{
		listarFunc: func(_ context.Context, _ especificacion.Especificacion, pag paginacion.Paginacion) ([]*domain.User, error) {
			pagRecibida = pag
			return []*domain.User{}, nil
		},
	}
	auth := &mockAuthSvc{verificarFunc: defaultMockAuth}
	uc := NewListarUsuariosCasoDeUso(repo, auth)

	cmd := &ComandoListarUsuarios{
		EjecutorID:     "exec",
		Especificacion: especificacion.Especificacion{},
		Paginacion:     paginacion.Paginacion{Pagina: 2, TamanoPagina: 25},
	}
	uc.Ejecutar(context.Background(), cmd)

	if pagRecibida.Pagina != 2 {
		t.Errorf("Pagina = %d, want 2", pagRecibida.Pagina)
	}
	if pagRecibida.TamanoPagina != 25 {
		t.Errorf("TamanoPagina = %d, want 25", pagRecibida.TamanoPagina)
	}
}

func TestListarUsuarios_DevuelveInfoPaginacion(t *testing.T) {
	users := []*domain.User{
		crearUsuario("u1", "u1@t.com", "Uno", true),
		crearUsuario("u2", "u2@t.com", "Dos", true),
	}
	repo := &mockUserRepo{
		listarFunc: func(_ context.Context, _ especificacion.Especificacion, _ paginacion.Paginacion) ([]*domain.User, error) {
			return users, nil
		},
	}
	auth := &mockAuthSvc{verificarFunc: defaultMockAuth}
	uc := NewListarUsuariosCasoDeUso(repo, auth)

	cmd := &ComandoListarUsuarios{
		EjecutorID:     "exec",
		Especificacion: especificacion.Especificacion{},
		Paginacion:     paginacion.Paginacion{Pagina: 1, TamanoPagina: 10},
	}
	resp, err := uc.Ejecutar(context.Background(), cmd)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Total != 2 {
		t.Errorf("Total = %d, want 2", resp.Total)
	}
	if resp.Pagina != 1 {
		t.Errorf("Pagina = %d, want 1", resp.Pagina)
	}
	if resp.TamanoPagina != 10 {
		t.Errorf("TamanoPagina = %d, want 10", resp.TamanoPagina)
	}
	if resp.TotalPages != 1 { // 2/10 = 0.2 → ceil = 1
		t.Errorf("TotalPages = %d, want 1", resp.TotalPages)
	}
}

func TestListarUsuarios_CalculaTotalPaginasCorrectamente(t *testing.T) {
	// Simulamos 25 usuarios con página size 10
	tests := []struct {
		total        int64
		tamanoPagina int
		expectedPages int
	}{
		{0, 10, 1},
		{1, 10, 1},
		{10, 10, 1},
		{11, 10, 2},
		{20, 10, 2},
		{21, 10, 3},
		{25, 10, 3},
	}

	for _, tc := range tests {
		t.Run("", func(t *testing.T) {
			// Simulamos retornar usuarios suficientes para el total
			users := make([]*domain.User, 0)
			for i := 0; i < int(tc.total); i++ {
				users = append(users, crearUsuario("u"+string(rune(i)), "u"+string(rune(i))+"@t.com", "User", true))
			}

			repo := &mockUserRepo{
				listarFunc: func(_ context.Context, _ especificacion.Especificacion, _ paginacion.Paginacion) ([]*domain.User, error) {
					return users, nil
				},
			}
			auth := &mockAuthSvc{verificarFunc: defaultMockAuth}
			uc := NewListarUsuariosCasoDeUso(repo, auth)

			cmd := &ComandoListarUsuarios{
				EjecutorID:     "exec",
				Especificacion: especificacion.Especificacion{},
				Paginacion:     paginacion.Paginacion{Pagina: 1, TamanoPagina: tc.tamanoPagina},
			}
			resp, _ := uc.Ejecutar(context.Background(), cmd)

			if resp.TotalPages != tc.expectedPages {
				t.Errorf("Total=%d, TamanoPagina=%d: TotalPages = %d, want %d",
					tc.total, tc.tamanoPagina, resp.TotalPages, tc.expectedPages)
			}
		})
	}
}
