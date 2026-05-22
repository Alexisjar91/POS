package changepassword

import (
	"github.com/Alexisjar91/POS/pkg/especificacion"
	"github.com/Alexisjar91/POS/pkg/paginacion"
	"context"
	"errors"
	"testing"

	"github.com/Alexisjar91/POS/internal/users/domain"
)

// --- Mocks ---

// mockUserRepo implementa domain.UserRepository para tests de change_password.
type mockUserRepo struct {
	obtenerPorIDFunc func(ctx context.Context, id string) (*domain.User, error)
	actualizarFunc   func(ctx context.Context, user *domain.User) (*domain.User, error)
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

func (m *mockUserRepo) Actualizar(ctx context.Context, user *domain.User) (*domain.User, error) {
	return m.actualizarFunc(ctx, user)
}

func (m *mockUserRepo) Listar(_ context.Context, _ especificacion.Especificacion, _ paginacion.Paginacion) ([]*domain.User, error) {
	panic("no debería llamarse")
}

func (m *mockUserRepo) ExistePorEmail(_ context.Context, _ string) (bool, error) {
	panic("no debería llamarse")
}

// mockPasswordHasher implementa domain.PasswordHasher para tests.
type mockPasswordHasher struct {
	compareFunc func(plain, hash string) error
	hashFunc    func(plain string) (string, error)
}

func (m *mockPasswordHasher) Hash(plain string) (string, error) {
	return m.hashFunc(plain)
}

func (m *mockPasswordHasher) Compare(plain, hash string) error {
	return m.compareFunc(plain, hash)
}

// --- Helpers ---

var usuarioVal = domain.NuevoUsuarioDesdeBD("user-1", "user@example.com", "Juan Pérez", "old_hash", true, nil, "2026-05-21T10:00:00Z")

func comandoValido() *ComandoCambiarContrasena {
	return &ComandoCambiarContrasena{
		UserID:          "user-1",
		CurrentPassword: "old_password",
		NewPassword:     "new_password",
	}
}

func defaultCompare(_ string, _ string) error {
	return nil
}

func defaultHash(_ string) (string, error) {
	return "new_hash", nil
}

// --- Tests ---

func TestCambiarContrasena_Exito(t *testing.T) {
	userRepo := &mockUserRepo{
		obtenerPorIDFunc: func(_ context.Context, _ string) (*domain.User, error) {
			return usuarioVal, nil
		},
		actualizarFunc: func(_ context.Context, user *domain.User) (*domain.User, error) {
			return user, nil
		},
	}
	hasher := &mockPasswordHasher{
		compareFunc: defaultCompare,
		hashFunc:    defaultHash,
	}
	uc := NewCambiarContrasenaCasoDeUso(userRepo, hasher)

	resp, err := uc.Ejecutar(context.Background(), comandoValido())
	if err != nil {
		t.Fatalf("Ejecutar returned unexpected error: %v", err)
	}
	if resp == nil {
		t.Fatal("Ejecutar returned nil response")
	}
	if resp.ID != "user-1" {
		t.Errorf("resp.ID = %q, want %q", resp.ID, "user-1")
	}
	if !resp.Active {
		t.Error("resp.Active should be true")
	}
}

func TestCambiarContrasena_ContrasenaActualIncorrecta(t *testing.T) {
	userRepo := &mockUserRepo{
		obtenerPorIDFunc: func(_ context.Context, _ string) (*domain.User, error) {
			return usuarioVal, nil
		},
	}
	hasher := &mockPasswordHasher{
		compareFunc: func(_, _ string) error {
			return errors.New("contraseña incorrecta")
		},
	}
	uc := NewCambiarContrasenaCasoDeUso(userRepo, hasher)

	resp, err := uc.Ejecutar(context.Background(), comandoValido())
	if resp != nil {
		t.Error("expected nil response on wrong password")
	}
	if !errors.Is(err, ErrContrasenaActualIncorrecta) {
		t.Errorf("err = %v, want %v", err, ErrContrasenaActualIncorrecta)
	}
}

func TestCambiarContrasena_ContrasenaIgual(t *testing.T) {
	userRepo := &mockUserRepo{
		obtenerPorIDFunc: func(_ context.Context, _ string) (*domain.User, error) {
			return usuarioVal, nil
		},
	}
	hasher := &mockPasswordHasher{}
	uc := NewCambiarContrasenaCasoDeUso(userRepo, hasher)

	cmd := &ComandoCambiarContrasena{
		UserID:          "user-1",
		CurrentPassword: "misma_contra",
		NewPassword:     "misma_contra",
	}

	resp, err := uc.Ejecutar(context.Background(), cmd)
	if resp != nil {
		t.Error("expected nil response on same password")
	}
	if !errors.Is(err, ErrContrasenaIgual) {
		t.Errorf("err = %v, want %v", err, ErrContrasenaIgual)
	}
}

func TestCambiarContrasena_ValidacionFallida(t *testing.T) {
	userRepo := &mockUserRepo{
		obtenerPorIDFunc: func(_ context.Context, _ string) (*domain.User, error) {
			return usuarioVal, nil
		},
	}
	hasher := &mockPasswordHasher{}
	uc := NewCambiarContrasenaCasoDeUso(userRepo, hasher)

	tests := []struct {
		name   string
		modify func(*ComandoCambiarContrasena)
		want   error
	}{
		{
			name:   "userID vacío",
			modify: func(cmd *ComandoCambiarContrasena) { cmd.UserID = "" },
			want:   ErrUsuarioRequerido,
		},
		{
			name:   "currentPassword vacío",
			modify: func(cmd *ComandoCambiarContrasena) { cmd.CurrentPassword = "" },
			want:   ErrContrasenaActualRequerida,
		},
		{
			name:   "newPassword vacío",
			modify: func(cmd *ComandoCambiarContrasena) { cmd.NewPassword = "" },
			want:   ErrNuevaContrasenaRequerida,
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

func TestCambiarContrasena_UsuarioNoEncontrado(t *testing.T) {
	userRepo := &mockUserRepo{
		obtenerPorIDFunc: func(_ context.Context, _ string) (*domain.User, error) {
			return nil, domain.ErrUsuarioNoEncontrado
		},
	}
	hasher := &mockPasswordHasher{}
	uc := NewCambiarContrasenaCasoDeUso(userRepo, hasher)

	resp, err := uc.Ejecutar(context.Background(), comandoValido())
	if resp != nil {
		t.Error("expected nil response when user not found")
	}
	if !errors.Is(err, domain.ErrUsuarioNoEncontrado) {
		t.Errorf("err = %v, want %v", err, domain.ErrUsuarioNoEncontrado)
	}
}

func TestCambiarContrasena_ErrorActualizar(t *testing.T) {
	userRepo := &mockUserRepo{
		obtenerPorIDFunc: func(_ context.Context, _ string) (*domain.User, error) {
			return usuarioVal, nil
		},
		actualizarFunc: func(_ context.Context, _ *domain.User) (*domain.User, error) {
			return nil, domain.ErrRepositorio
		},
	}
	hasher := &mockPasswordHasher{
		compareFunc: defaultCompare,
		hashFunc:    defaultHash,
	}
	uc := NewCambiarContrasenaCasoDeUso(userRepo, hasher)

	resp, err := uc.Ejecutar(context.Background(), comandoValido())
	if resp != nil {
		t.Error("expected nil response on update error")
	}
	if !errors.Is(err, domain.ErrRepositorio) {
		t.Errorf("err = %v, want %v", err, domain.ErrRepositorio)
	}
}

func TestCambiarContrasena_ErrorHash(t *testing.T) {
	userRepo := &mockUserRepo{
		obtenerPorIDFunc: func(_ context.Context, _ string) (*domain.User, error) {
			return usuarioVal, nil
		},
	}
	hasher := &mockPasswordHasher{
		compareFunc: defaultCompare,
		hashFunc: func(_ string) (string, error) {
			return "", errors.New("error de hashing")
		},
	}
	uc := NewCambiarContrasenaCasoDeUso(userRepo, hasher)

	resp, err := uc.Ejecutar(context.Background(), comandoValido())
	if resp != nil {
		t.Error("expected nil response on hash error")
	}
	if err == nil || err.Error() != "error de hashing" {
		t.Errorf("err = %v, want %v", err, "error de hashing")
	}
}

func TestCambiarContrasena_ErrorCompare(t *testing.T) {
	userRepo := &mockUserRepo{
		obtenerPorIDFunc: func(_ context.Context, _ string) (*domain.User, error) {
			return usuarioVal, nil
		},
	}
	hasher := &mockPasswordHasher{
		compareFunc: func(_, _ string) error {
			return errors.New("error tecnico de bcrypt")
		},
	}
	uc := NewCambiarContrasenaCasoDeUso(userRepo, hasher)

	resp, err := uc.Ejecutar(context.Background(), comandoValido())
	if resp != nil {
		t.Error("expected nil response on compare error")
	}
	if !errors.Is(err, ErrContrasenaActualIncorrecta) {
		t.Errorf("err = %v, want %v", err, ErrContrasenaActualIncorrecta)
	}
}

func TestCambiarContrasena_NoLlamaActualizarSiCompareFalla(t *testing.T) {
	actualizarCalled := false
	userRepo := &mockUserRepo{
		obtenerPorIDFunc: func(_ context.Context, _ string) (*domain.User, error) {
			return usuarioVal, nil
		},
		actualizarFunc: func(_ context.Context, _ *domain.User) (*domain.User, error) {
			actualizarCalled = true
			return nil, nil
		},
	}
	hasher := &mockPasswordHasher{
		compareFunc: func(_, _ string) error {
			return errors.New("contraseña incorrecta")
		},
	}
	uc := NewCambiarContrasenaCasoDeUso(userRepo, hasher)

	resp, err := uc.Ejecutar(context.Background(), comandoValido())
	if resp != nil {
		t.Error("expected nil response on wrong password")
	}
	if !errors.Is(err, ErrContrasenaActualIncorrecta) {
		t.Errorf("err = %v, want %v", err, ErrContrasenaActualIncorrecta)
	}
	if actualizarCalled {
		t.Error("Actualizar should not be called when compare fails")
	}
}
