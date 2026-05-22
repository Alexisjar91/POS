package domain_test

import (
	"errors"
	"testing"

	"github.com/Alexisjar91/POS/internal/users/domain"
)

func TestPermission_NuevoPermiso_Exito(t *testing.T) {
	perm, err := domain.NuevoPermiso("create_user", "Crear usuarios", "users")
	if err != nil {
		t.Fatalf("NuevoPermiso returned unexpected error: %v", err)
	}
	if perm == nil {
		t.Fatal("NuevoPermiso returned nil")
	}
	if perm.Code() != "create_user" {
		t.Errorf("Code() = %q, want %q", perm.Code(), "create_user")
	}
	if perm.Module() != "users" {
		t.Errorf("Module() = %q, want %q", perm.Module(), "users")
	}
}

func TestPermission_NuevoPermiso_CodeVacio(t *testing.T) {
	perm, err := domain.NuevoPermiso("", "desc", "users")
	if !errors.Is(err, domain.ErrCodeRequerido) {
		t.Errorf("err = %v, want %v", err, domain.ErrCodeRequerido)
	}
	if perm != nil {
		t.Errorf("perm = %v, want nil", perm)
	}
}

func TestPermission_NuevoPermiso_ModuleVacio(t *testing.T) {
	perm, err := domain.NuevoPermiso("create_user", "desc", "")
	if !errors.Is(err, domain.ErrModuleRequerido) {
		t.Errorf("err = %v, want %v", err, domain.ErrModuleRequerido)
	}
	if perm != nil {
		t.Errorf("perm = %v, want nil", perm)
	}
}

func TestPermission_NuevoPermisoDesdeBD_Exito(t *testing.T) {
	perm := domain.NuevoPermisoDesdeBD("01J", "create_user", "Crear usuarios", "users")
	if perm == nil {
		t.Fatal("NuevoPermisoDesdeBD returned nil")
	}
	if perm.Code() != "create_user" {
		t.Errorf("Code() = %q, want %q", perm.Code(), "create_user")
	}
	if perm.Module() != "users" {
		t.Errorf("Module() = %q, want %q", perm.Module(), "users")
	}
}

func TestPermission_BelongsToModule(t *testing.T) {
	perm, err := domain.NuevoPermiso("create_user", "Crear usuarios", "users")
	if err != nil {
		t.Fatalf("NuevoPermiso returned unexpected error: %v", err)
	}

	tests := []struct {
		name       string
		moduleName string
		want       bool
	}{
		{name: "mismo modulo", moduleName: "users", want: true},
		{name: "distinto modulo", moduleName: "sales", want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := perm.BelongsToModule(tt.moduleName)
			if got != tt.want {
				t.Errorf("BelongsToModule(%q) = %v, want %v", tt.moduleName, got, tt.want)
			}
		})
	}
}
