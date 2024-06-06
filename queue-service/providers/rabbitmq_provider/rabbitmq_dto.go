package rabbitmq_provider

type DocumentMessage struct {
	DocumentID       string `json:"document_id"`
	OriginalFileName string `json:"original_file_name"`
	S3Path           string `json:"s3_path"`
	SessionID        string `json:"session_id"`
}
