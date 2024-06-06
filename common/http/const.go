package http

type Methods = int

// Константы типа методов
const (
	GetMethod Methods = iota
	PutMethod
	PostMethod
	DeleteMethod
	PatchMethod
)

// MaskHeaders карта заголовков
var MaskHeaders = map[string]struct{}{AccessToken: {}}

const (
	AccessToken = "AccessToken"
)

// SwaggerRoute путь до сваггера
const (
	SwaggerRoute = "/swagger/*"
)
