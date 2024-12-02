// utils/pagination.go

package utils

type PaginatedResponse struct {
	Total int64       `json:"total"`
	Data  interface{} `json:"data"`
}
