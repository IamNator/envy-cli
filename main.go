package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type Secret struct {
	Key         string `json:"key"`
	Value       string `json:"value"`
	Description string `json:"description"`
}

type uploadRequest struct {
	Secrets []Secret `json:"secrets"`
}

var (
	MAC  = "secret" //message authentication code (MAC)
	HOST = "https://iamnator-super-space-carnival-64g6qw5px9rh4759-8080.preview.app.github.dev"
)

var (
	source      = flag.String("source", "", "file of env variables from")
	set         = flag.String("set", "", "set a single secret")
	get         = flag.String("get", "", "get a single secret")
	environment = flag.String("environment", "", "development, staging, production or any other environment")
	dir         = flag.String("dir", "", "directory to store .env file")
)

func check(s *string) bool {
	if s == nil {
		return false
	}
	return *s != ""
}

func main() {
	flag.Parse()

	// upload file
	if check(source) {
		println("uploading file")
		secrets, err := readFromFile(*source)
		if err != nil {
			log.Fatal(err)
		}
		if err := setSecrets(secrets); err != nil {
			log.Fatal(err)
		}
		println("done")
		return

	}

	// set single secret

	if check(set) {
		secret, err := parseEnv(*set)
		if err != nil {
			log.Fatal(err)
		}
		if err := setSecrets([]Secret{secret}); err != nil {
			log.Fatal(err)
		}
		println("done")

		return
	}

	// get all secrets and write to .env file
	if check(dir) {
		secrets, err := getSecrets()
		if err != nil {
			log.Fatal(err)
		}

		if err := creatEnvFile(*dir, secrets); err != nil {
			log.Fatal(err)
		}
	}

	// get single secret
	if check(get) {
		secret, err := getSecrets(*get)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s=%s \n", secret[0].Key, decrypt(secret[0].Value, MAC))
	}

}

// parseEnv parses a string of the form "key=value" and returns a Secret
func parseEnv(arg string) (Secret, error) {

	parts := strings.SplitN(arg, "=", 2)
	if len(parts) != 2 {
		return Secret{}, fmt.Errorf("invalid argument: %s", arg)
	}

	//clean up the key
	key := strings.TrimSpace(parts[0])
	//clean up the value
	value := strings.TrimSpace(parts[1])

	//remove quotes
	if strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"") {
		value = value[1 : len(value)-1]
	}

	return Secret{
		Key:   key,
		Value: value,
	}, nil
}

func encrypt(b string, key string) string {
	return b
}

func decrypt(b string, key string) string {
	return b
}

func readFromFile(filename string) ([]Secret, error) {
	// read file
	fileReader, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer fileReader.Close()

	var secrets []Secret

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
			secrets = append(secrets, Secret{Key: key, Value: encrypt(value, MAC)})
		}
	}

	// Check for any scanner errors
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return secrets, nil
}

// setSecrets uploads the secrets to the server
func setSecrets(secrets []Secret) error {

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

	if err := checkError(resp); err != nil {
		return err
	}

	return nil

}

type downloadResponse struct {
	Secrets []Secret `json:"secrets"`
}

// getSecrets downloads the secrets from the server and writes them to a file
func getSecrets(keys ...string) ([]Secret, error) {
	resp, err := http.Get(HOST + "/get?key=" + strings.Join(keys, ","))
	if err != nil {
		return nil, err
	}

	if err := checkError(resp); err != nil {
		return nil, err
	}

	var body downloadResponse

	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return nil, err
	}

	return body.Secrets, nil
}

func creatEnvFile(dir string, secrets []Secret) error {
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
		str += fmt.Sprintf("%s=%s # %s \n", secret.Key, decrypt(secret.Value, MAC), secret.Description)
		writer.WriteString(str)
	}
	writer.Flush()
	println(".env file created successfully")

	return nil
}

func checkError(resp *http.Response) error {
	if resp.StatusCode > 300 {
		msg, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf("unexpected status code: %d: %s", resp.StatusCode, msg)
	}
	return nil
}

func setEnv(secrets []Secret) {

	for _, secret := range secrets {
		os.Setenv(secret.Key, decrypt(secret.Value, MAC))
	}

	println("env set")

}
