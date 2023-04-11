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

	if resp.ContentLength < 1 { //no body
		return fmt.Errorf("no body in response, check network connection")
	}

	return nil
}
