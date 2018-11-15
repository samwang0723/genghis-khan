package facebook

import (
	"encoding/json"
	"fmt"
	"io"
)

// QueryResponse is the response sent back by Facebook when setting up things
// like greetings or call-to-actions
type QueryResponse struct {
	Error  *QueryError `json:"error,omitempty"`
	Result string      `json:"result,omitempty"`
}

// QueryError is representing an error sent back by Facebook
type QueryError struct {
	Message      string `json:"message"`
	Type         string `json:"type"`
	Code         int    `json:"code"`
	ErrorSubCode int    `json:"error_subcode"`
	FBTraceID    string `json:"fbtrace_id"`
}

// CheckFacebookError parse error response
func CheckFacebookError(r io.Reader) error {
	var err error
	qr := QueryResponse{}
	err = json.NewDecoder(r).Decode(&qr)
	if qr.Error != nil {
		err = fmt.Errorf("Facebook error: %s", qr.Error.Message)
		return err
	}
	return nil
}
