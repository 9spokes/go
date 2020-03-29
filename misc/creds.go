package misc

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

// LoadCredentials reads a JSON object from a file representing API client_id/client_secret pairs and returns a map[string]string with the results
func LoadCredentials(path string) (map[string]string, error) {

	if path != "" {

		creds := make(map[string]string)

		credsFile, err := os.Open(path)
		if err != nil {
			return nil, err
		}
		defer credsFile.Close()

		byteValue, _ := ioutil.ReadAll(credsFile)
		err = json.Unmarshal(byteValue, &creds)
		if err != nil {
			return nil, err
		}
		return creds, err
	}

	return nil, nil
}
