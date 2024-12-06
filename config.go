package copilot

import (
	"errors"
	"os"
	"strconv"

	"github.com/colbylwilliams/copilot-go/azure"
	"github.com/joho/godotenv"
)

const (
	OpenAIChatModelDefault string = "gpt-4o"
)

const (
	GitHubAppIDKey             string = "GITHUB_APP_ID"
	GitHubAppClientIDKey       string = "GITHUB_APP_CLIENT_ID"
	GitHubAppClientSecretKey   string = "GITHUB_APP_CLIENT_SECRET"
	GitHubAppPrivateKeyPathKey string = "GITHUB_APP_PRIVATE_KEY_PATH"
	GitHubAppWebhookSecretKey  string = "GITHUB_APP_WEBHOOK_SECRET"
	GitHubAppFQDNKey           string = "GITHUB_APP_FQDN"
	GitHubAppUserAgentKey      string = "GITHUB_APP_USER_AGENT"
	OpenAIChatModelKey         string = "OPENAI_CHAT_MODEL"
)

// Config represents the configuration of the app.
type Config struct {
	// Environment is the environment the app is running in (development, production, etc).
	// It is resolved from the ENVIRONMENT environment variable.
	Environment string
	// HTTPPort is the port the HTTP server will listen on.
	// It is resolved from the PORT environment variable.
	HTTPPort string
	// GitHubAppFQDN is the fully qualified domain name of the GitHub App.
	// It is resolved from the GITHUB_APP_FQDN environment variable.
	// If using devoutness or ngrok, this should be the public URL.
	GitHubAppFQDN string
	// GitHubAppID is the app ID of the GitHub App.
	// It is resolved from the GITHUB_APP_ID environment variable.
	// Note that this is not required for the GitHub App to function.
	GitHubAppID int64
	// GitHubAppClientID is the client ID of the GitHub App.
	// It is resolved from the GITHUB_APP_CLIENT_ID environment variable.
	GitHubAppClientID string
	// GitHubAppClientSecret is the client secret of the GitHub App.
	// It is resolved from the GITHUB_APP_CLIENT_SECRET environment variable.
	GitHubAppClientSecret string
	// GitHubAppPrivateKeyPath is the path to the private key of the GitHub App.
	// It is resolved from the GITHUB_APP_PRIVATE_KEY_PATH environment variable.
	// The file should be a PEM file.
	GitHubAppPrivateKeyPath string
	// GitHubAppPrivateKey is the private key of the GitHub App.
	// It is resolved from the pem file at GitHubAppPrivateKeyPath.
	GitHubAppPrivateKey []byte
	// GitHubAppWebhookSecret is the secret used to validate GitHub App webhooks.
	// It is resolved from the GITHUB_APP_WEBHOOK_SECRET environment variable.
	GitHubAppWebhookSecret string
	// GitHubAppUserAgent is the user agent to use when making requests to the GitHub API.
	// It is resolved from the GITHUB_APP_USER_AGENT environment variable.
	GitHubAppUserAgent string
	// ChatModel is the OpenAI chat model to use.
	// It is resolved from the OPENAI_CHAT_MODEL environment variable.
	// If not set, it defaults to "gpt-4o".
	ChatModel string
	// Azure is the configuration for Azure OpenAI.
	Azure *azure.Config
}

// LoadConfig reads the environment variables and returns a new Config.
// env is a list of .env files to load. If none are provided,
// it will default to loading .env in the current path.
func LoadConfig(env ...string) (*Config, error) {
	// Load environment variables from .env files.
	// Load doesn't really return an error, so we ignore it.
	_ = godotenv.Load(env...)

	cfg := &Config{}

	cfg.Environment = getEnvOrDefault("ENVIRONMENT", "development")
	cfg.HTTPPort = getEnvOrDefault("PORT", "")

	// github
	if appId := os.Getenv(GitHubAppIDKey); appId != "" {
		if id, err := strconv.ParseInt(appId, 10, 64); err == nil {
			cfg.GitHubAppID = id
		}
	}

	cfg.GitHubAppClientID = getRequiredEnv(GitHubAppClientIDKey)
	cfg.GitHubAppClientSecret = getRequiredEnv(GitHubAppClientSecretKey)
	cfg.GitHubAppPrivateKeyPath = getRequiredEnv(GitHubAppPrivateKeyPathKey)

	// TODO - allow for directly setting the private key with GITHUB_APP_PRIVATE_KEY
	// Read key from pem file
	cfg.GitHubAppPrivateKey = getGitHubPrivateKey(cfg.GitHubAppPrivateKeyPath)

	cfg.GitHubAppUserAgent = getRequiredEnv(GitHubAppUserAgentKey)
	cfg.GitHubAppWebhookSecret = getRequiredEnv(GitHubAppWebhookSecretKey)
	cfg.GitHubAppFQDN = getRequiredEnv(GitHubAppFQDNKey)

	// chat
	cfg.ChatModel = getEnvOrDefault(OpenAIChatModelKey, OpenAIChatModelDefault)

	// azure
	cfg.Azure = azure.LoadConfig()

	return cfg, nil
}

// IsProduction returns true if the environment is production (or staging).
func (cfg *Config) IsProduction() bool {
	return !cfg.IsDevelopment()
}

// IsDevelopment returns true if the environment is development.
func (cfg *Config) IsDevelopment() bool {
	return cfg.Environment == "development"
}

func getRequiredEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		panic(errors.New("Missing required environment variable: " + key))
	}
	return value
}

func getEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value != "" {
		return value
	}
	return defaultValue
}

func getGitHubPrivateKey(pemFile string) []byte {
	// Read key from pem file
	if _, err := os.Stat(pemFile); err == nil {
		pemBytes, err := os.ReadFile(pemFile)
		if err != nil {
			panic(err)
		}
		return pemBytes
	}
	panic("GitHub App private key not found")
}
