package utils

type ReaponseError struct {
	Error string `json:"error"`
}

type ReaponseID struct {
	ID uint `json:"id"`
}

type Reaponse struct {
	Error string `json:"error,omitempty"`
	Msg   string `json:"msg,omitempty"`
	ID    uint   `json:"id,omitempty"`
}
