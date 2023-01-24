package response

type ID struct {
	ID uint `json:"id"`
}

type Error struct {
	Error string `json:"error"`
}

type Data[T any] struct {
	Data  []T   `json:"data"`
	Page  int   `json:"page"`
	Limit int   `json:"limit"`
	Total int64 `json:"total"`
}
