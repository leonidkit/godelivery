package xml2json

import (
	"fmt"
	"strings"

	xtoj "github.com/basgys/goxml2json"
)

type XML2JSON struct{}

func New() *XML2JSON {
	return &XML2JSON{}
}

func (xj *XML2JSON) Convert(xml string, formatType string) ([]byte, error) {
	rd := strings.NewReader(xml)

	res, err := xtoj.Convert(rd)
	if err != nil {
		return []byte{}, fmt.Errorf("xml converting to json error: %v", err)
	}

	return res.Bytes(), err
}
