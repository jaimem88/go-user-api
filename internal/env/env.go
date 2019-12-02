package env

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

// ServiceAddr - the address used by the grpc server, eg. ":8899"
func ServiceAddr() string {
	return GetFromEnvOrDefault("SERVICE_ADDR", ":8899")
}

// Version reads the version number from the version.json file
// panics if no file exists
func Version() string {
	s := struct {
		Version string `json:"version"`
	}{}
	f, err := ioutil.ReadFile("version.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(f, &s)
	if err != nil {
		panic(err)
	}
	return s.Version
}

// GetFromEnvOrDefault - returns the environment variable if it exists else returns to default
func GetFromEnvOrDefault(name, defaultValue string) string {
	fromEnv := os.Getenv(name)
	if fromEnv == "" {
		return defaultValue
	}
	return fromEnv
}
