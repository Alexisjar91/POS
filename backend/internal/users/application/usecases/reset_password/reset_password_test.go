package resetpassword

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

// mockPasswordHasher implementa domain.PasswordHasher para tests.
type mockPasswordHasher struct {
	hashFunc    func(plain string) (string, error)
	compareFunc func(plain, hash string) error
}

func (m *mockPasswordHasher) Hash(plain string) (string, error) {
	if m.hashFunc != nil {
		return m.hashFunc(plain)
	}
	return "", nil
}

func (m *mockPasswordHasher) Compare(plain, hash string) error {
	if m.compareFunc != nil {
		return m.compareFunc(plain, hash)
	}
	return nil
}

// --- Helpers ---

func strPtr(s string) *string {
	return &s
}

func comandoValido() *ComandoResetearContrasena {
	return &ComandoResetearContrasena{
		UserID:      "user-target",
		NewPassword: "new_password",
		EjecutorID:  "user-executor",
	}
}

func usuarioActivo() *domain.User {
	return domain.NuevoUsuarioDesdeBD(
		"user-target",
		"target@test.com",
		"Target User",
		"old_hash",
		true,
		strPtr("creator"),
		"2024-01-01",
	)
}

func defaultMockAuth(_ context.Context, _ string, _ string) (bool, error) {
	return true, nil
}

func defaultMockHash(_ string) (string, error) {
	return "new_hash", nil
}

func defaultMockActualizar(_ context.Context, user *domain.User) (*domain.User, error) {
	return user, nil
}

// --- Tests ---

func TestResetearContrasena_Exito(t *testing.T) {
	repo := &mockUserRepo{
		obtenerPorIDFunc: func(_ context.Context, _ string) (*domain.User, error) {
			return usuarioActivo(), nil
		},
		actualizarFunc: defaultMockActualizar,
	}
	auth := &mockAuthSvc{
		verificarFunc: defaultMockAuth,
	}
	hasher := &mockPasswordHasher{
		hashFunc: defaultMockHash,
	}
	uc := NewResetearContrasenaCasoDeUso(repo, hasher, auth)

	resp, err := uc.Ejecutar(context.Background(), comandoValido())

	if err != nil {
		t.Fatalf("Ejecutar returned unexpected error: %v", err)
	}
	if resp == nil {
		t.Fatal("Ejecutar returned nil response")
	}
	if resp.ID != "user-target" {
		t.Errorf("resp.ID = %q, want %q", resp.ID, "user-target")
	}
	if !resp.Active {
		t.Error("resp.Active should be true")
	}
}

func TestResetearContrasena_SinPermiso(t *testing.T) {
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
	hasher := &mockPasswordHasher{}
	uc := NewResetearContrasenaCasoDeUso(repo, hasher, auth)

	resp, err := uc.Ejecutar(context.Background(), comandoValido())

	if resp != nil {
		t.Error("expected nil response on access denied")
	}
	if !errors.Is(err, domain.ErrAccesoDenegado) {
		t.Errorf("err = %v, want %v", err, domain.ErrAccesoDenegado)
	}
}

func TestResetearContrasena_ValidacionFallida(t *testing.T) {
	repo := &mockUserRepo{
		obtenerPorIDFunc: func(_ context.Context, _ string) (*domain.User, error) {
			t.Error("ObtenerPorID should not be called on validation failure")
			return nil, nil
		},
	}
	auth := &mockAuthSvc{
		verificarFunc: defaultMockAuth,
	}
	hasher := &mockPasswordHasher{}
	uc := NewResetearContrasenaCasoDeUso(repo, hasher, auth)

	tests := []struct {
		name string
		cmd  *ComandoResetearContrasena
		want error
	}{
		{
			name: "UserID vacío",
			cmd:  &ComandoResetearContrasena{UserID: "", NewPassword: "new_pass", EjecutorID: "executor"},
			want: ErrUsuarioRequerido,
		},
		{
			name: "NewPassword vacío",
			cmd:  &ComandoResetearContrasena{UserID: "target", NewPassword: "", EjecutorID: "executor"},
			want: ErrNuevaContrasenaRequerida,
		},
		{
			name: "EjecutorID vacío",
			cmd:  &ComandoResetearContrasena{UserID: "target", NewPassword: "new_pass", EjecutorID: ""},
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

func TestResetearContrasena_UsuarioNoEncontrado(t *testing.T) {
	repo := &mockUserRepo{
		obtenerPorIDFunc: func(_ context.Context, _ string) (*domain.User, error) {
			return nil, domain.ErrUsuarioNoEncontrado
		},
	}
	auth := &mockAuthSvc{
		verificarFunc: defaultMockAuth,
	}
	hasher := &mockPasswordHasher{
		hashFunc: defaultMockHash,
	}
	uc := NewResetearContrasenaCasoDeUso(repo, hasher, auth)

	resp, err := uc.Ejecutar(context.Background(), comandoValido())

	if resp != nil {
		t.Error("expected nil response on user not found")
	}
	if !errors.Is(err, domain.ErrUsuarioNoEncontrado) {
		t.Errorf("err = %v, want %v", err, domain.ErrUsuarioNoEncontrado)
	}
}

func TestResetearContrasena_ErrorHash(t *testing.T) {
	repo := &mockUserRepo{
		obtenerPorIDFunc: func(_ context.Context, _ string) (*domain.User, error) {
			return usuarioActivo(), nil
		},
	}
	auth := &mockAuthSvc{
		verificarFunc: defaultMockAuth,
	}
	hasher := &mockPasswordHasher{
		hashFunc: func(_ string) (string, error) {
			return "", errors.New("error de hashing")
		},
	}
	uc := NewResetearContrasenaCasoDeUso(repo, hasher, auth)

	resp, err := uc.Ejecutar(context.Background(), comandoValido())

	if resp != nil {
		t.Error("expected nil response on hash error")
	}
	if err == nil || err.Error() != "error de hashing" {
		t.Errorf("err = %v, want %v", err, "error de hashing")
	}
}

func TestResetearContrasena_ErrorActualizar(t *testing.T) {
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
	}
	hasher := &mockPasswordHasher{
		hashFunc: defaultMockHash,
	}
	uc := NewResetearContrasenaCasoDeUso(repo, hasher, auth)

	resp, err := uc.Ejecutar(context.Background(), comandoValido())

	if resp != nil {
		t.Error("expected nil response on update error")
	}
	if !errors.Is(err, domain.ErrRepositorio) {
		t.Errorf("err = %v, want %v", err, domain.ErrRepositorio)
	}
}

func TestResetearContrasena_NoLlamaRepoSiNoAutorizado(t *testing.T) {
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
	hasher := &mockPasswordHasher{}
	uc := NewResetearContrasenaCasoDeUso(repo, hasher, auth)

	resp, _ := uc.Ejecutar(context.Background(), comandoValido())

	if resp != nil {
		t.Error("expected nil response on access denied")
	}
	if obtenerPorIDLlamado {
		t.Error("ObtenerPorID should not be called when authorization fails")
	}
}

func TestResetearContrasena_NoLlamaActualizarSiHashFalla(t *testing.T) {
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
	}
	hasher := &mockPasswordHasher{
		hashFunc: func(_ string) (string, error) {
			return "", errors.New("error de hashing")
		},
	}
	uc := NewResetearContrasenaCasoDeUso(repo, hasher, auth)

	resp, err := uc.Ejecutar(context.Background(), comandoValido())

	if resp != nil {
		t.Error("expected nil response on hash error")
	}
	if err == nil || err.Error() != "error de hashing" {
		t.Errorf("err = %v, want %v", err, "error de hashing")
	}
	if actualizarLlamado {
		t.Error("Actualizar should not be called when hash fails")
	}
}
