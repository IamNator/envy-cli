package internal

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"github.com/iamnator/envy/internal/model"
	"github.com/iamnator/envy/internal/pkg/errors"
)

func encrytSecrets(secrets []model.Secret) error {

	for i, secret := range secrets {
		encrypted, err := encrypt(secret.Value)
		if err != nil {
			return err
		}
		secrets[i].Value = encrypted
	}

	return nil
}

// setSecretOnHost uploads the secrets to the server
func setSecretOnHost(secrets []model.Secret) error {

	//make a post request
	buf, err := json.Marshal(uploadRequest{
		Secrets: secrets,
	})
	if err != nil {
		return err
	}

	resp, err := http.Post(
		HOST+"/set",
		"json/application",
		bytes.NewReader(buf))
	if err != nil {
		log.Fatal(err)
	}

	if err := errors.CheckError(resp); err != nil {
		return err
	}

	return nil

}
