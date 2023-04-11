package envfile

import (
	"bufio"
	"os"
	"strings"

	"github.com/iamnator/envy/internal/model"
)

func ReadFromFile(filename string) ([]model.Secret, error) {
	// read file
	fileReader, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer fileReader.Close()

	var secrets []model.Secret

	// Create a scanner to read the file line by line
	scanner := bufio.NewScanner(fileReader)
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "#") && strings.Contains(line, "=") {
			parts := strings.SplitN(line, "=", 2)
			if len(parts) != 2 {
				continue
			}
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			secrets = append(secrets, model.Secret{Key: key, Value: value})
		}
	}

	// Check for any scanner errors
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return secrets, nil
}
