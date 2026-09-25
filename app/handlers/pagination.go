package handlers

type PaginatedResponse[T any] struct {
	Data   []T `json:"data"`
	Total  int `json:"total"`
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}
type PaginatedRequest struct {
	Offset int `query:"offset"`
	Limit  int `query:"limit"`
}
