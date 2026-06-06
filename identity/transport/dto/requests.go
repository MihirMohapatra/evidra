package dto

type CreateOrganizationRequest struct {
	Name string `json:"name" validate:"required,max=255"`
	Slug string `json:"slug" validate:"required,alphanum,max=100"`
}

type UpdateOrganizationRequest struct {
	Name string `json:"name" validate:"required,max=255"`
}

type CreateUserRequest struct {
	OrganizationID string `json:"organization_id" validate:"required,uuid"`
	Email          string `json:"email" validate:"required,email"`
	Password       string `json:"password" validate:"required,min=8"`
	Role           string `json:"role" validate:"required,oneof=admin reviewer"`
}

type UpdateUserRequest struct {
	Email string `json:"email" validate:"omitempty,email"`
	Role  string `json:"role" validate:"omitempty,oneof=admin reviewer"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type CreateAPIKeyRequest struct {
	Name string `json:"name" validate:"required,max=255"`
}
