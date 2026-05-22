package removerole

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
	removerFunc func(ctx context.Context, userID, roleID string) error
}

func (m *mockUserRoleRepo) Asignar(_ context.Context, _, _ string) error {
	panic("no debería llamarse")
}

func (m *mockUserRoleRepo) Remover(ctx context.Context, userID, roleID string) error {
	return m.removerFunc(ctx, userID, roleID)
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
	panic("no debería llamarse en tests de remove_role")
}

// --- Helpers ---

var (
	rolAdmin   = domain.NuevoRolDesdeBD("role-1", "ADMIN", "Administrador", true)
	rolOWNER   = domain.NuevoRolDesdeBD("role-2", "OWNER", "Propietario", true)
	usuarioVal = domain.NuevoUsuarioDesdeBD("user-1", "user@example.com", "Juan Pérez", "hash", true, nil, "2026-05-21T10:00:00Z")
)

func comandoValido() *ComandoRemoverRol {
	return &ComandoRemoverRol{
		UserID:     "user-1",
		RoleID:     "role-1",
		EjecutorID: "ejecutor-1",
	}
}

func defaultMockAuth(_ context.Context, _ string, _ string) (bool, error) {
	return true, nil
}

// --- Tests ---

func TestRemoverRol_Exito(t *testing.T) {
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
		removerFunc: func(_ context.Context, _, _ string) error {
			return nil
		},
	}
	authSvc := &mockAuthSvc{verificarFunc: defaultMockAuth}
	uc := NewRemoverRolCasoDeUso(userRepo, roleRepo, userRoleRepo, authSvc)

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

func TestRemoverRol_SinPermiso(t *testing.T) {
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
		removerFunc: func(_ context.Context, _, _ string) error {
			return nil
		},
	}
	authSvc := &mockAuthSvc{
		verificarFunc: func(_ context.Context, _, _ string) (bool, error) {
			return false, nil
		},
	}
	uc := NewRemoverRolCasoDeUso(userRepo, roleRepo, userRoleRepo, authSvc)

	resp, err := uc.Ejecutar(context.Background(), comandoValido())
	if resp != nil {
		t.Error("expected nil response on access denied")
	}
	if !errors.Is(err, domain.ErrAccesoDenegado) {
		t.Errorf("err = %v, want %v", err, domain.ErrAccesoDenegado)
	}
}

func TestRemoverRol_ErrorAutorizacion(t *testing.T) {
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
		removerFunc: func(_ context.Context, _, _ string) error {
			return nil
		},
	}
	authSvc := &mockAuthSvc{
		verificarFunc: func(_ context.Context, _, _ string) (bool, error) {
			return false, errors.New("error de infraestructura")
		},
	}
	uc := NewRemoverRolCasoDeUso(userRepo, roleRepo, userRoleRepo, authSvc)

	resp, err := uc.Ejecutar(context.Background(), comandoValido())
	if resp != nil {
		t.Error("expected nil response on auth error")
	}
	if err == nil || err.Error() != "error de infraestructura" {
		t.Errorf("err = %v, want %v", err, "error de infraestructura")
	}
}

func TestRemoverRol_UsuarioNoEncontrado(t *testing.T) {
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
		removerFunc: func(_ context.Context, _, _ string) error {
			return nil
		},
	}
	authSvc := &mockAuthSvc{verificarFunc: defaultMockAuth}
	uc := NewRemoverRolCasoDeUso(userRepo, roleRepo, userRoleRepo, authSvc)

	resp, err := uc.Ejecutar(context.Background(), comandoValido())
	if resp != nil {
		t.Error("expected nil response when user not found")
	}
	if !errors.Is(err, domain.ErrUsuarioNoEncontrado) {
		t.Errorf("err = %v, want %v", err, domain.ErrUsuarioNoEncontrado)
	}
}

func TestRemoverRol_RolNoEncontrado(t *testing.T) {
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
		removerFunc: func(_ context.Context, _, _ string) error {
			return nil
		},
	}
	authSvc := &mockAuthSvc{verificarFunc: defaultMockAuth}
	uc := NewRemoverRolCasoDeUso(userRepo, roleRepo, userRoleRepo, authSvc)

	resp, err := uc.Ejecutar(context.Background(), comandoValido())
	if resp != nil {
		t.Error("expected nil response when role not found")
	}
	if !errors.Is(err, domain.ErrRolNoEncontrado) {
		t.Errorf("err = %v, want %v", err, domain.ErrRolNoEncontrado)
	}
}

func TestRemoverRol_RolOWNERNoRemovible(t *testing.T) {
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
		removerFunc: func(_ context.Context, _, _ string) error {
			return nil
		},
	}
	authSvc := &mockAuthSvc{verificarFunc: defaultMockAuth}
	uc := NewRemoverRolCasoDeUso(userRepo, roleRepo, userRoleRepo, authSvc)

	resp, err := uc.Ejecutar(context.Background(), comandoValido())
	if resp != nil {
		t.Error("expected nil response when role is OWNER")
	}
	if !errors.Is(err, ErrRolOWNERNoRemovible) {
		t.Errorf("err = %v, want %v", err, ErrRolOWNERNoRemovible)
	}
}

func TestRemoverRol_ValidacionFallida(t *testing.T) {
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
		removerFunc: func(_ context.Context, _, _ string) error {
			return nil
		},
	}
	authSvc := &mockAuthSvc{verificarFunc: defaultMockAuth}
	uc := NewRemoverRolCasoDeUso(userRepo, roleRepo, userRoleRepo, authSvc)

	tests := []struct {
		name   string
		modify func(*ComandoRemoverRol)
		want   error
	}{
		{
			name:   "userID vacío",
			modify: func(cmd *ComandoRemoverRol) { cmd.UserID = "" },
			want:   ErrUsuarioRequerido,
		},
		{
			name:   "roleID vacío",
			modify: func(cmd *ComandoRemoverRol) { cmd.RoleID = "" },
			want:   ErrRolRequerido,
		},
		{
			name:   "ejecutorID vacío",
			modify: func(cmd *ComandoRemoverRol) { cmd.EjecutorID = "" },
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

func TestRemoverRol_ErrorRemover(t *testing.T) {
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
		removerFunc: func(_ context.Context, _, _ string) error {
			return domain.ErrRepositorio
		},
	}
	authSvc := &mockAuthSvc{verificarFunc: defaultMockAuth}
	uc := NewRemoverRolCasoDeUso(userRepo, roleRepo, userRoleRepo, authSvc)

	resp, err := uc.Ejecutar(context.Background(), comandoValido())
	if resp != nil {
		t.Error("expected nil response on repository error")
	}
	if !errors.Is(err, domain.ErrRepositorio) {
		t.Errorf("err = %v, want %v", err, domain.ErrRepositorio)
	}
}

func TestRemoverRol_NoLlamaRoleRepoSiNoAutorizado(t *testing.T) {
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
		removerFunc: func(_ context.Context, _, _ string) error {
			return nil
		},
	}
	authSvc := &mockAuthSvc{
		verificarFunc: func(_ context.Context, _, _ string) (bool, error) {
			return false, nil
		},
	}
	uc := NewRemoverRolCasoDeUso(userRepo, roleRepo, userRoleRepo, authSvc)

	resp, _ := uc.Ejecutar(context.Background(), comandoValido())
	if resp != nil {
		t.Error("expected nil response on access denied")
	}
	if roleRepoCalled {
		t.Error("role repository should not be called when authorization fails")
	}
}

func TestRemoverRol_NoRemueveSiRolEsOWNER(t *testing.T) {
	removerCalled := false
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
		removerFunc: func(_ context.Context, _, _ string) error {
			removerCalled = true
			return nil
		},
	}
	authSvc := &mockAuthSvc{verificarFunc: defaultMockAuth}
	uc := NewRemoverRolCasoDeUso(userRepo, roleRepo, userRoleRepo, authSvc)

	resp, err := uc.Ejecutar(context.Background(), comandoValido())
	if resp != nil {
		t.Error("expected nil response when role is OWNER")
	}
	if !errors.Is(err, ErrRolOWNERNoRemovible) {
		t.Errorf("err = %v, want %v", err, ErrRolOWNERNoRemovible)
	}
	if removerCalled {
		t.Error("Remover should not be called when role is OWNER")
	}
}
