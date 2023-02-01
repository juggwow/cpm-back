package response

type ID struct {
	ID uint `json:"id"`
}

type Error struct {
	Error string `json:"error"`
}

// type Delete struct {
// 	Massage string `json:"del" default:"success"`
// }

type Data[T any] struct {
	Data  []T   `json:"data"`
	Page  int   `json:"page"`
	Limit int   `json:"limit"`
	Total int64 `json:"total"`
}

type DataDoc[I any, D any] struct {
	Item  I     `json:"item"`
	DOC   D     `json:"doc"`
	Page  int   `json:"page"`
	Limit int   `json:"limit"`
	Total int64 `json:"total"`
}
