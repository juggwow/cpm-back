package response

type ResponseID struct {
	ID uint `json:"id"`
}

type ResponseError struct {
	Error string `json:"error"`
}

type ResponseData[T any] struct {
	Data  []T   `json:"data"`
	Page  int   `json:"page"`
	Limit int   `json:"limit"`
	Total int64 `json:"total"`
}
