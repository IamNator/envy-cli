package errors

import (
	"fmt"
	"io"
	"net/http"
)

func CheckError(resp *http.Response) error {
	if resp.StatusCode > 201 {
		msg, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf("unexpected status code: %d: %s", resp.StatusCode, string(msg))
	}
	return nil
}
