package envfile

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/iamnator/envy/internal/model"
)

func CreatEnvFile(dir string, secrets []model.Secret) error {
	//open file to write to
	//create file if it doesn't exist

	fileName := dir
	//if .env file is not present in dir name then add it
	if !strings.HasSuffix(fileName, ".env") {
		fileName += "/.env"
	}

	//create directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(fileName), 0755); err != nil {
		return err
	}

	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	var str string
	println("writing secrets to file")
	for _, secret := range secrets {
		str += fmt.Sprintf("%s=%s # %s \n", secret.Key, secret.Value, secret.Description)
		writer.WriteString(str)
	}
	writer.Flush()
	println(".env file created successfully")

	return nil
}
