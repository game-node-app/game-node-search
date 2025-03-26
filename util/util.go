package util

import (
	Manticoresearch "github.com/manticoresoftware/manticoresearch-go"
	"os"
)

func GetEnv(env, fallback string) string {
	v, ok := os.LookupEnv(env)
	if !ok {
		return fallback
	}
	return v

}

func Contains(slice []string, item string) bool {
	for _, a := range slice {
		if a == item {
			return true
		}
	}
	return false
}

func GetManticoreInstance() *Manticoresearch.APIClient {
	// Create an instance of API client
	configuration := Manticoresearch.NewConfiguration()
	url := GetEnv("MANTICORE_URL", "http://localhost:9308")

	configuration.Servers[0].URL = url
	apiClient := Manticoresearch.NewAPIClient(configuration)

	return apiClient
}
