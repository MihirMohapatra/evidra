package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRolePermissions(t *testing.T) {
	adminPerms := RoleAdmin.Permissions()
	assert.Contains(t, adminPerms, PermissionCreateOrganization)
	assert.Contains(t, adminPerms, PermissionDeleteUser)
	assert.Contains(t, adminPerms, PermissionReviewEvidence)
	assert.Len(t, adminPerms, 14)

	reviewerPerms := RoleReviewer.Permissions()
	assert.Contains(t, reviewerPerms, PermissionReadOrganization)
	assert.Contains(t, reviewerPerms, PermissionViewEvidence)
	assert.Contains(t, reviewerPerms, PermissionReviewEvidence)
	assert.NotContains(t, reviewerPerms, PermissionCreateUser)
	assert.NotContains(t, reviewerPerms, PermissionCreateOrganization)
	assert.Len(t, reviewerPerms, 4)

	var unknown Role = "unknown"
	assert.Empty(t, unknown.Permissions())
}

func TestHasPermission(t *testing.T) {
	assert.True(t, RoleAdmin.HasPermission(PermissionCreateOrganization))
	assert.True(t, RoleAdmin.HasPermission(PermissionReviewEvidence))
	assert.False(t, RoleAdmin.HasPermission("nonexistent"))

	assert.True(t, RoleReviewer.HasPermission(PermissionViewEvidence))
	assert.False(t, RoleReviewer.HasPermission(PermissionCreateUser))
}

func TestValidRoles(t *testing.T) {
	roles := ValidRoles()
	assert.ElementsMatch(t, []Role{RoleAdmin, RoleReviewer}, roles)
}

func TestParseRole(t *testing.T) {
	r, ok := ParseRole("admin")
	assert.True(t, ok)
	assert.Equal(t, RoleAdmin, r)

	r, ok = ParseRole("reviewer")
	assert.True(t, ok)
	assert.Equal(t, RoleReviewer, r)

	_, ok = ParseRole("unknown")
	assert.False(t, ok)
}
