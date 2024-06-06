package anonymizer_provider

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"

	"gitlab.com/docshade/common/core"
)

type Anonymizer interface {
	InitAnonymizer() error
	AnonymizeDocument(document []byte, filename string) ([]byte, error)
}

type anonymizer struct {
	cfg core.AnonymizerConfig
}

func NewAnonymizer(cfg core.AnonymizerConfig) Anonymizer {
	return &anonymizer{
		cfg: cfg,
	}
}

func (a *anonymizer) InitAnonymizer() error {
	return nil
}

func (a *anonymizer) AnonymizeDocument(document []byte, filename string) ([]byte, error) {
	url := a.cfg.URI

	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	// Создание поля формы для файла
	fw, err := w.CreateFormFile("file", filename)
	if err != nil {
		return nil, fmt.Errorf("failed to create form file: %v", err)
	}

	// Установка правильного Content-Type для файла
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="file"; filename="%s"`, filename))
	h.Set("Content-Type", "application/pdf")

	fw, err = w.CreatePart(h)
	if err != nil {
		return nil, fmt.Errorf("failed to create form part: %v", err)
	}

	// Запись содержимого файла в поле формы
	if _, err = io.Copy(fw, bytes.NewReader(document)); err != nil {
		return nil, fmt.Errorf("failed to copy document to form file: %v", err)
	}

	// Завершение записи multipart/form-data
	w.Close()

	request, err := http.NewRequest("POST", url, &b)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	// Установка заголовков
	request.Header.Set("Content-Type", w.FormDataContentType())

	// Логирование заголовков запроса и тела запроса для отладки
	// fmt.Printf("Request headers: %v\n", request.Header)
	// fmt.Printf("Request body: %s\n", b.String())

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(response.Body)
		bodyString := string(bodyBytes)
		return nil, fmt.Errorf("received non-200 response: %d, body: %s", response.StatusCode, bodyString)
	}

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	return responseBody, nil
}
