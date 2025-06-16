package dto

// PaginationRequest represents common pagination query parameters
type PaginationRequest struct {
	Page   int    `json:"page" form:"page" query:"page" validate:"omitempty,min=1" example:"1"`
	Limit  int    `json:"limit" form:"limit" query:"limit" validate:"omitempty,min=1,max=100" example:"10"`
	Search string `json:"search" form:"search" query:"search" validate:"omitempty,max=255" example:"john"`
	Order  string `json:"order" form:"order" query:"order" validate:"omitempty,max=50" example:"created_at"`
	Sort   string `json:"sort" form:"sort" query:"sort" validate:"omitempty,oneof=asc desc" example:"asc"`
}
