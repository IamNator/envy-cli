package internal

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/iamnator/envy/internal/model"
	"github.com/iamnator/envy/internal/pkg/envfile"
	"github.com/iamnator/envy/internal/pkg/parser"
	"github.com/iamnator/envy/pkg/encryption"
)

type uploadRequest struct {
	Secrets []model.Secret `json:"secrets"`
}

var (
	ENVY_SECRET = "secret" // this is the secret used to encrypt the secrets
	HOST        = "https://iamnator-super-space-carnival-64g6qw5px9rh4759-8080.preview.app.github.dev"
)

var (
	source = flag.String("source", "", "file of env variables from")
	set    = flag.String("set", "", "set a single secret")
	get    = flag.String("get", "", "get a single secret")
	// environment = flag.String("environment", "", "development, staging, production or any other environment")
	dir    = flag.String("dir", "", "directory to store .env file")
	secret = flag.String("secret", "", "secret to encrypt/decrypt secrets")
)

func check(s *string) bool {
	if s == nil {
		return false
	}
	return *s != ""
}

var decrypt = encryption.MakeDecrypter(ENVY_SECRET)
var encrypt = encryption.MakeEncrypter(ENVY_SECRET)

func Run() {
	flag.Parse()

	if check(secret) {
		ENVY_SECRET = *secret
	} else {
		ENVY_SECRET = os.Getenv("ENVY_SECRET")
		if ENVY_SECRET == "" {
			log.Println(`
			ENVY_SECRET not set in environment or secret flag
			please set the ENVY_SECRET environment variable or use the -mac flag

			e.g  export ENVY_SECRET=<secret> 
			or 
			envy --secret=<secret>`)
		}
	}

	// upload file
	if check(source) {
		println("uploading file")
		secrets, err := envfile.ReadFromFile(*source)
		if err != nil {
			log.Fatal(err)
		}

		if err := encrytSecrets(secrets); err != nil {
			log.Fatal(err)
		}

		if err := setSecretOnHost(secrets); err != nil {
			log.Fatal(err)
		}
		println("done")
		return

	}

	// set single secret

	if check(set) {
		secret, err := parser.ParseEnv(*set)
		if err != nil {
			log.Fatal(err)
		}

		secret.Value, err = encrypt(secret.Value)
		if err != nil {
			log.Fatal(err)
		}

		if err := setSecretOnHost([]model.Secret{secret}); err != nil {
			log.Fatal(err)
		}
		println("done")

		return
	}

	// download all secrets to env file
	if check(dir) {
		secrets, err := getSecretFromHost()
		if err != nil {
			log.Fatal(err)
		}

		if err := envfile.CreatEnvFile(*dir, secrets); err != nil {
			log.Fatal(err)
		}
	}

	// get single secret
	if check(get) {
		secret, err := getSecretFromHost(*get)
		if err != nil {
			log.Fatal(err)
		}

		if len(secret) < 1 {
			log.Fatal("secret not found")
		}

		for _, s := range secret {

			value, err := decrypt(s.Value)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Printf("%s=%s \n", secret[0].Key, value)
		}
	}

}
