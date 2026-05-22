package domain_test

import (
	"errors"
	"testing"

	"github.com/Alexisjar91/POS/internal/users/domain"
)

func TestRolePermission_NuevoRolePermission_Exito(t *testing.T) {
	rp, err := domain.NuevoRolePermission("01J", "02K")
	if err != nil {
		t.Fatalf("NuevoRolePermission returned unexpected error: %v", err)
	}
	if rp == nil {
		t.Fatal("NuevoRolePermission returned nil")
	}
	if rp.RoleID() != "01J" {
		t.Errorf("RoleID() = %q, want %q", rp.RoleID(), "01J")
	}
	if rp.PermissionID() != "02K" {
		t.Errorf("PermissionID() = %q, want %q", rp.PermissionID(), "02K")
	}
}

func TestRolePermission_NuevoRolePermission_RoleIDVacio(t *testing.T) {
	rp, err := domain.NuevoRolePermission("", "02K")
	if !errors.Is(err, domain.ErrRolNoEncontrado) {
		t.Errorf("err = %v, want %v", err, domain.ErrRolNoEncontrado)
	}
	if rp != nil {
		t.Errorf("rp = %v, want nil", rp)
	}
}

func TestRolePermission_NuevoRolePermission_PermissionIDVacio(t *testing.T) {
	rp, err := domain.NuevoRolePermission("01J", "")
	if !errors.Is(err, domain.ErrPermisoNoEncontrado) {
		t.Errorf("err = %v, want %v", err, domain.ErrPermisoNoEncontrado)
	}
	if rp != nil {
		t.Errorf("rp = %v, want nil", rp)
	}
}

func TestRolePermission_NuevoRolePermissionDesdeBD_Exito(t *testing.T) {
	rp := domain.NuevoRolePermissionDesdeBD("01J", "02K")
	if rp == nil {
		t.Fatal("NuevoRolePermissionDesdeBD returned nil")
	}
	if rp.RoleID() != "01J" {
		t.Errorf("RoleID() = %q, want %q", rp.RoleID(), "01J")
	}
	if rp.PermissionID() != "02K" {
		t.Errorf("PermissionID() = %q, want %q", rp.PermissionID(), "02K")
	}
}
