package domain_test

import (
	"errors"
	"testing"

	"github.com/Alexisjar91/POS/internal/users/domain"
)

func TestUserRole_NuevoUserRole_Exito(t *testing.T) {
	ur, err := domain.NuevoUserRole("01J-user", "02K-role")
	if err != nil {
		t.Fatalf("NuevoUserRole returned unexpected error: %v", err)
	}
	if ur == nil {
		t.Fatal("NuevoUserRole returned nil")
	}
	if ur.UserID() != "01J-user" {
		t.Errorf("UserID() = %q, want %q", ur.UserID(), "01J-user")
	}
	if ur.RoleID() != "02K-role" {
		t.Errorf("RoleID() = %q, want %q", ur.RoleID(), "02K-role")
	}
}

func TestUserRole_NuevoUserRole_UserIDVacio(t *testing.T) {
	ur, err := domain.NuevoUserRole("", "02K-role")
	if !errors.Is(err, domain.ErrUsuarioNoEncontrado) {
		t.Errorf("err = %v, want %v", err, domain.ErrUsuarioNoEncontrado)
	}
	if ur != nil {
		t.Errorf("ur = %v, want nil", ur)
	}
}

func TestUserRole_NuevoUserRole_RoleIDVacio(t *testing.T) {
	ur, err := domain.NuevoUserRole("01J-user", "")
	if !errors.Is(err, domain.ErrRolNoEncontrado) {
		t.Errorf("err = %v, want %v", err, domain.ErrRolNoEncontrado)
	}
	if ur != nil {
		t.Errorf("ur = %v, want nil", ur)
	}
}

func TestUserRole_NuevoUserRoleDesdeBD_Exito(t *testing.T) {
	ur := domain.NuevoUserRoleDesdeBD("01J-user", "02K-role")
	if ur == nil {
		t.Fatal("NuevoUserRoleDesdeBD returned nil")
	}
	if ur.UserID() != "01J-user" {
		t.Errorf("UserID() = %q, want %q", ur.UserID(), "01J-user")
	}
	if ur.RoleID() != "02K-role" {
		t.Errorf("RoleID() = %q, want %q", ur.RoleID(), "02K-role")
	}
}
