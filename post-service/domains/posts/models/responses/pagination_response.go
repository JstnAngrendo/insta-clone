package responses

type PaginationResponse struct {
	Page       int         `json:"page"`
	Size       int         `json:"size"`
	TotalPages int         `json:"totalPages"`
	TotalItems int64       `json:"totalItems"`
	Data       interface{} `json:"data"`
}
