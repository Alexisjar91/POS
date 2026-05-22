package manageroles

import (
	"context"
	"errors"
	"testing"

	"github.com/Alexisjar91/POS/internal/users/domain"
)

// --- Mocks ---

// mockRoleRepo implementa domain.RoleRepository para tests.
type mockRoleRepo struct {
	crearFunc        func(ctx context.Context, role *domain.Role) (*domain.Role, error)
	obtenerPorIDFunc func(ctx context.Context, id string) (*domain.Role, error)
	obtenerPorNombreFunc func(ctx context.Context, name string) (*domain.Role, error)
	actualizarFunc   func(ctx context.Context, role *domain.Role) (*domain.Role, error)
	eliminarFunc     func(ctx context.Context, id string) error
	listarFunc       func(ctx context.Context) ([]*domain.Role, error)
}

func (m *mockRoleRepo) Crear(ctx context.Context, role *domain.Role) (*domain.Role, error) {
	if m.crearFunc != nil {
		return m.crearFunc(ctx, role)
	}
	panic("no debería llamarse")
}

func (m *mockRoleRepo) ObtenerPorID(ctx context.Context, id string) (*domain.Role, error) {
	if m.obtenerPorIDFunc != nil {
		return m.obtenerPorIDFunc(ctx, id)
	}
	panic("no debería llamarse")
}

func (m *mockRoleRepo) ObtenerPorNombre(ctx context.Context, name string) (*domain.Role, error) {
	if m.obtenerPorNombreFunc != nil {
		return m.obtenerPorNombreFunc(ctx, name)
	}
	panic("no debería llamarse")
}

func (m *mockRoleRepo) Actualizar(ctx context.Context, role *domain.Role) (*domain.Role, error) {
	if m.actualizarFunc != nil {
		return m.actualizarFunc(ctx, role)
	}
	panic("no debería llamarse")
}

func (m *mockRoleRepo) Eliminar(ctx context.Context, id string) error {
	if m.eliminarFunc != nil {
		return m.eliminarFunc(ctx, id)
	}
	panic("no debería llamarse")
}

func (m *mockRoleRepo) Listar(ctx context.Context) ([]*domain.Role, error) {
	if m.listarFunc != nil {
		return m.listarFunc(ctx)
	}
	panic("no debería llamarse")
}

// mockAuthSvc implementa domain.AuthorizationService para tests.
type mockAuthSvc struct {
	verificarFunc func(ctx context.Context, userID, permissionCode string) (bool, error)
	esOwnerFunc   func(ctx context.Context, userID string) (bool, error)
}

func (m *mockAuthSvc) VerificarPermiso(ctx context.Context, userID, permissionCode string) (bool, error) {
	if m.verificarFunc != nil {
		return m.verificarFunc(ctx, userID, permissionCode)
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

var (
	rolNoSistema = domain.NuevoRolDesdeBD("role-1", "cashier", "Cajero", false)
	rolSistema   = domain.NuevoRolDesdeBD("role-sys", "ADMIN", "Administrador", true)
)

func defaultMockAuth(_ context.Context, _ string, _ string) (bool, error) {
	return true, nil
}

// ============================================================
// Tests para CrearRol
// ============================================================

func TestCrearRol_Exito(t *testing.T) {
	repo := &mockRoleRepo{
		crearFunc: func(_ context.Context, role *domain.Role) (*domain.Role, error) {
			return domain.NuevoRolDesdeBD("new-role-id", role.Name(), "Cajero del sistema", false), nil
		},
	}
	auth := &mockAuthSvc{verificarFunc: defaultMockAuth}
	uc := NewCrearRolCasoDeUso(repo, auth)

	cmd := &ComandoCrearRol{
		Nombre:      "cashier",
		Descripcion: "Cajero del sistema",
		EjecutorID:  "ejecutor-1",
	}

	resp, err := uc.Ejecutar(context.Background(), cmd)
	if err != nil {
		t.Fatalf("Ejecutar returned unexpected error: %v", err)
	}
	if resp == nil {
		t.Fatal("Ejecutar returned nil response")
	}
	if resp.ID != "new-role-id" {
		t.Errorf("resp.ID = %q, want %q", resp.ID, "new-role-id")
	}
	if resp.Nombre != "cashier" {
		t.Errorf("resp.Nombre = %q, want %q", resp.Nombre, "cashier")
	}
	if resp.Descripcion != "Cajero del sistema" {
		t.Errorf("resp.Descripcion = %q, want %q", resp.Descripcion, "Cajero del sistema")
	}
	if resp.IsSystem {
		t.Error("resp.IsSystem should be false for a new role")
	}
}

func TestCrearRol_SinPermiso(t *testing.T) {
	repo := &mockRoleRepo{
		crearFunc: func(_ context.Context, _ *domain.Role) (*domain.Role, error) {
			t.Error("Crear should not be called when not authorized")
			return nil, nil
		},
	}
	auth := &mockAuthSvc{
		verificarFunc: func(_ context.Context, _, _ string) (bool, error) {
			return false, nil
		},
	}
	uc := NewCrearRolCasoDeUso(repo, auth)

	cmd := &ComandoCrearRol{
		Nombre:      "cashier",
		Descripcion: "Cajero",
		EjecutorID:  "ejecutor-1",
	}

	resp, err := uc.Ejecutar(context.Background(), cmd)
	if resp != nil {
		t.Error("expected nil response on access denied")
	}
	if !errors.Is(err, domain.ErrAccesoDenegado) {
		t.Errorf("err = %v, want %v", err, domain.ErrAccesoDenegado)
	}
}

func TestCrearRol_ValidacionFallida(t *testing.T) {
	repo := &mockRoleRepo{
		crearFunc: func(_ context.Context, _ *domain.Role) (*domain.Role, error) {
			t.Error("Crear should not be called on validation error")
			return nil, nil
		},
	}
	auth := &mockAuthSvc{verificarFunc: defaultMockAuth}
	uc := NewCrearRolCasoDeUso(repo, auth)

	tests := []struct {
		name   string
		modify func(*ComandoCrearRol)
		want   error
	}{
		{
			name:   "nombre vacío",
			modify: func(cmd *ComandoCrearRol) { cmd.Nombre = "" },
			want:   ErrNombreRequerido,
		},
		{
			name:   "ejecutorID vacío",
			modify: func(cmd *ComandoCrearRol) { cmd.EjecutorID = "" },
			want:   ErrEjecutorRequerido,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &ComandoCrearRol{
				Nombre:      "cashier",
				Descripcion: "Cajero",
				EjecutorID:  "ejecutor-1",
			}
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

func TestCrearRol_ErrorRepositorio(t *testing.T) {
	repo := &mockRoleRepo{
		crearFunc: func(_ context.Context, _ *domain.Role) (*domain.Role, error) {
			return nil, domain.ErrRolDuplicado
		},
	}
	auth := &mockAuthSvc{verificarFunc: defaultMockAuth}
	uc := NewCrearRolCasoDeUso(repo, auth)

	cmd := &ComandoCrearRol{
		Nombre:      "cashier",
		Descripcion: "Cajero",
		EjecutorID:  "ejecutor-1",
	}

	resp, err := uc.Ejecutar(context.Background(), cmd)
	if resp != nil {
		t.Error("expected nil response on repository error")
	}
	if !errors.Is(err, domain.ErrRolDuplicado) {
		t.Errorf("err = %v, want %v", err, domain.ErrRolDuplicado)
	}
}

func TestCrearRol_NoLlamaRepoSiNoAutorizado(t *testing.T) {
	llamado := false
	repo := &mockRoleRepo{
		crearFunc: func(_ context.Context, _ *domain.Role) (*domain.Role, error) {
			llamado = true
			return nil, nil
		},
	}
	auth := &mockAuthSvc{
		verificarFunc: func(_ context.Context, _, _ string) (bool, error) {
			return false, nil
		},
	}
	uc := NewCrearRolCasoDeUso(repo, auth)

	cmd := &ComandoCrearRol{
		Nombre:      "cashier",
		Descripcion: "Cajero",
		EjecutorID:  "ejecutor-1",
	}

	resp, _ := uc.Ejecutar(context.Background(), cmd)
	if resp != nil {
		t.Error("expected nil response on access denied")
	}
	if llamado {
		t.Error("Crear should not be called when authorization fails")
	}
}

// ============================================================
// Tests para ActualizarRol
// ============================================================

func TestActualizarRol_Exito(t *testing.T) {
	repo := &mockRoleRepo{
		obtenerPorIDFunc: func(_ context.Context, _ string) (*domain.Role, error) {
			return rolNoSistema, nil
		},
		actualizarFunc: func(_ context.Context, role *domain.Role) (*domain.Role, error) {
			return role, nil
		},
	}
	auth := &mockAuthSvc{verificarFunc: defaultMockAuth}
	uc := NewActualizarRolCasoDeUso(repo, auth)

	cmd := &ComandoActualizarRol{
		RoleID:      "role-1",
		Nombre:      "supervisor",
		Descripcion: "Supervisor de caja",
		EjecutorID:  "ejecutor-1",
	}

	resp, err := uc.Ejecutar(context.Background(), cmd)
	if err != nil {
		t.Fatalf("Ejecutar returned unexpected error: %v", err)
	}
	if resp == nil {
		t.Fatal("Ejecutar returned nil response")
	}
	if resp.ID != "role-1" {
		t.Errorf("resp.ID = %q, want %q", resp.ID, "role-1")
	}
	if resp.Nombre != "supervisor" {
		t.Errorf("resp.Nombre = %q, want %q", resp.Nombre, "supervisor")
	}
	if resp.Descripcion != "Supervisor de caja" {
		t.Errorf("resp.Descripcion = %q, want %q", resp.Descripcion, "Supervisor de caja")
	}
	if resp.IsSystem {
		t.Error("resp.IsSystem should be false")
	}
}

func TestActualizarRol_SinPermiso(t *testing.T) {
	repo := &mockRoleRepo{
		obtenerPorIDFunc: func(_ context.Context, _ string) (*domain.Role, error) {
			t.Error("ObtenerPorID should not be called when not authorized")
			return nil, nil
		},
	}
	auth := &mockAuthSvc{
		verificarFunc: func(_ context.Context, _, _ string) (bool, error) {
			return false, nil
		},
	}
	uc := NewActualizarRolCasoDeUso(repo, auth)

	cmd := &ComandoActualizarRol{
		RoleID:      "role-1",
		Nombre:      "supervisor",
		Descripcion: "Supervisor",
		EjecutorID:  "ejecutor-1",
	}

	resp, err := uc.Ejecutar(context.Background(), cmd)
	if resp != nil {
		t.Error("expected nil response on access denied")
	}
	if !errors.Is(err, domain.ErrAccesoDenegado) {
		t.Errorf("err = %v, want %v", err, domain.ErrAccesoDenegado)
	}
}

func TestActualizarRol_SistemaInmutable(t *testing.T) {
	repo := &mockRoleRepo{
		obtenerPorIDFunc: func(_ context.Context, _ string) (*domain.Role, error) {
			return rolSistema, nil
		},
	}
	auth := &mockAuthSvc{verificarFunc: defaultMockAuth}
	uc := NewActualizarRolCasoDeUso(repo, auth)

	cmd := &ComandoActualizarRol{
		RoleID:      "role-sys",
		Nombre:      "SUPERADMIN",
		Descripcion: "Super administrador",
		EjecutorID:  "ejecutor-1",
	}

	resp, err := uc.Ejecutar(context.Background(), cmd)
	if resp != nil {
		t.Error("expected nil response on system role modification")
	}
	if !errors.Is(err, ErrRolSistemaInmutable) {
		t.Errorf("err = %v, want %v", err, ErrRolSistemaInmutable)
	}
}

func TestActualizarRol_RolNoEncontrado(t *testing.T) {
	repo := &mockRoleRepo{
		obtenerPorIDFunc: func(_ context.Context, _ string) (*domain.Role, error) {
			return nil, domain.ErrRolNoEncontrado
		},
	}
	auth := &mockAuthSvc{verificarFunc: defaultMockAuth}
	uc := NewActualizarRolCasoDeUso(repo, auth)

	cmd := &ComandoActualizarRol{
		RoleID:      "role-inexistente",
		Nombre:      "supervisor",
		Descripcion: "Supervisor",
		EjecutorID:  "ejecutor-1",
	}

	resp, err := uc.Ejecutar(context.Background(), cmd)
	if resp != nil {
		t.Error("expected nil response on role not found")
	}
	if !errors.Is(err, domain.ErrRolNoEncontrado) {
		t.Errorf("err = %v, want %v", err, domain.ErrRolNoEncontrado)
	}
}

func TestActualizarRol_ValidacionFallida(t *testing.T) {
	repo := &mockRoleRepo{
		obtenerPorIDFunc: func(_ context.Context, _ string) (*domain.Role, error) {
			t.Error("ObtenerPorID should not be called on validation error")
			return nil, nil
		},
	}
	auth := &mockAuthSvc{verificarFunc: defaultMockAuth}
	uc := NewActualizarRolCasoDeUso(repo, auth)

	tests := []struct {
		name   string
		modify func(*ComandoActualizarRol)
		want   error
	}{
		{
			name:   "RoleID vacío",
			modify: func(cmd *ComandoActualizarRol) { cmd.RoleID = "" },
			want:   ErrRolRequerido,
		},
		{
			name:   "nombre vacío",
			modify: func(cmd *ComandoActualizarRol) { cmd.Nombre = "" },
			want:   ErrNombreRequerido,
		},
		{
			name:   "ejecutorID vacío",
			modify: func(cmd *ComandoActualizarRol) { cmd.EjecutorID = "" },
			want:   ErrEjecutorRequerido,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &ComandoActualizarRol{
				RoleID:      "role-1",
				Nombre:      "supervisor",
				Descripcion: "Supervisor",
				EjecutorID:  "ejecutor-1",
			}
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

func TestActualizarRol_ErrorRepositorio(t *testing.T) {
	repo := &mockRoleRepo{
		obtenerPorIDFunc: func(_ context.Context, _ string) (*domain.Role, error) {
			return rolNoSistema, nil
		},
		actualizarFunc: func(_ context.Context, _ *domain.Role) (*domain.Role, error) {
			return nil, domain.ErrRolDuplicado
		},
	}
	auth := &mockAuthSvc{verificarFunc: defaultMockAuth}
	uc := NewActualizarRolCasoDeUso(repo, auth)

	cmd := &ComandoActualizarRol{
		RoleID:      "role-1",
		Nombre:      "admin", // nombre ya existente
		Descripcion: "Administrador",
		EjecutorID:  "ejecutor-1",
	}

	resp, err := uc.Ejecutar(context.Background(), cmd)
	if resp != nil {
		t.Error("expected nil response on repository error")
	}
	if !errors.Is(err, domain.ErrRolDuplicado) {
		t.Errorf("err = %v, want %v", err, domain.ErrRolDuplicado)
	}
}

// ============================================================
// Tests para EliminarRol
// ============================================================

func TestEliminarRol_Exito(t *testing.T) {
	repo := &mockRoleRepo{
		obtenerPorIDFunc: func(_ context.Context, _ string) (*domain.Role, error) {
			return rolNoSistema, nil
		},
		eliminarFunc: func(_ context.Context, _ string) error {
			return nil
		},
	}
	auth := &mockAuthSvc{verificarFunc: defaultMockAuth}
	uc := NewEliminarRolCasoDeUso(repo, auth)

	cmd := &ComandoEliminarRol{
		RoleID:     "role-1",
		EjecutorID: "ejecutor-1",
	}

	resp, err := uc.Ejecutar(context.Background(), cmd)
	if err != nil {
		t.Fatalf("Ejecutar returned unexpected error: %v", err)
	}
	if resp == nil {
		t.Fatal("Ejecutar returned nil response")
	}
	if resp.ID != "role-1" {
		t.Errorf("resp.ID = %q, want %q", resp.ID, "role-1")
	}
	if !resp.Exitoso {
		t.Error("resp.Exitoso should be true")
	}
}

func TestEliminarRol_SinPermiso(t *testing.T) {
	repo := &mockRoleRepo{
		obtenerPorIDFunc: func(_ context.Context, _ string) (*domain.Role, error) {
			t.Error("ObtenerPorID should not be called when not authorized")
			return nil, nil
		},
	}
	auth := &mockAuthSvc{
		verificarFunc: func(_ context.Context, _, _ string) (bool, error) {
			return false, nil
		},
	}
	uc := NewEliminarRolCasoDeUso(repo, auth)

	cmd := &ComandoEliminarRol{
		RoleID:     "role-1",
		EjecutorID: "ejecutor-1",
	}

	resp, err := uc.Ejecutar(context.Background(), cmd)
	if resp != nil {
		t.Error("expected nil response on access denied")
	}
	if !errors.Is(err, domain.ErrAccesoDenegado) {
		t.Errorf("err = %v, want %v", err, domain.ErrAccesoDenegado)
	}
}

func TestEliminarRol_SistemaInmutable(t *testing.T) {
	repo := &mockRoleRepo{
		obtenerPorIDFunc: func(_ context.Context, _ string) (*domain.Role, error) {
			return rolSistema, nil
		},
	}
	auth := &mockAuthSvc{verificarFunc: defaultMockAuth}
	uc := NewEliminarRolCasoDeUso(repo, auth)

	cmd := &ComandoEliminarRol{
		RoleID:     "role-sys",
		EjecutorID: "ejecutor-1",
	}

	resp, err := uc.Ejecutar(context.Background(), cmd)
	if resp != nil {
		t.Error("expected nil response on system role deletion")
	}
	if !errors.Is(err, ErrRolSistemaInmutable) {
		t.Errorf("err = %v, want %v", err, ErrRolSistemaInmutable)
	}
}

func TestEliminarRol_RolNoEncontrado(t *testing.T) {
	repo := &mockRoleRepo{
		obtenerPorIDFunc: func(_ context.Context, _ string) (*domain.Role, error) {
			return nil, domain.ErrRolNoEncontrado
		},
	}
	auth := &mockAuthSvc{verificarFunc: defaultMockAuth}
	uc := NewEliminarRolCasoDeUso(repo, auth)

	cmd := &ComandoEliminarRol{
		RoleID:     "role-inexistente",
		EjecutorID: "ejecutor-1",
	}

	resp, err := uc.Ejecutar(context.Background(), cmd)
	if resp != nil {
		t.Error("expected nil response on role not found")
	}
	if !errors.Is(err, domain.ErrRolNoEncontrado) {
		t.Errorf("err = %v, want %v", err, domain.ErrRolNoEncontrado)
	}
}

func TestEliminarRol_ErrorRepositorio(t *testing.T) {
	repo := &mockRoleRepo{
		obtenerPorIDFunc: func(_ context.Context, _ string) (*domain.Role, error) {
			return rolNoSistema, nil
		},
		eliminarFunc: func(_ context.Context, _ string) error {
			return domain.ErrRolConUsuarios
		},
	}
	auth := &mockAuthSvc{verificarFunc: defaultMockAuth}
	uc := NewEliminarRolCasoDeUso(repo, auth)

	cmd := &ComandoEliminarRol{
		RoleID:     "role-1",
		EjecutorID: "ejecutor-1",
	}

	resp, err := uc.Ejecutar(context.Background(), cmd)
	if resp != nil {
		t.Error("expected nil response on repository error")
	}
	if !errors.Is(err, domain.ErrRolConUsuarios) {
		t.Errorf("err = %v, want %v", err, domain.ErrRolConUsuarios)
	}
}

func TestEliminarRol_ValidacionFallida(t *testing.T) {
	repo := &mockRoleRepo{
		obtenerPorIDFunc: func(_ context.Context, _ string) (*domain.Role, error) {
			t.Error("ObtenerPorID should not be called on validation error")
			return nil, nil
		},
	}
	auth := &mockAuthSvc{verificarFunc: defaultMockAuth}
	uc := NewEliminarRolCasoDeUso(repo, auth)

	tests := []struct {
		name   string
		modify func(*ComandoEliminarRol)
		want   error
	}{
		{
			name:   "RoleID vacío",
			modify: func(cmd *ComandoEliminarRol) { cmd.RoleID = "" },
			want:   ErrRolRequerido,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &ComandoEliminarRol{
				RoleID:     "role-1",
				EjecutorID: "ejecutor-1",
			}
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
