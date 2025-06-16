package dto

type PaginationMeta struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	TotalRows  int64 `json:"totalRows"`
	TotalPages int   `json:"totalPages"`
}

type PaginationResponseDTO[T any] struct {
	Data []T            `json:"data"`
	Meta PaginationMeta `json:"meta"`
}