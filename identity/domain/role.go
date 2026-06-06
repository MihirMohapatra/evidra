package domain

type Role string

const (
	RoleAdmin    Role = "admin"
	RoleReviewer Role = "reviewer"
)

type Permission string

const (
	PermissionCreateOrganization Permission = "organization:create"
	PermissionReadOrganization   Permission = "organization:read"
	PermissionUpdateOrganization Permission = "organization:update"
	PermissionDeleteOrganization Permission = "organization:delete"

	PermissionCreateUser Permission = "user:create"
	PermissionReadUser   Permission = "user:read"
	PermissionUpdateUser Permission = "user:update"
	PermissionDeleteUser Permission = "user:delete"

	PermissionSubmitEvidence Permission = "evidence:submit"
	PermissionReviewEvidence Permission = "evidence:review"
	PermissionViewEvidence   Permission = "evidence:view"

	PermissionCreateAPIKey Permission = "apikey:create"
	PermissionReadAPIKey   Permission = "apikey:read"
	PermissionDeleteAPIKey Permission = "apikey:delete"
)

var rolePermissions = map[Role][]Permission{
	RoleAdmin: {
		PermissionCreateOrganization,
		PermissionReadOrganization,
		PermissionUpdateOrganization,
		PermissionDeleteOrganization,
		PermissionCreateUser,
		PermissionReadUser,
		PermissionUpdateUser,
		PermissionDeleteUser,
		PermissionSubmitEvidence,
		PermissionReviewEvidence,
		PermissionViewEvidence,
		PermissionCreateAPIKey,
		PermissionReadAPIKey,
		PermissionDeleteAPIKey,
	},
	RoleReviewer: {
		PermissionReadOrganization,
		PermissionReadUser,
		PermissionViewEvidence,
		PermissionReviewEvidence,
	},
}

func (r Role) Permissions() []Permission {
	return rolePermissions[r]
}

func (r Role) HasPermission(p Permission) bool {
	for _, perm := range r.Permissions() {
		if perm == p {
			return true
		}
	}
	return false
}

func ValidRoles() []Role {
	return []Role{RoleAdmin, RoleReviewer}
}

func ParseRole(s string) (Role, bool) {
	for _, r := range ValidRoles() {
		if string(r) == s {
			return r, true
		}
	}
	return "", false
}
