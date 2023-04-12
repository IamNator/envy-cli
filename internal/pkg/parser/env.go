package parser

import (
	"fmt"
	"strings"

	"github.com/iamnator/envy-cli/internal/model"
)

// parseEnv parses a string of the form "key=value" and returns a Secret
func ParseEnv(arg string) (model.Secret, error) {

	parts := strings.SplitN(arg, "=", 2)
	if len(parts) != 2 {
		return model.Secret{}, fmt.Errorf("invalid argument: %s", arg)
	}

	//clean up the key
	key := strings.TrimSpace(parts[0])
	//clean up the value
	value := strings.TrimSpace(parts[1])

	//remove quotes
	if strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"") {
		value = value[1 : len(value)-1]
	}

	return model.Secret{
		Key:   key,
		Value: value,
	}, nil
}
