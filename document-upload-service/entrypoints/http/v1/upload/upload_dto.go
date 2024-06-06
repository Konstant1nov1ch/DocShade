package upload

// DtoIn Input data (может остаться пустым, если не требуется дополнительных данных в теле запроса)
type DtoIn struct {
}

// DtoOut Output data
type DtoOut struct {
	SessionID  string `json:"session_id"`
	DocumentID string `json:"document_id"`
	Message    string `json:"message"`
}
