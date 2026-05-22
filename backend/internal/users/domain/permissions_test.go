package domain_test

import (
	"testing"

	"github.com/Alexisjar91/POS/internal/users/domain"
)

func TestPermissionConstants_UserModule(t *testing.T) {
	tests := []struct {
		name     string
		constant string
		want     string
	}{
		{name: "CreateUser", constant: domain.CreateUser, want: "create_user"},
		{name: "DisableUser", constant: domain.DisableUser, want: "disable_user"},
		{name: "EnableUser", constant: domain.EnableUser, want: "enable_user"},
		{name: "AssignRole", constant: domain.AssignRole, want: "assign_role"},
		{name: "ManageRoles", constant: domain.ManageRoles, want: "manage_roles"},
		{name: "ViewUsers", constant: domain.ViewUsers, want: "view_users"},
		{name: "ResetUserPassword", constant: domain.ResetUserPassword, want: "reset_user_password"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.constant != tt.want {
				t.Errorf("domain.%s = %q, want %q", tt.name, tt.constant, tt.want)
			}
		})
	}
}


