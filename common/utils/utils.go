package utils

const (
	maskText = "***"
)

// MaskText маскирует строку заменяя середину на "***"
func MaskText(s string) string {
	if len(s) == 0 {
		return ""
	}
	if len(s) == 1 {
		s = s + s[:1]
	}

	return s[:2] + maskText + s[len(s)-2:]
}
