package createuser

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
	crearFunc func(ctx context.Context, user *domain.User) (*domain.User, error)
}

func (m *mockUserRepo) Crear(ctx context.Context, user *domain.User) (*domain.User, error) {
	if m.crearFunc != nil {
		return m.crearFunc(ctx, user)
	}
	return nil, nil
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

func (m *mockUserRepo) Listar(_ context.Context, _ especificacion.Especificacion, _ paginacion.Paginacion) ([]*domain.User, error) {
	panic("no debería llamarse")
}

func (m *mockUserRepo) ExistePorEmail(_ context.Context, _ string) (bool, error) {
	panic("no debería llamarse")
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
	panic("no debería llamarse en tests de create_user")
}

// --- Helpers ---

func crearCmdValido() *ComandoCrearUsuario {
	return &ComandoCrearUsuario{
		Email:      "user@example.com",
		FullName:   "Juan Pérez",
		Password:   "$2a$10$rOzR0a5jN6cX1YbYbYbYbO",
		CreatedBy:  "01J-creator",
		EjecutorID: "01J-executor",
	}
}

func defaultMockCrear(_ context.Context, user *domain.User) (*domain.User, error) {
	createdBy := user.CreatedBy()
	return domain.NuevoUsuarioDesdeBD(
		"01J-new-user",
		user.Email(),
		user.FullName(),
		user.PasswordHash(),
		true,
		createdBy,
		"2026-05-21T10:00:00Z",
	), nil
}

func defaultMockAuth(_ context.Context, _ string, _ string) (bool, error) {
	return true, nil
}

// --- Tests ---

func TestCrearUsuario_Exito(t *testing.T) {
	repo := &mockUserRepo{crearFunc: defaultMockCrear}
	auth := &mockAuthSvc{verificarFunc: defaultMockAuth}
	uc := NewCrearUsuarioCasoDeUso(repo, auth)

	resp, err := uc.Ejecutar(context.Background(), crearCmdValido())

	if err != nil {
		t.Fatalf("Ejecutar returned unexpected error: %v", err)
	}
	if resp == nil {
		t.Fatal("Ejecutar returned nil response")
	}
	if resp.ID != "01J-new-user" {
		t.Errorf("resp.ID = %q, want %q", resp.ID, "01J-new-user")
	}
	if resp.Email != "user@example.com" {
		t.Errorf("resp.Email = %q, want %q", resp.Email, "user@example.com")
	}
	if resp.FullName != "Juan Pérez" {
		t.Errorf("resp.FullName = %q, want %q", resp.FullName, "Juan Pérez")
	}
	if !resp.Active {
		t.Error("resp.Active should be true")
	}
	if resp.CreatedAt != "2026-05-21T10:00:00Z" {
		t.Errorf("resp.CreatedAt = %q, want %q", resp.CreatedAt, "2026-05-21T10:00:00Z")
	}
}

func TestCrearUsuario_EmailInvalido(t *testing.T) {
	repo := &mockUserRepo{crearFunc: defaultMockCrear}
	auth := &mockAuthSvc{verificarFunc: defaultMockAuth}
	uc := NewCrearUsuarioCasoDeUso(repo, auth)

	tests := []struct {
		name  string
		email string
		want  error
	}{
		{name: "email vacío", email: "", want: ErrEmailRequerido},
		{name: "email sin arroba", email: "usuario", want: ErrEmailInvalido},
		{name: "email sin dominio", email: "usuario@", want: ErrEmailInvalido},
		{name: "email con espacios", email: " ", want: ErrEmailRequerido},
		{name: "email inválido", email: "@example.com", want: ErrEmailInvalido},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := crearCmdValido()
			cmd.Email = tt.email

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

func TestCrearUsuario_CamposRequeridos(t *testing.T) {
	repo := &mockUserRepo{crearFunc: defaultMockCrear}
	auth := &mockAuthSvc{verificarFunc: defaultMockAuth}
	uc := NewCrearUsuarioCasoDeUso(repo, auth)

	tests := []struct {
		name    string
		modify  func(*ComandoCrearUsuario)
		want    error
	}{
		{
			name:    "nombre vacío",
			modify:  func(cmd *ComandoCrearUsuario) { cmd.FullName = "" },
			want:    ErrNombreRequerido,
		},
		{
			name:    "nombre solo espacios",
			modify:  func(cmd *ComandoCrearUsuario) { cmd.FullName = "   " },
			want:    ErrNombreRequerido,
		},
		{
			name:    "password vacío",
			modify:  func(cmd *ComandoCrearUsuario) { cmd.Password = "" },
			want:    ErrPasswordRequerido,
		},
		{
			name:    "createdBy vacío",
			modify:  func(cmd *ComandoCrearUsuario) { cmd.CreatedBy = "" },
			want:    ErrCreadorRequerido,
		},
		{
			name:    "ejecutorID vacío",
			modify:  func(cmd *ComandoCrearUsuario) { cmd.EjecutorID = "" },
			want:    ErrEjecutorRequerido,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := crearCmdValido()
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

func TestCrearUsuario_SinPermiso(t *testing.T) {
	repo := &mockUserRepo{crearFunc: defaultMockCrear}
	auth := &mockAuthSvc{
		verificarFunc: func(_ context.Context, _, _ string) (bool, error) {
			return false, nil
		},
	}
	uc := NewCrearUsuarioCasoDeUso(repo, auth)

	resp, err := uc.Ejecutar(context.Background(), crearCmdValido())

	if resp != nil {
		t.Error("expected nil response on access denied")
	}
	if !errors.Is(err, domain.ErrAccesoDenegado) {
		t.Errorf("err = %v, want %v", err, domain.ErrAccesoDenegado)
	}
}

func TestCrearUsuario_ErrorAutorizacion(t *testing.T) {
	repo := &mockUserRepo{crearFunc: defaultMockCrear}
	auth := &mockAuthSvc{
		verificarFunc: func(_ context.Context, _, _ string) (bool, error) {
			return false, errors.New("error de infraestructura")
		},
	}
	uc := NewCrearUsuarioCasoDeUso(repo, auth)

	resp, err := uc.Ejecutar(context.Background(), crearCmdValido())

	if resp != nil {
		t.Error("expected nil response on auth error")
	}
	if err == nil || err.Error() != "error de infraestructura" {
		t.Errorf("err = %v, want %v", err, "error de infraestructura")
	}
}

func TestCrearUsuario_ErrorRepositorio(t *testing.T) {
	repo := &mockUserRepo{
		crearFunc: func(_ context.Context, _ *domain.User) (*domain.User, error) {
			return nil, domain.ErrEmailDuplicado
		},
	}
	auth := &mockAuthSvc{verificarFunc: defaultMockAuth}
	uc := NewCrearUsuarioCasoDeUso(repo, auth)

	resp, err := uc.Ejecutar(context.Background(), crearCmdValido())

	if resp != nil {
		t.Error("expected nil response on repository error")
	}
	if !errors.Is(err, domain.ErrEmailDuplicado) {
		t.Errorf("err = %v, want %v", err, domain.ErrEmailDuplicado)
	}
}

func TestCrearUsuario_NoLlamaRepositorioSiValidacionFalla(t *testing.T) {
	llamado := false
	repo := &mockUserRepo{
		crearFunc: func(_ context.Context, _ *domain.User) (*domain.User, error) {
			llamado = true
			return nil, nil
		},
	}
	auth := &mockAuthSvc{verificarFunc: defaultMockAuth}
	uc := NewCrearUsuarioCasoDeUso(repo, auth)

	cmd := crearCmdValido()
	cmd.Email = "invalido" // sin arroba

	resp, _ := uc.Ejecutar(context.Background(), cmd)

	if resp != nil {
		t.Error("expected nil response on validation error")
	}
	if llamado {
		t.Error("repository should not be called when validation fails")
	}
}

func TestCrearUsuario_NoLlamaRepositorioSiNoAutorizado(t *testing.T) {
	llamado := false
	repo := &mockUserRepo{
		crearFunc: func(_ context.Context, _ *domain.User) (*domain.User, error) {
			llamado = true
			return nil, nil
		},
	}
	auth := &mockAuthSvc{
		verificarFunc: func(_ context.Context, _, _ string) (bool, error) {
			return false, nil
		},
	}
	uc := NewCrearUsuarioCasoDeUso(repo, auth)

	resp, _ := uc.Ejecutar(context.Background(), crearCmdValido())

	if resp != nil {
		t.Error("expected nil response on access denied")
	}
	if llamado {
		t.Error("repository should not be called when authorization fails")
	}
}
