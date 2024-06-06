package http

// ErrorHttp описание ошибки
type ErrorHttp struct {
	ErrorText string   `json:"errorText"`
	Details   []string `json:"details"`
}
