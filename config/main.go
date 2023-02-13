package config

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path"
	"runtime"

	"github.com/nemoden/chat/renderer"
)

var (
	ConfigDir      string = getConfigDir()
	CacheDir       string = getCacheDir()
	ApiKeyFilePath string = getApiKeyFilePath()
	InitPrompt     string = `Please initialise chat using the init command:

$ chat init
`
	// @TODO revisit instructions
	TokenFileDoesntExistPrompt string = fmt.Sprintf(`Oops. Looks like your token file %s doesn't exist.

Please obtain the openai API key from their website https://platform.openai.com/account/api-keys

Once you have the API key, you can either add it manually to %s:

$ echo "<your-api-key>" > %s

Or if you didn't run chat init, it's adviced that you do this instead:

$ chat init

Alternatively, you can provide a token using environment variable %s. How to set it depends on your shell, i.e. 

- for bash shell, add export %s=<your-api-key> to your ~/.bash_profile
- for zsh shell, add export %s=<your-api-key> to your ~/.zsh_profile
- for fish shell, add set -x %s <your-api-key> to your ~/.config/fish/config.fish
`, ApiKeyFilePath, ApiKeyFilePath, ApiKeyFilePath, API_KEY_ENV_VAR_NAME, API_KEY_ENV_VAR_NAME, API_KEY_ENV_VAR_NAME, API_KEY_ENV_VAR_NAME)

	TokenFileNotReadablePrompt string = fmt.Sprintf(`Oops. Looks like your token file %s is not readable.

chat stores the api token in that file.

Another option is to store the API key in the environment variable. How to set it depends on your shell, i.e. 

- for bash shell, add export %s=<your-api-key> to your ~/.bash_profile
- for zsh shell, add export %s=<your-api-key> to your ~/.zsh_profile
- for fish shell, add set -x %s <your-api-key> to your ~/.config/fish/config.fish
`, ApiKeyFilePath, API_KEY_ENV_VAR_NAME, API_KEY_ENV_VAR_NAME, API_KEY_ENV_VAR_NAME)
)

var (
	ErrApiTokenFileDoesntExist        = errors.New(fmt.Sprintf("API key file %s doesn't exist", ApiKeyFilePath))
	ErrCantOpenApiTokenFileForReading = errors.New(fmt.Sprintf("Can not open API key file %s for reading", ApiKeyFilePath))
)

const (
	APP_NAME             = "chat"
	API_KEY_ENV_VAR_NAME = "CHAT_GPT_API_TOKEN"
	API_KEY_SOURCE_ENV   = "env"
	API_KEY_SOURCE_FILE  = "file"
)

type Config struct {
	Renderer renderer.Renderer
	Format   string
}

func inSlice(what string, slice []string) bool {
	for _, i := range slice {
		if what == i {
			return true
		}
	}
	return false
}
func LoadConfig(optionsOverride []string) *Config {
	var r renderer.Renderer
	var format string
	if inSlice("--md", optionsOverride) {
		r = renderer.NewMarkdownRenderer(os.Stdout, "ChatGPT: ")
		format = "markdown"
	} else {
		r = renderer.NewPassthruRenderer(os.Stdout, "ChatGPT: ")
		format = ""
	}
	return &Config{
		Renderer: r,
		Format:   format,
	}
}

func LoadApiKey() (string, string, error) {
	apiKey := os.Getenv(API_KEY_ENV_VAR_NAME)
	if apiKey != "" {
		return apiKey, API_KEY_SOURCE_ENV, nil
	}
	_, err := os.Stat(ApiKeyFilePath)
	if err != nil {
		return "", "", ErrApiTokenFileDoesntExist
	}
	file, err := os.Open(ApiKeyFilePath)
	if err != nil {
		return "", "", ErrCantOpenApiTokenFileForReading
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)
	scanner.Scan()
	return scanner.Text(), API_KEY_SOURCE_FILE, nil
}

func getCacheDir() string {
	var dir string
	switch runtime.GOOS {
	case "darwin", "linux", "freebsd", "openbsd":
		dir = os.Getenv("XDG_CACHE_HOME")
		if dir != "" {
			return dir
		}
		home, _ := os.UserHomeDir()
		return path.Join(home, ".cache", APP_NAME)
	default:
		dir, _ := os.UserConfigDir()
		return dir
	}
}

func getConfigDir() string {
	var dir string
	switch runtime.GOOS {
	case "darwin", "linux", "freebsd", "openbsd":
		dir = os.Getenv("XDG_CONFIG_HOME")
		if dir != "" {
			return dir
		}
		home, _ := os.UserHomeDir()
		return path.Join(home, ".config", APP_NAME)
	default:
		dir, _ := os.UserConfigDir()
		return dir
	}
}

func getApiKeyFilePath() string {
	home, _ := os.UserHomeDir()
	return path.Join(home, "."+APP_NAME)
}

func StoreApiToken(token string) bool {
	return true
}
