package httputil

import (
	"net/http"
	"encoding/json"
	"fmt"
)

type HTTPError struct {
	Code int
	Message string
}

func (e *HTTPError) Error() string {
	return e.Message
}

func (e *HTTPError) WriteTo(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(e.Code)
	b, err := json.Marshal(e)
	if err != nil {
		return fmt.Errorf("marshal HTTPError should never fail (%v)", err)
	}
	if _, err := w.Write(b); err != nil {
		return err
	}
	return nil
}

func NewHTTPError(code int, m string) *HTTPError {
	return &HTTPError{
		Message: m,
		Code:    code,
	}
}
