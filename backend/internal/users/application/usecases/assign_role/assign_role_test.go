package assignrole

import (
	"context"
	"errors"
	"testing"

	"github.com/Alexisjar91/POS/internal/users/domain"
	"github.com/Alexisjar91/POS/pkg/especificacion"
	"github.com/Alexisjar91/POS/pkg/paginacion"
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
	return m.obtenerPorIDFunc(ctx, id)
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

// mockRoleRepo implementa domain.RoleRepository para tests.
type mockRoleRepo struct {
	obtenerPorIDFunc func(ctx context.Context, id string) (*domain.Role, error)
}

func (m *mockRoleRepo) Crear(_ context.Context, _ *domain.Role) (*domain.Role, error) {
	panic("no debería llamarse")
}

func (m *mockRoleRepo) ObtenerPorID(ctx context.Context, id string) (*domain.Role, error) {
	return m.obtenerPorIDFunc(ctx, id)
}

func (m *mockRoleRepo) ObtenerPorNombre(_ context.Context, _ string) (*domain.Role, error) {
	panic("no debería llamarse")
}

func (m *mockRoleRepo) Actualizar(_ context.Context, _ *domain.Role) (*domain.Role, error) {
	panic("no debería llamarse")
}

func (m *mockRoleRepo) Eliminar(_ context.Context, _ string) error {
	panic("no debería llamarse")
}

func (m *mockRoleRepo) Listar(_ context.Context) ([]*domain.Role, error) {
	panic("no debería llamarse")
}

// mockUserRoleRepo implementa domain.UserRoleRepository para tests.
type mockUserRoleRepo struct {
	asignarFunc func(ctx context.Context, userID, roleID string) error
}

func (m *mockUserRoleRepo) Asignar(ctx context.Context, userID, roleID string) error {
	return m.asignarFunc(ctx, userID, roleID)
}

func (m *mockUserRoleRepo) Remover(_ context.Context, _, _ string) error {
	panic("no debería llamarse en tests de assign_role")
}

// mockAuthSvc implementa domain.AuthorizationService para tests.
type mockAuthSvc struct {
	verificarFunc func(ctx context.Context, userID, permissionCode string) (bool, error)
}

func (m *mockAuthSvc) VerificarPermiso(ctx context.Context, userID, permissionCode string) (bool, error) {
	if m.verificarFunc != nil {
		return m.verificarFunc(ctx, userID, permissionCode)
	}
	return false, nil
}

func (m *mockAuthSvc) EsUsuarioOWNER(_ context.Context, _ string) (bool, error) {
	panic("no debería llamarse en tests de assign_role")
}

// --- Helpers ---

var (
	rolAdmin     = domain.NuevoRolDesdeBD("role-1", "ADMIN", "Administrador", true)
	rolOWNER     = domain.NuevoRolDesdeBD("role-2", "OWNER", "Propietario", true)
	usuarioVal   = domain.NuevoUsuarioDesdeBD("user-1", "user@example.com", "Juan Pérez", "hash", true, nil, "2026-05-21T10:00:00Z")
)

func comandoValido() *ComandoAsignarRol {
	return &ComandoAsignarRol{
		UserID:     "user-1",
		RoleID:     "role-1",
		EjecutorID: "ejecutor-1",
	}
}

func defaultMockAuth(_ context.Context, _ string, _ string) (bool, error) {
	return true, nil
}

// --- Tests ---

func TestAsignarRol_Exito(t *testing.T) {
	userRepo := &mockUserRepo{
		obtenerPorIDFunc: func(_ context.Context, _ string) (*domain.User, error) {
			return usuarioVal, nil
		},
	}
	roleRepo := &mockRoleRepo{
		obtenerPorIDFunc: func(_ context.Context, _ string) (*domain.Role, error) {
			return rolAdmin, nil
		},
	}
	userRoleRepo := &mockUserRoleRepo{
		asignarFunc: func(_ context.Context, _, _ string) error {
			return nil
		},
	}
	authSvc := &mockAuthSvc{verificarFunc: defaultMockAuth}
	uc := NewAsignarRolCasoDeUso(userRepo, roleRepo, userRoleRepo, authSvc)

	resp, err := uc.Ejecutar(context.Background(), comandoValido())
	if err != nil {
		t.Fatalf("Ejecutar returned unexpected error: %v", err)
	}
	if resp == nil {
		t.Fatal("Ejecutar returned nil response")
	}
	if resp.UserID != "user-1" {
		t.Errorf("resp.UserID = %q, want %q", resp.UserID, "user-1")
	}
	if resp.RoleID != "role-1" {
		t.Errorf("resp.RoleID = %q, want %q", resp.RoleID, "role-1")
	}
}

func TestAsignarRol_SinPermiso(t *testing.T) {
	userRepo := &mockUserRepo{
		obtenerPorIDFunc: func(_ context.Context, _ string) (*domain.User, error) {
			return usuarioVal, nil
		},
	}
	roleRepo := &mockRoleRepo{
		obtenerPorIDFunc: func(_ context.Context, _ string) (*domain.Role, error) {
			return rolAdmin, nil
		},
	}
	userRoleRepo := &mockUserRoleRepo{
		asignarFunc: func(_ context.Context, _, _ string) error {
			return nil
		},
	}
	authSvc := &mockAuthSvc{
		verificarFunc: func(_ context.Context, _, _ string) (bool, error) {
			return false, nil
		},
	}
	uc := NewAsignarRolCasoDeUso(userRepo, roleRepo, userRoleRepo, authSvc)

	resp, err := uc.Ejecutar(context.Background(), comandoValido())
	if resp != nil {
		t.Error("expected nil response on access denied")
	}
	if !errors.Is(err, domain.ErrAccesoDenegado) {
		t.Errorf("err = %v, want %v", err, domain.ErrAccesoDenegado)
	}
}

func TestAsignarRol_ErrorAutorizacion(t *testing.T) {
	userRepo := &mockUserRepo{
		obtenerPorIDFunc: func(_ context.Context, _ string) (*domain.User, error) {
			return usuarioVal, nil
		},
	}
	roleRepo := &mockRoleRepo{
		obtenerPorIDFunc: func(_ context.Context, _ string) (*domain.Role, error) {
			return rolAdmin, nil
		},
	}
	userRoleRepo := &mockUserRoleRepo{
		asignarFunc: func(_ context.Context, _, _ string) error {
			return nil
		},
	}
	authSvc := &mockAuthSvc{
		verificarFunc: func(_ context.Context, _, _ string) (bool, error) {
			return false, errors.New("error de infraestructura")
		},
	}
	uc := NewAsignarRolCasoDeUso(userRepo, roleRepo, userRoleRepo, authSvc)

	resp, err := uc.Ejecutar(context.Background(), comandoValido())
	if resp != nil {
		t.Error("expected nil response on auth error")
	}
	if err == nil || err.Error() != "error de infraestructura" {
		t.Errorf("err = %v, want %v", err, "error de infraestructura")
	}
}

func TestAsignarRol_UsuarioNoEncontrado(t *testing.T) {
	userRepo := &mockUserRepo{
		obtenerPorIDFunc: func(_ context.Context, _ string) (*domain.User, error) {
			return nil, domain.ErrUsuarioNoEncontrado
		},
	}
	roleRepo := &mockRoleRepo{
		obtenerPorIDFunc: func(_ context.Context, _ string) (*domain.Role, error) {
			return rolAdmin, nil
		},
	}
	userRoleRepo := &mockUserRoleRepo{
		asignarFunc: func(_ context.Context, _, _ string) error {
			return nil
		},
	}
	authSvc := &mockAuthSvc{verificarFunc: defaultMockAuth}
	uc := NewAsignarRolCasoDeUso(userRepo, roleRepo, userRoleRepo, authSvc)

	resp, err := uc.Ejecutar(context.Background(), comandoValido())
	if resp != nil {
		t.Error("expected nil response when user not found")
	}
	if !errors.Is(err, domain.ErrUsuarioNoEncontrado) {
		t.Errorf("err = %v, want %v", err, domain.ErrUsuarioNoEncontrado)
	}
}

func TestAsignarRol_RolNoEncontrado(t *testing.T) {
	userRepo := &mockUserRepo{
		obtenerPorIDFunc: func(_ context.Context, _ string) (*domain.User, error) {
			return usuarioVal, nil
		},
	}
	roleRepo := &mockRoleRepo{
		obtenerPorIDFunc: func(_ context.Context, _ string) (*domain.Role, error) {
			return nil, domain.ErrRolNoEncontrado
		},
	}
	userRoleRepo := &mockUserRoleRepo{
		asignarFunc: func(_ context.Context, _, _ string) error {
			return nil
		},
	}
	authSvc := &mockAuthSvc{verificarFunc: defaultMockAuth}
	uc := NewAsignarRolCasoDeUso(userRepo, roleRepo, userRoleRepo, authSvc)

	resp, err := uc.Ejecutar(context.Background(), comandoValido())
	if resp != nil {
		t.Error("expected nil response when role not found")
	}
	if !errors.Is(err, domain.ErrRolNoEncontrado) {
		t.Errorf("err = %v, want %v", err, domain.ErrRolNoEncontrado)
	}
}

func TestAsignarRol_RolOWNERNoAsignable(t *testing.T) {
	userRepo := &mockUserRepo{
		obtenerPorIDFunc: func(_ context.Context, _ string) (*domain.User, error) {
			return usuarioVal, nil
		},
	}
	roleRepo := &mockRoleRepo{
		obtenerPorIDFunc: func(_ context.Context, _ string) (*domain.Role, error) {
			return rolOWNER, nil
		},
	}
	userRoleRepo := &mockUserRoleRepo{
		asignarFunc: func(_ context.Context, _, _ string) error {
			return nil
		},
	}
	authSvc := &mockAuthSvc{verificarFunc: defaultMockAuth}
	uc := NewAsignarRolCasoDeUso(userRepo, roleRepo, userRoleRepo, authSvc)

	resp, err := uc.Ejecutar(context.Background(), comandoValido())
	if resp != nil {
		t.Error("expected nil response when role is OWNER")
	}
	if !errors.Is(err, ErrRolOWNERNoAsignable) {
		t.Errorf("err = %v, want %v", err, ErrRolOWNERNoAsignable)
	}
}

func TestAsignarRol_ValidacionFallida(t *testing.T) {
	userRepo := &mockUserRepo{
		obtenerPorIDFunc: func(_ context.Context, _ string) (*domain.User, error) {
			return usuarioVal, nil
		},
	}
	roleRepo := &mockRoleRepo{
		obtenerPorIDFunc: func(_ context.Context, _ string) (*domain.Role, error) {
			return rolAdmin, nil
		},
	}
	userRoleRepo := &mockUserRoleRepo{
		asignarFunc: func(_ context.Context, _, _ string) error {
			return nil
		},
	}
	authSvc := &mockAuthSvc{verificarFunc: defaultMockAuth}
	uc := NewAsignarRolCasoDeUso(userRepo, roleRepo, userRoleRepo, authSvc)

	tests := []struct {
		name   string
		modify func(*ComandoAsignarRol)
		want   error
	}{
		{
			name:   "userID vacío",
			modify: func(cmd *ComandoAsignarRol) { cmd.UserID = "" },
			want:   ErrUsuarioRequerido,
		},
		{
			name:   "roleID vacío",
			modify: func(cmd *ComandoAsignarRol) { cmd.RoleID = "" },
			want:   ErrRolRequerido,
		},
		{
			name:   "ejecutorID vacío",
			modify: func(cmd *ComandoAsignarRol) { cmd.EjecutorID = "" },
			want:   ErrEjecutorRequerido,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := comandoValido()
			tt.modify(cmd)

			resp, err := uc.Ejecutar(context.Background(), cmd)
			if resp != nil {
				t.Error("expected nil response on validation error")
			}
			if !errors.Is(err, tt.want) {
				t.Errorf("err = %v, want %v", err, tt.want)
			}
		})
	}
}

func TestAsignarRol_ErrorAsignar(t *testing.T) {
	userRepo := &mockUserRepo{
		obtenerPorIDFunc: func(_ context.Context, _ string) (*domain.User, error) {
			return usuarioVal, nil
		},
	}
	roleRepo := &mockRoleRepo{
		obtenerPorIDFunc: func(_ context.Context, _ string) (*domain.Role, error) {
			return rolAdmin, nil
		},
	}
	userRoleRepo := &mockUserRoleRepo{
		asignarFunc: func(_ context.Context, _, _ string) error {
			return domain.ErrRepositorio
		},
	}
	authSvc := &mockAuthSvc{verificarFunc: defaultMockAuth}
	uc := NewAsignarRolCasoDeUso(userRepo, roleRepo, userRoleRepo, authSvc)

	resp, err := uc.Ejecutar(context.Background(), comandoValido())
	if resp != nil {
		t.Error("expected nil response on repository error")
	}
	if !errors.Is(err, domain.ErrRepositorio) {
		t.Errorf("err = %v, want %v", err, domain.ErrRepositorio)
	}
}

func TestAsignarRol_NoLlamaRoleRepoSiNoAutorizado(t *testing.T) {
	roleRepoCalled := false
	userRepo := &mockUserRepo{
		obtenerPorIDFunc: func(_ context.Context, _ string) (*domain.User, error) {
			return usuarioVal, nil
		},
	}
	roleRepo := &mockRoleRepo{
		obtenerPorIDFunc: func(_ context.Context, _ string) (*domain.Role, error) {
			roleRepoCalled = true
			return rolAdmin, nil
		},
	}
	userRoleRepo := &mockUserRoleRepo{
		asignarFunc: func(_ context.Context, _, _ string) error {
			return nil
		},
	}
	authSvc := &mockAuthSvc{
		verificarFunc: func(_ context.Context, _, _ string) (bool, error) {
			return false, nil
		},
	}
	uc := NewAsignarRolCasoDeUso(userRepo, roleRepo, userRoleRepo, authSvc)

	resp, _ := uc.Ejecutar(context.Background(), comandoValido())
	if resp != nil {
		t.Error("expected nil response on access denied")
	}
	if roleRepoCalled {
		t.Error("role repository should not be called when authorization fails")
	}
}

func TestAsignarRol_NoAsignaSiRolEsOWNER(t *testing.T) {
	asignarCalled := false
	userRepo := &mockUserRepo{
		obtenerPorIDFunc: func(_ context.Context, _ string) (*domain.User, error) {
			return usuarioVal, nil
		},
	}
	roleRepo := &mockRoleRepo{
		obtenerPorIDFunc: func(_ context.Context, _ string) (*domain.Role, error) {
			return rolOWNER, nil
		},
	}
	userRoleRepo := &mockUserRoleRepo{
		asignarFunc: func(_ context.Context, _, _ string) error {
			asignarCalled = true
			return nil
		},
	}
	authSvc := &mockAuthSvc{verificarFunc: defaultMockAuth}
	uc := NewAsignarRolCasoDeUso(userRepo, roleRepo, userRoleRepo, authSvc)

	resp, err := uc.Ejecutar(context.Background(), comandoValido())
	if resp != nil {
		t.Error("expected nil response when role is OWNER")
	}
	if !errors.Is(err, ErrRolOWNERNoAsignable) {
		t.Errorf("err = %v, want %v", err, ErrRolOWNERNoAsignable)
	}
	if asignarCalled {
		t.Error("Asignar should not be called when role is OWNER")
	}
}
