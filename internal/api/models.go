package api

import "fmt"

type Request struct {
	ID         int64  `json:"request_id"`
	FormatType string `json:"format_type"`
	Format     string `json:"format"`
}

func (r *Request) validate() error {
	if r == nil {
		return fmt.Errorf("the request cannot be equal to nil")
	}

	if r.ID == 0 || r.FormatType == "" || r.Format == "" {
		return fmt.Errorf("the request_id, format_type, foramt must be in the request")
	}

	return nil
}

type Response struct {
	Message string `json:"message,omitempty"`
	Format  string `json:"format,omitempty"`
	Error   string `json:"error,omitempty"`
}
