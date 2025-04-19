package common

type PaginationParams struct {
	Limit int
}

type Pagination[T any] struct {
	Items []T  `json:"items"`
	Total int  `json:"total"`
	Limit int  `json:"limit"`
	Page  int  `json:"page"`
	Next  bool `json:"next"`
	Prev  bool `json:"prev"`
}
