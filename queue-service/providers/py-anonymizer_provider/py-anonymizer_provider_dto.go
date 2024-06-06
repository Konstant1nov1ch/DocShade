package anonymizer_provider

// AnonymizeDocumentRequest структура для запроса
type AnonymizeDocumentRequest struct {
	Document []byte `json:"document"`
}

// AnonymizeDocumentResponse структура для ответа
type AnonymizeDocumentResponse struct {
	AnonymizedDocument []byte `json:"anonymized_document"`
}
