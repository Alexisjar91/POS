package domain_test

import (
	"errors"
	"testing"

	"github.com/Alexisjar91/POS/internal/users/domain"
)

func TestUser_NuevoUsuario_Exito(t *testing.T) {
	createdBy := "01J-creator"
	user, err := domain.NuevoUsuario("test@test.com", "Test User", "$2a$10$hash", createdBy)
	if err != nil {
		t.Fatalf("NuevoUsuario returned unexpected error: %v", err)
	}
	if user == nil {
		t.Fatal("NuevoUsuario returned nil")
	}
	if user.Email() != "test@test.com" {
		t.Errorf("Email() = %q, want %q", user.Email(), "test@test.com")
	}
	if user.FullName() != "Test User" {
		t.Errorf("FullName() = %q, want %q", user.FullName(), "Test User")
	}
	if user.PasswordHash() != "$2a$10$hash" {
		t.Errorf("PasswordHash() = %q, want %q", user.PasswordHash(), "$2a$10$hash")
	}
	if !user.IsActive() {
		t.Error("IsActive() should be true")
	}
	if user.CreatedBy() == nil {
		t.Fatal("CreatedBy() should not be nil")
	}
	if *user.CreatedBy() != createdBy {
		t.Errorf("CreatedBy() = %q, want %q", *user.CreatedBy(), createdBy)
	}
	if user.CreatedAt() != "" {
		t.Errorf("CreatedAt() = %q, want empty string (infrastructure sets it)", user.CreatedAt())
	}
}

func TestUser_NuevoUsuario_EmailVacio(t *testing.T) {
	user, err := domain.NuevoUsuario("", "Test User", "hash", "creator")
	if !errors.Is(err, domain.ErrEmailRequerido) {
		t.Errorf("err = %v, want %v", err, domain.ErrEmailRequerido)
	}
	if user != nil {
		t.Errorf("user = %v, want nil", user)
	}
}

func TestUser_NuevoUsuario_NombreVacio(t *testing.T) {
	user, err := domain.NuevoUsuario("test@test.com", "", "hash", "creator")
	if !errors.Is(err, domain.ErrNombreRequerido) {
		t.Errorf("err = %v, want %v", err, domain.ErrNombreRequerido)
	}
	if user != nil {
		t.Errorf("user = %v, want nil", user)
	}
}

func TestUser_NuevoUsuario_PasswordHashVacio(t *testing.T) {
	user, err := domain.NuevoUsuario("test@test.com", "Test User", "", "creator")
	if !errors.Is(err, domain.ErrPasswordHashRequerido) {
		t.Errorf("err = %v, want %v", err, domain.ErrPasswordHashRequerido)
	}
	if user != nil {
		t.Errorf("user = %v, want nil", user)
	}
}

func TestUser_NuevoUsuario_CreatedByVacio(t *testing.T) {
	user, err := domain.NuevoUsuario("test@test.com", "Test User", "hash", "")
	if !errors.Is(err, domain.ErrCreatedByRequerido) {
		t.Errorf("err = %v, want %v", err, domain.ErrCreatedByRequerido)
	}
	if user != nil {
		t.Errorf("user = %v, want nil", user)
	}
}

func TestUser_NuevoUsuarioOwner_Exito(t *testing.T) {
	user, err := domain.NuevoUsuarioOwner("owner@test.com", "Owner User", "$2a$10$hash")
	if err != nil {
		t.Fatalf("NuevoUsuarioOwner returned unexpected error: %v", err)
	}
	if user == nil {
		t.Fatal("NuevoUsuarioOwner returned nil")
	}
	if user.Email() != "owner@test.com" {
		t.Errorf("Email() = %q, want %q", user.Email(), "owner@test.com")
	}
	if user.FullName() != "Owner User" {
		t.Errorf("FullName() = %q, want %q", user.FullName(), "Owner User")
	}
	if user.PasswordHash() != "$2a$10$hash" {
		t.Errorf("PasswordHash() = %q, want %q", user.PasswordHash(), "$2a$10$hash")
	}
	if !user.IsActive() {
		t.Error("IsActive() should be true")
	}
	if user.CreatedBy() != nil {
		t.Error("CreatedBy() should be nil for OWNER")
	}
}

func TestUser_NuevoUsuarioOwner_EmailVacio(t *testing.T) {
	user, err := domain.NuevoUsuarioOwner("", "Owner", "hash")
	if !errors.Is(err, domain.ErrEmailRequerido) {
		t.Errorf("err = %v, want %v", err, domain.ErrEmailRequerido)
	}
	if user != nil {
		t.Errorf("user = %v, want nil", user)
	}
}

func TestUser_NuevoUsuarioDesdeBD_Exito(t *testing.T) {
	createdBy := "creator-id"
	user := domain.NuevoUsuarioDesdeBD("01J-user", "test@test.com", "Test User", "hash123", true, &createdBy, "2024-01-01T00:00:00Z")
	if user == nil {
		t.Fatal("NuevoUsuarioDesdeBD returned nil")
	}
	if user.ID() != "01J-user" {
		t.Errorf("ID() = %q, want %q", user.ID(), "01J-user")
	}
	if user.Email() != "test@test.com" {
		t.Errorf("Email() = %q, want %q", user.Email(), "test@test.com")
	}
	if user.FullName() != "Test User" {
		t.Errorf("FullName() = %q, want %q", user.FullName(), "Test User")
	}
	if user.PasswordHash() != "hash123" {
		t.Errorf("PasswordHash() = %q, want %q", user.PasswordHash(), "hash123")
	}
	if !user.IsActive() {
		t.Error("IsActive() should be true")
	}
	if user.CreatedBy() == nil || *user.CreatedBy() != "creator-id" {
		t.Errorf("CreatedBy() = %v, want %v", user.CreatedBy(), &createdBy)
	}
	if user.CreatedAt() != "2024-01-01T00:00:00Z" {
		t.Errorf("CreatedAt() = %q, want %q", user.CreatedAt(), "2024-01-01T00:00:00Z")
	}
}

func TestUser_Disable(t *testing.T) {
	t.Run("desactivar usuario activo", func(t *testing.T) {
		user, err := domain.NuevoUsuario("test@test.com", "Test", "hash", "creator")
		if err != nil {
			t.Fatal(err)
		}
		if !user.IsActive() {
			t.Fatal("precondition: user should be active")
		}
		err = user.Disable()
		if err != nil {
			t.Errorf("Disable() returned unexpected error: %v", err)
		}
		if user.IsActive() {
			t.Error("user should be inactive after Disable()")
		}
	})

	t.Run("desactivar usuario ya inactivo", func(t *testing.T) {
		createdBy := "creator"
		user := domain.NuevoUsuarioDesdeBD("01J", "test@test.com", "Test", "hash", false, &createdBy, "")
		if user.IsActive() {
			t.Fatal("precondition: user should be inactive")
		}
		err := user.Disable()
		if !errors.Is(err, domain.ErrUsuarioYaInactivo) {
			t.Errorf("err = %v, want %v", err, domain.ErrUsuarioYaInactivo)
		}
	})
}

func TestUser_Enable(t *testing.T) {
	t.Run("activar usuario inactivo", func(t *testing.T) {
		createdBy := "creator"
		user := domain.NuevoUsuarioDesdeBD("01J", "test@test.com", "Test", "hash", false, &createdBy, "")
		if user.IsActive() {
			t.Fatal("precondition: user should be inactive")
		}
		err := user.Enable()
		if err != nil {
			t.Errorf("Enable() returned unexpected error: %v", err)
		}
		if !user.IsActive() {
			t.Error("user should be active after Enable()")
		}
	})

	t.Run("activar usuario ya activo", func(t *testing.T) {
		user, err := domain.NuevoUsuario("test@test.com", "Test", "hash", "creator")
		if err != nil {
			t.Fatal(err)
		}
		if !user.IsActive() {
			t.Fatal("precondition: user should be active")
		}
		err = user.Enable()
		if !errors.Is(err, domain.ErrUsuarioYaActivo) {
			t.Errorf("err = %v, want %v", err, domain.ErrUsuarioYaActivo)
		}
	})
}

func TestUser_SetPasswordHash(t *testing.T) {
	user, err := domain.NuevoUsuario("test@test.com", "Test", "oldhash", "creator")
	if err != nil {
		t.Fatal(err)
	}
	if user.PasswordHash() != "oldhash" {
		t.Fatalf("precondition: PasswordHash() = %q, want %q", user.PasswordHash(), "oldhash")
	}
	user.SetPasswordHash("newhash123")
	if user.PasswordHash() != "newhash123" {
		t.Errorf("after SetPasswordHash, PasswordHash() = %q, want %q", user.PasswordHash(), "newhash123")
	}
}

func TestUser_CreatedBy(t *testing.T) {
	t.Run("usuario normal tiene creador", func(t *testing.T) {
		creator := "01J-creator"
		user, err := domain.NuevoUsuario("test@test.com", "Test", "hash", creator)
		if err != nil {
			t.Fatal(err)
		}
		if user.CreatedBy() == nil {
			t.Fatal("CreatedBy() should not be nil for normal user")
		}
		if *user.CreatedBy() != creator {
			t.Errorf("CreatedBy() = %q, want %q", *user.CreatedBy(), creator)
		}
	})

	t.Run("usuario OWNER tiene createdBy nil", func(t *testing.T) {
		owner, err := domain.NuevoUsuarioOwner("owner@test.com", "Owner", "hash")
		if err != nil {
			t.Fatal(err)
		}
		if owner.CreatedBy() != nil {
			t.Error("CreatedBy() should be nil for OWNER")
		}
	})
}
