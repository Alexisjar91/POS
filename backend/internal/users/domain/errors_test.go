package domain_test

import (
	"errors"
	"testing"

	"github.com/Alexisjar91/POS/internal/users/domain"
)

func TestAutorizacionErrors_Existen(t *testing.T) {
	tests := []struct {
		name string
		err  error
	}{
		{name: "ErrAccesoDenegado", err: domain.ErrAccesoDenegado},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err == nil {
				t.Fatalf("%s should not be nil", tt.name)
			}
			if tt.err.Error() == "" {
				t.Fatalf("%s message should not be empty", tt.name)
			}
		})
	}
}

func TestRepositorioErrors_Existen(t *testing.T) {
	tests := []struct {
		name string
		err  error
	}{
		{name: "ErrEmailDuplicado", err: domain.ErrEmailDuplicado},
		{name: "ErrUsuarioNoEncontrado", err: domain.ErrUsuarioNoEncontrado},
		{name: "ErrRolNoEncontrado", err: domain.ErrRolNoEncontrado},
		{name: "ErrPermisoNoEncontrado", err: domain.ErrPermisoNoEncontrado},
		{name: "ErrRolDuplicado", err: domain.ErrRolDuplicado},
		{name: "ErrPermisoDuplicado", err: domain.ErrPermisoDuplicado},
		{name: "ErrRolSistemaInmutable", err: domain.ErrRolSistemaInmutable},
		{name: "ErrRolConUsuarios", err: domain.ErrRolConUsuarios},
		{name: "ErrPermisoEnUso", err: domain.ErrPermisoEnUso},
		{name: "ErrRepositorio", err: domain.ErrRepositorio},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err == nil {
				t.Fatalf("%s should not be nil", tt.name)
			}
			if tt.err.Error() == "" {
				t.Fatalf("%s message should not be empty", tt.name)
			}
		})
	}
}

func TestPermissionErrors_Existen(t *testing.T) {
	tests := []struct {
		name string
		err  error
	}{
		{name: "ErrCodeRequerido", err: domain.ErrCodeRequerido},
		{name: "ErrModuleRequerido", err: domain.ErrModuleRequerido},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err == nil {
				t.Fatalf("%s should not be nil", tt.name)
			}
			if tt.err.Error() == "" {
				t.Fatalf("%s message should not be empty", tt.name)
			}
		})
	}
}

func TestRoleErrors_Existen(t *testing.T) {
	tests := []struct {
		name string
		err  error
	}{
		{name: "ErrNombreRolRequerido", err: domain.ErrNombreRolRequerido},
		{name: "ErrPermisoYaAsignado", err: domain.ErrPermisoYaAsignado},
		{name: "ErrPermisoNoAsignado", err: domain.ErrPermisoNoAsignado},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err == nil {
				t.Fatalf("%s should not be nil", tt.name)
			}
			if tt.err.Error() == "" {
				t.Fatalf("%s message should not be empty", tt.name)
			}
		})
	}
}

func TestAllErrors_MensajesUnicos(t *testing.T) {
	allErrors := []struct {
		name string
		err  error
	}{
		{name: "ErrAccesoDenegado", err: domain.ErrAccesoDenegado},
		{name: "ErrEmailDuplicado", err: domain.ErrEmailDuplicado},
		{name: "ErrUsuarioNoEncontrado", err: domain.ErrUsuarioNoEncontrado},
		{name: "ErrRolNoEncontrado", err: domain.ErrRolNoEncontrado},
		{name: "ErrPermisoNoEncontrado", err: domain.ErrPermisoNoEncontrado},
		{name: "ErrRolDuplicado", err: domain.ErrRolDuplicado},
		{name: "ErrPermisoDuplicado", err: domain.ErrPermisoDuplicado},
		{name: "ErrRolSistemaInmutable", err: domain.ErrRolSistemaInmutable},
		{name: "ErrRolConUsuarios", err: domain.ErrRolConUsuarios},
		{name: "ErrPermisoEnUso", err: domain.ErrPermisoEnUso},
		{name: "ErrRepositorio", err: domain.ErrRepositorio},
		{name: "ErrCodeRequerido", err: domain.ErrCodeRequerido},
		{name: "ErrModuleRequerido", err: domain.ErrModuleRequerido},
		{name: "ErrNombreRolRequerido", err: domain.ErrNombreRolRequerido},
		{name: "ErrPermisoYaAsignado", err: domain.ErrPermisoYaAsignado},
		{name: "ErrPermisoNoAsignado", err: domain.ErrPermisoNoAsignado},
		{name: "ErrEmailRequerido", err: domain.ErrEmailRequerido},
		{name: "ErrNombreRequerido", err: domain.ErrNombreRequerido},
		{name: "ErrPasswordHashRequerido", err: domain.ErrPasswordHashRequerido},
		{name: "ErrCreatedByRequerido", err: domain.ErrCreatedByRequerido},
		{name: "ErrUsuarioYaActivo", err: domain.ErrUsuarioYaActivo},
		{name: "ErrUsuarioYaInactivo", err: domain.ErrUsuarioYaInactivo},
	}
	seen := make(map[string]string, len(allErrors))
	for _, e := range allErrors {
		msg := e.err.Error()
		if prev, ok := seen[msg]; ok {
			t.Errorf("mensaje duplicado entre %q y %q: %q", prev, e.name, msg)
		}
		seen[msg] = e.name
	}
}

func TestErrorsIs_Comparable(t *testing.T) {
	t.Run("ErrEmailRequerido es igual a sí mismo", func(t *testing.T) {
		if !errors.Is(domain.ErrEmailRequerido, domain.ErrEmailRequerido) {
			t.Error("errors.Is(ErrEmailRequerido, ErrEmailRequerido) debe ser true")
		}
	})

	t.Run("ErrEmailRequerido no es igual a ErrNombreRequerido", func(t *testing.T) {
		if errors.Is(domain.ErrEmailRequerido, domain.ErrNombreRequerido) {
			t.Error("errors.Is(ErrEmailRequerido, ErrNombreRequerido) debe ser false")
		}
	})

	t.Run("ErrAccesoDenegado es igual a sí mismo", func(t *testing.T) {
		if !errors.Is(domain.ErrAccesoDenegado, domain.ErrAccesoDenegado) {
			t.Error("errors.Is(ErrAccesoDenegado, ErrAccesoDenegado) debe ser true")
		}
	})

	t.Run("ErrRepositorio no es igual a ErrEmailDuplicado", func(t *testing.T) {
		if errors.Is(domain.ErrRepositorio, domain.ErrEmailDuplicado) {
			t.Error("errors.Is(ErrRepositorio, ErrEmailDuplicado) debe ser false")
		}
	})
}

func TestUserErrors_Existen(t *testing.T) {
	tests := []struct {
		name string
		err  error
	}{
		{name: "ErrEmailRequerido", err: domain.ErrEmailRequerido},
		{name: "ErrNombreRequerido", err: domain.ErrNombreRequerido},
		{name: "ErrPasswordHashRequerido", err: domain.ErrPasswordHashRequerido},
		{name: "ErrCreatedByRequerido", err: domain.ErrCreatedByRequerido},
		{name: "ErrUsuarioYaActivo", err: domain.ErrUsuarioYaActivo},
		{name: "ErrUsuarioYaInactivo", err: domain.ErrUsuarioYaInactivo},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err == nil {
				t.Fatalf("%s should not be nil", tt.name)
			}
			if tt.err.Error() == "" {
				t.Fatalf("%s message should not be empty", tt.name)
			}
		})
	}
}
