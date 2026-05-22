package getuser

import (
	"github.com/Alexisjar91/POS/pkg/especificacion"
	"github.com/Alexisjar91/POS/pkg/paginacion"
	"context"
	"errors"
	"testing"

	"github.com/Alexisjar91/POS/internal/users/domain"
)

// --- Mocks ---

// mockUserRepo implementa domain.UserRepository para tests.
type mockUserRepo struct {
	obtenerPorIDFunc func(ctx context.Context, id string) (*domain.User, error)
}

func (m *mockUserRepo) Crear(_ context.Context, _ *domain.User) (*domain.User, error) {
	panic("no debería llamarse")
}

func (m *mockUserRepo) ObtenerPorID(ctx context.Context, id string) (*domain.User, error) {
	if m.obtenerPorIDFunc != nil {
		return m.obtenerPorIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *mockUserRepo) ObtenerPorEmail(_ context.Context, _ string) (*domain.User, error) {
	panic("no debería llamarse")
}

func (m *mockUserRepo) Actualizar(_ context.Context, _ *domain.User) (*domain.User, error) {
	panic("no debería llamarse")
}

func (m *mockUserRepo) Listar(_ context.Context, _ especificacion.Especificacion, _ paginacion.Paginacion) ([]*domain.User, error) {
	panic("no debería llamarse")
}

func (m *mockUserRepo) ExistePorEmail(_ context.Context, _ string) (bool, error) {
	panic("no debería llamarse")
}

// mockAuthSvc implementa domain.AuthorizationService para tests.
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

func crearCmdValido() *ComandoObtenerUsuario {
	return &ComandoObtenerUsuario{
		UserID:     "user-target",
		EjecutorID: "user-executor",
	}
}

func usuarioExistente() *domain.User {
	return domain.NuevoUsuarioDesdeBD(
		"user-target",
		"target@test.com",
		"Target User",
		"hash",
		true,
		strPtr("creator"),
		"2024-01-01",
	)
}

func strPtr(s string) *string {
	return &s
}

func defaultMockAuth(_ context.Context, _ string, _ string) (bool, error) {
	return true, nil
}

// --- Tests ---

func TestObtenerUsuario_Exito(t *testing.T) {
	repo := &mockUserRepo{
		obtenerPorIDFunc: func(_ context.Context, _ string) (*domain.User, error) {
			return usuarioExistente(), nil
		},
	}
	auth := &mockAuthSvc{
		verificarFunc: defaultMockAuth,
	}
	uc := NewObtenerUsuarioCasoDeUso(repo, auth)

	resp, err := uc.Ejecutar(context.Background(), crearCmdValido())

	if err != nil {
		t.Fatalf("Ejecutar returned unexpected error: %v", err)
	}
	if resp == nil {
		t.Fatal("Ejecutar returned nil response")
	}
	if resp.ID != "user-target" {
		t.Errorf("resp.ID = %q, want %q", resp.ID, "user-target")
	}
	if resp.Email != "target@test.com" {
		t.Errorf("resp.Email = %q, want %q", resp.Email, "target@test.com")
	}
	if resp.FullName != "Target User" {
		t.Errorf("resp.FullName = %q, want %q", resp.FullName, "Target User")
	}
	if !resp.Active {
		t.Error("resp.Active should be true")
	}
	if resp.CreatedAt != "2024-01-01" {
		t.Errorf("resp.CreatedAt = %q, want %q", resp.CreatedAt, "2024-01-01")
	}
}

func TestObtenerUsuario_SinPermiso(t *testing.T) {
	repo := &mockUserRepo{
		obtenerPorIDFunc: func(_ context.Context, _ string) (*domain.User, error) {
			t.Error("ObtenerPorID should not be called when not authorized")
			return nil, nil
		},
	}
	auth := &mockAuthSvc{
		verificarFunc: func(_ context.Context, _, _ string) (bool, error) {
			return false, nil
		},
	}
	uc := NewObtenerUsuarioCasoDeUso(repo, auth)

	resp, err := uc.Ejecutar(context.Background(), crearCmdValido())

	if resp != nil {
		t.Error("expected nil response on access denied")
	}
	if !errors.Is(err, domain.ErrAccesoDenegado) {
		t.Errorf("err = %v, want %v", err, domain.ErrAccesoDenegado)
	}
}

func TestObtenerUsuario_UsuarioNoEncontrado(t *testing.T) {
	repo := &mockUserRepo{
		obtenerPorIDFunc: func(_ context.Context, _ string) (*domain.User, error) {
			return nil, domain.ErrUsuarioNoEncontrado
		},
	}
	auth := &mockAuthSvc{
		verificarFunc: defaultMockAuth,
	}
	uc := NewObtenerUsuarioCasoDeUso(repo, auth)

	resp, err := uc.Ejecutar(context.Background(), crearCmdValido())

	if resp != nil {
		t.Error("expected nil response on user not found")
	}
	if !errors.Is(err, domain.ErrUsuarioNoEncontrado) {
		t.Errorf("err = %v, want %v", err, domain.ErrUsuarioNoEncontrado)
	}
}

func TestObtenerUsuario_ValidacionFallida(t *testing.T) {
	repo := &mockUserRepo{
		obtenerPorIDFunc: func(_ context.Context, _ string) (*domain.User, error) {
			t.Error("ObtenerPorID should not be called on validation failure")
			return nil, nil
		},
	}
	auth := &mockAuthSvc{
		verificarFunc: defaultMockAuth,
	}
	uc := NewObtenerUsuarioCasoDeUso(repo, auth)

	tests := []struct {
		name string
		cmd  *ComandoObtenerUsuario
		want error
	}{
		{
			name: "UserID vacío",
			cmd:  &ComandoObtenerUsuario{UserID: "", EjecutorID: "executor"},
			want: ErrUsuarioRequerido,
		},
		{
			name: "EjecutorID vacío",
			cmd:  &ComandoObtenerUsuario{UserID: "target", EjecutorID: ""},
			want: ErrEjecutorRequerido,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := uc.Ejecutar(context.Background(), tt.cmd)

			if resp != nil {
				t.Error("expected nil response on validation error")
			}
			if !errors.Is(err, tt.want) {
				t.Errorf("err = %v, want %v", err, tt.want)
			}
		})
	}
}

func TestObtenerUsuario_ErrorAutorizacion(t *testing.T) {
	repo := &mockUserRepo{
		obtenerPorIDFunc: func(_ context.Context, _ string) (*domain.User, error) {
			t.Error("ObtenerPorID should not be called when auth returns error")
			return nil, nil
		},
	}
	auth := &mockAuthSvc{
		verificarFunc: func(_ context.Context, _, _ string) (bool, error) {
			return false, errors.New("error de infraestructura")
		},
	}
	uc := NewObtenerUsuarioCasoDeUso(repo, auth)

	resp, err := uc.Ejecutar(context.Background(), crearCmdValido())

	if resp != nil {
		t.Error("expected nil response on auth error")
	}
	if err == nil || err.Error() != "error de infraestructura" {
		t.Errorf("err = %v, want %v", err, "error de infraestructura")
	}
}

func TestObtenerUsuario_ErrorRepositorio(t *testing.T) {
	repo := &mockUserRepo{
		obtenerPorIDFunc: func(_ context.Context, _ string) (*domain.User, error) {
			return nil, domain.ErrRepositorio
		},
	}
	auth := &mockAuthSvc{
		verificarFunc: defaultMockAuth,
	}
	uc := NewObtenerUsuarioCasoDeUso(repo, auth)

	resp, err := uc.Ejecutar(context.Background(), crearCmdValido())

	if resp != nil {
		t.Error("expected nil response on repository error")
	}
	if !errors.Is(err, domain.ErrRepositorio) {
		t.Errorf("err = %v, want %v", err, domain.ErrRepositorio)
	}
}

func TestObtenerUsuario_NoLlamaRepoSiNoAutorizado(t *testing.T) {
	obtenerPorIDLlamado := false
	repo := &mockUserRepo{
		obtenerPorIDFunc: func(_ context.Context, _ string) (*domain.User, error) {
			obtenerPorIDLlamado = true
			return nil, nil
		},
	}
	auth := &mockAuthSvc{
		verificarFunc: func(_ context.Context, _, _ string) (bool, error) {
			return false, nil
		},
	}
	uc := NewObtenerUsuarioCasoDeUso(repo, auth)

	resp, _ := uc.Ejecutar(context.Background(), crearCmdValido())

	if resp != nil {
		t.Error("expected nil response on access denied")
	}
	if obtenerPorIDLlamado {
		t.Error("ObtenerPorID should not be called when authorization fails")
	}
}
