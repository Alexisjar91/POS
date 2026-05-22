package disableuser

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
	actualizarFunc   func(ctx context.Context, user *domain.User) (*domain.User, error)
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

func (m *mockUserRepo) Actualizar(ctx context.Context, user *domain.User) (*domain.User, error) {
	if m.actualizarFunc != nil {
		return m.actualizarFunc(ctx, user)
	}
	return nil, nil
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

func strPtr(s string) *string {
	return &s
}

func crearCmdValido() *ComandoDesactivarUsuario {
	return &ComandoDesactivarUsuario{
		UserID:     "user-target",
		EjecutorID: "user-executor",
	}
}

func usuarioActivo() *domain.User {
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

func usuarioInactivo() *domain.User {
	return domain.NuevoUsuarioDesdeBD(
		"user-target",
		"target@test.com",
		"Target User",
		"hash",
		false,
		strPtr("creator"),
		"2024-01-01",
	)
}

func defaultMockAuth(_ context.Context, _ string, _ string) (bool, error) {
	return true, nil
}

func defaultMockOwner(_ context.Context, _ string) (bool, error) {
	return false, nil
}

func defaultMockActualizar(_ context.Context, user *domain.User) (*domain.User, error) {
	return user, nil
}

// --- Tests ---

func TestDesactivarUsuario_Exito(t *testing.T) {
	repo := &mockUserRepo{
		obtenerPorIDFunc: func(_ context.Context, _ string) (*domain.User, error) {
			return usuarioActivo(), nil
		},
		actualizarFunc: defaultMockActualizar,
	}
	auth := &mockAuthSvc{
		verificarFunc: defaultMockAuth,
		esOwnerFunc:   defaultMockOwner,
	}
	uc := NewDesactivarUsuarioCasoDeUso(repo, auth)

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
	if resp.Active {
		t.Error("resp.Active should be false after disable")
	}
}

func TestDesactivarUsuario_SinPermiso(t *testing.T) {
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
		esOwnerFunc: defaultMockOwner,
	}
	uc := NewDesactivarUsuarioCasoDeUso(repo, auth)

	resp, err := uc.Ejecutar(context.Background(), crearCmdValido())

	if resp != nil {
		t.Error("expected nil response on access denied")
	}
	if !errors.Is(err, domain.ErrAccesoDenegado) {
		t.Errorf("err = %v, want %v", err, domain.ErrAccesoDenegado)
	}
}

func TestDesactivarUsuario_AutoDesactivacion(t *testing.T) {
	repo := &mockUserRepo{
		obtenerPorIDFunc: func(_ context.Context, _ string) (*domain.User, error) {
			return usuarioActivo(), nil
		},
	}
	auth := &mockAuthSvc{
		verificarFunc: defaultMockAuth,
		esOwnerFunc:   defaultMockOwner,
	}
	uc := NewDesactivarUsuarioCasoDeUso(repo, auth)

	cmd := &ComandoDesactivarUsuario{
		UserID:     "same-user",
		EjecutorID: "same-user",
	}

	resp, err := uc.Ejecutar(context.Background(), cmd)

	if resp != nil {
		t.Error("expected nil response on self-deactivation")
	}
	if !errors.Is(err, ErrAutoDesactivacion) {
		t.Errorf("err = %v, want %v", err, ErrAutoDesactivacion)
	}
}

func TestDesactivarUsuario_OWNERInmune(t *testing.T) {
	repo := &mockUserRepo{
		obtenerPorIDFunc: func(_ context.Context, _ string) (*domain.User, error) {
			return usuarioActivo(), nil
		},
	}
	auth := &mockAuthSvc{
		verificarFunc: defaultMockAuth,
		esOwnerFunc: func(_ context.Context, userID string) (bool, error) {
			return true, nil
		},
	}
	uc := NewDesactivarUsuarioCasoDeUso(repo, auth)

	resp, err := uc.Ejecutar(context.Background(), crearCmdValido())

	if resp != nil {
		t.Error("expected nil response when target is OWNER")
	}
	if !errors.Is(err, ErrOWNERInmune) {
		t.Errorf("err = %v, want %v", err, ErrOWNERInmune)
	}
}

func TestDesactivarUsuario_UsuarioNoEncontrado(t *testing.T) {
	repo := &mockUserRepo{
		obtenerPorIDFunc: func(_ context.Context, _ string) (*domain.User, error) {
			return nil, domain.ErrUsuarioNoEncontrado
		},
	}
	auth := &mockAuthSvc{
		verificarFunc: defaultMockAuth,
		esOwnerFunc:   defaultMockOwner,
	}
	uc := NewDesactivarUsuarioCasoDeUso(repo, auth)

	resp, err := uc.Ejecutar(context.Background(), crearCmdValido())

	if resp != nil {
		t.Error("expected nil response on user not found")
	}
	if !errors.Is(err, domain.ErrUsuarioNoEncontrado) {
		t.Errorf("err = %v, want %v", err, domain.ErrUsuarioNoEncontrado)
	}
}

func TestDesactivarUsuario_ValidacionFallida(t *testing.T) {
	repo := &mockUserRepo{
		obtenerPorIDFunc: func(_ context.Context, _ string) (*domain.User, error) {
			t.Error("ObtenerPorID should not be called on validation failure")
			return nil, nil
		},
	}
	auth := &mockAuthSvc{
		verificarFunc: defaultMockAuth,
		esOwnerFunc:   defaultMockOwner,
	}
	uc := NewDesactivarUsuarioCasoDeUso(repo, auth)

	tests := []struct {
		name string
		cmd  *ComandoDesactivarUsuario
		want error
	}{
		{
			name: "UserID vacío",
			cmd:  &ComandoDesactivarUsuario{UserID: "", EjecutorID: "executor"},
			want: ErrUsuarioRequerido,
		},
		{
			name: "EjecutorID vacío",
			cmd:  &ComandoDesactivarUsuario{UserID: "target", EjecutorID: ""},
			want: ErrEjecutorRequerido,
		},
		{
			name: "ambos vacíos",
			cmd:  &ComandoDesactivarUsuario{UserID: "", EjecutorID: ""},
			want: ErrUsuarioRequerido, // primer error
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

func TestDesactivarUsuario_NoLlamaRepoSiNoAutorizado(t *testing.T) {
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
		esOwnerFunc: defaultMockOwner,
	}
	uc := NewDesactivarUsuarioCasoDeUso(repo, auth)

	resp, _ := uc.Ejecutar(context.Background(), crearCmdValido())

	if resp != nil {
		t.Error("expected nil response on access denied")
	}
	if obtenerPorIDLlamado {
		t.Error("ObtenerPorID should not be called when authorization fails")
	}
}

func TestDesactivarUsuario_NoLLamaActualizarSiFallaOWNERCheck(t *testing.T) {
	actualizarLlamado := false
	repo := &mockUserRepo{
		obtenerPorIDFunc: func(_ context.Context, _ string) (*domain.User, error) {
			return usuarioActivo(), nil
		},
		actualizarFunc: func(_ context.Context, _ *domain.User) (*domain.User, error) {
			actualizarLlamado = true
			return nil, nil
		},
	}
	auth := &mockAuthSvc{
		verificarFunc: defaultMockAuth,
		esOwnerFunc: func(_ context.Context, _ string) (bool, error) {
			return false, errors.New("error de infraestructura")
		},
	}
	uc := NewDesactivarUsuarioCasoDeUso(repo, auth)

	resp, err := uc.Ejecutar(context.Background(), crearCmdValido())

	if resp != nil {
		t.Error("expected nil response on owner check error")
	}
	if err == nil || err.Error() != "error de infraestructura" {
		t.Errorf("err = %v, want %v", err, "error de infraestructura")
	}
	if actualizarLlamado {
		t.Error("Actualizar should not be called when EsUsuarioOWNER fails")
	}
}

func TestDesactivarUsuario_ErrorAutorizacion(t *testing.T) {
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
		esOwnerFunc: defaultMockOwner,
	}
	uc := NewDesactivarUsuarioCasoDeUso(repo, auth)

	resp, err := uc.Ejecutar(context.Background(), crearCmdValido())

	if resp != nil {
		t.Error("expected nil response on auth error")
	}
	if err == nil || err.Error() != "error de infraestructura" {
		t.Errorf("err = %v, want %v", err, "error de infraestructura")
	}
}

func TestDesactivarUsuario_ErrorRepositorio(t *testing.T) {
	repo := &mockUserRepo{
		obtenerPorIDFunc: func(_ context.Context, _ string) (*domain.User, error) {
			return nil, domain.ErrRepositorio
		},
	}
	auth := &mockAuthSvc{
		verificarFunc: defaultMockAuth,
		esOwnerFunc:   defaultMockOwner,
	}
	uc := NewDesactivarUsuarioCasoDeUso(repo, auth)

	resp, err := uc.Ejecutar(context.Background(), crearCmdValido())

	if resp != nil {
		t.Error("expected nil response on repository error")
	}
	if !errors.Is(err, domain.ErrRepositorio) {
		t.Errorf("err = %v, want %v", err, domain.ErrRepositorio)
	}
}

func TestDesactivarUsuario_UsuarioYaInactivo(t *testing.T) {
	repo := &mockUserRepo{
		obtenerPorIDFunc: func(_ context.Context, _ string) (*domain.User, error) {
			return usuarioInactivo(), nil
		},
	}
	auth := &mockAuthSvc{
		verificarFunc: defaultMockAuth,
		esOwnerFunc:   defaultMockOwner,
	}
	uc := NewDesactivarUsuarioCasoDeUso(repo, auth)

	resp, err := uc.Ejecutar(context.Background(), crearCmdValido())

	if resp != nil {
		t.Error("expected nil response on already inactive user")
	}
	if !errors.Is(err, domain.ErrUsuarioYaInactivo) {
		t.Errorf("err = %v, want %v", err, domain.ErrUsuarioYaInactivo)
	}
}

func TestDesactivarUsuario_ErrorActualizar(t *testing.T) {
	repo := &mockUserRepo{
		obtenerPorIDFunc: func(_ context.Context, _ string) (*domain.User, error) {
			return usuarioActivo(), nil
		},
		actualizarFunc: func(_ context.Context, _ *domain.User) (*domain.User, error) {
			return nil, domain.ErrRepositorio
		},
	}
	auth := &mockAuthSvc{
		verificarFunc: defaultMockAuth,
		esOwnerFunc:   defaultMockOwner,
	}
	uc := NewDesactivarUsuarioCasoDeUso(repo, auth)

	resp, err := uc.Ejecutar(context.Background(), crearCmdValido())

	if resp != nil {
		t.Error("expected nil response on update error")
	}
	if !errors.Is(err, domain.ErrRepositorio) {
		t.Errorf("err = %v, want %v", err, domain.ErrRepositorio)
	}
}
