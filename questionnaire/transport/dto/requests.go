package dto

type UploadRequest struct {
	Title string `form:"title" validate:"required,max=255"`
}

type UpdateStatusRequest struct {
	Status string `json:"status" validate:"required,oneof=uploaded queued parsing parsed failed assigned"`
}
