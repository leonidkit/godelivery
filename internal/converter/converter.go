package converter

type Interface interface {
	// Convert xml string to json string
	Convert(xml string, formatType string) ([]byte, error)
}
