package domain_test

import (
	"errors"
	"testing"

	"github.com/Alexisjar91/POS/internal/users/domain"
)

func TestRole_NuevoRol_Exito(t *testing.T) {
	role, err := domain.NuevoRol("cashier", "Cajero")
	if err != nil {
		t.Fatalf("NuevoRol returned unexpected error: %v", err)
	}
	if role == nil {
		t.Fatal("NuevoRol returned nil")
	}
	if role.Name() != "cashier" {
		t.Errorf("Name() = %q, want %q", role.Name(), "cashier")
	}
	if role.IsSystem() != false {
		t.Errorf("IsSystem() = %v, want false", role.IsSystem())
	}
	if role.IsOwner() != false {
		t.Errorf("IsOwner() = %v, want false", role.IsOwner())
	}
	if role.IsAdmin() != false {
		t.Errorf("IsAdmin() = %v, want false", role.IsAdmin())
	}
}

func TestRole_NuevoRol_NombreVacio(t *testing.T) {
	role, err := domain.NuevoRol("", "desc")
	if !errors.Is(err, domain.ErrNombreRolRequerido) {
		t.Errorf("err = %v, want %v", err, domain.ErrNombreRolRequerido)
	}
	if role != nil {
		t.Errorf("role = %v, want nil", role)
	}
}

func TestRole_NuevoRolSistema_Exito(t *testing.T) {
	role, err := domain.NuevoRolSistema("OWNER", "Propietario del sistema")
	if err != nil {
		t.Fatalf("NuevoRolSistema returned unexpected error: %v", err)
	}
	if role == nil {
		t.Fatal("NuevoRolSistema returned nil")
	}
	if role.Name() != "OWNER" {
		t.Errorf("Name() = %q, want %q", role.Name(), "OWNER")
	}
	if role.IsSystem() != true {
		t.Errorf("IsSystem() = %v, want true", role.IsSystem())
	}
	if role.IsOwner() != true {
		t.Errorf("IsOwner() = %v, want true", role.IsOwner())
	}
	if role.IsAdmin() != false {
		t.Errorf("IsAdmin() = %v, want false", role.IsAdmin())
	}
}

func TestRole_NuevoRolSistema_NombreVacio(t *testing.T) {
	role, err := domain.NuevoRolSistema("", "desc")
	if !errors.Is(err, domain.ErrNombreRolRequerido) {
		t.Errorf("err = %v, want %v", err, domain.ErrNombreRolRequerido)
	}
	if role != nil {
		t.Errorf("role = %v, want nil", role)
	}
}

func TestRole_NuevoRolDesdeBD_Exito(t *testing.T) {
	role := domain.NuevoRolDesdeBD("01J", "admin", "Administrador", true)
	if role == nil {
		t.Fatal("NuevoRolDesdeBD returned nil")
	}
	if role.Name() != "admin" {
		t.Errorf("Name() = %q, want %q", role.Name(), "admin")
	}
	if role.IsSystem() != true {
		t.Errorf("IsSystem() = %v, want true", role.IsSystem())
	}
}
