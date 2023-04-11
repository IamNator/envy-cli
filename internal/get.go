package internal

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/iamnator/envy/internal/model"
	"github.com/iamnator/envy/internal/pkg/errors"
)

type downloadResponse struct {
	Secrets []model.Secret `json:"secrets"`
}

// getSecretFromHost downloads the secrets from the server and writes them to a file
func getSecretFromHost(keys ...string) ([]model.Secret, error) {
	resp, err := http.Get(HOST + "/get?key=" + strings.Join(keys, ","))
	if err != nil {
		return nil, err
	}

	if err := errors.CheckError(resp); err != nil {
		return nil, err
	}

	var body downloadResponse
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return nil, err
	}

	return body.Secrets, nil
}
