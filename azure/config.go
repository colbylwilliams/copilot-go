// Package azure provides helpful types and functions for use with Azure OpenAI.
package azure

import "os"

const (
	AzureOpenAIAPIVersionDefault string = "2024-07-01-preview" //"2024-06-01"
)

const (
	AzureTenantIDKey         string = "AZURE_TENANT_ID"
	AzureOpenAIAPIKey        string = "AZURE_OPENAI_API_KEY"
	AzureOpenAIEndpointKey   string = "AZURE_OPENAI_ENDPOINT"
	AzureOpenAIAPIVersionKey string = "OPENAI_API_VERSION"
)

// Config represents the configuration needed for Azure OpenAI.
type Config struct {
	TenantID         string
	OpenAIEndpoint   string
	OpenAIAPIKey     string
	OpenAIAPIVersion string
}

// LoadConfig reads the environment variables and returns a new AzureConfig.
// This function assumes that the environment variables (for example, from
// and .env file) are already loaded.
func LoadConfig() *Config {
	tenantID := getEnvOrDefault(AzureTenantIDKey, "")
	endpoint := getEnvOrDefault(AzureOpenAIEndpointKey, "")

	if tenantID == "" || endpoint == "" {
		return nil
	}

	// for now we won't require the api key
	apiKey := getEnvOrDefault(AzureOpenAIAPIKey, "")

	return &Config{
		TenantID:         tenantID,
		OpenAIEndpoint:   endpoint,
		OpenAIAPIKey:     apiKey,
		OpenAIAPIVersion: getEnvOrDefault(AzureOpenAIAPIVersionKey, AzureOpenAIAPIVersionDefault),
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value != "" {
		return value
	}
	return defaultValue
}
