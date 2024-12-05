package copilot

import (
	"errors"
	"os"

	"github.com/colbylwilliams/copilot-go/azure"
	"github.com/joho/godotenv"
)

const (
	OpenAIChatModelDefault string = "gpt-4o"
)

const (
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
	Environment             string
	HTTPPort                string
	GitHubAppFQDN           string
	GitHubAppClientID       string
	GitHubAppClientSecret   string
	GitHubAppPrivateKeyPath string
	GitHubAppPrivateKey     []byte
	GitHubAppWebhookSecret  string
	GitHubAppUserAgent      string
	ChatModel               string
	Azure                   *azure.Config
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
