package config

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path"
	"runtime"
	"strings"

	"github.com/nemoden/chat/renderer"
	gogpt "github.com/sashabaranov/go-gpt3"
	"gopkg.in/yaml.v2"
)

var (
	ConfigDir      string = getConfigDir()
	CacheDir       string = getCacheDir()
	ApiKeyFilePath string = getApiKeyFilePath()
	InitPrompt     string = `Please initialise chat using the config command:

$ chat config
`
	// @TODO revisit instructions
	TokenFileDoesntExistPrompt string = fmt.Sprintf(`Oops. Looks like your token file %s doesn't exist.

Please obtain the openai API key from their website https://platform.openai.com/account/api-keys

Once you have the API key, you can either add it manually to %s:

$ echo "<your-api-key>" > %s

Or, if you didn't run chat config, it's adviced that you do this instead:

$ chat config

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
)

type ApiKeySource int
type Format int

const (
	FormatMarkdown Format = iota
	FormatPlain
)

const (
	SourceNone ApiKeySource = iota
	SourceEnv
	SourceFile
)

type ApiKey struct {
	key    string
	Source ApiKeySource
}

func (ak ApiKey) String() string {
	return ak.key
}

func (ak ApiKey) IsEmpty() bool {
	return ak.key == ""
}

func (ak ApiKey) Mask() string {
	head := 5
	tail := 4
	mask := 5
	if len(ak.key) < (head + tail + mask + 1) {
		if len(ak.key) < 4 {
			return strings.Repeat("*", len(ak.key))
		}
		return strings.Join([]string{ak.key[:1], strings.Repeat("*", len(ak.key)-1)}, "")
	}
	return strings.Join([]string{
		ak.key[:head],
		strings.Repeat("*", mask),
		ak.key[len(ak.key)-tail+1:],
	}, "")
}

type Config struct {
	Renderer     renderer.Renderer `yaml:"-"`
	RendererRef  string            `yaml:"renderer,omitempty"`
	Format       Format            `yaml:"-"`
	PromptPrefix string            `yaml:"-"`
	Model        string            `yaml:"model,omitempty"`
	MaxTokens    int               `yaml:"max_tokens,omitempty"`
	Temperature  float32           `yaml:"temperature,omitempty"`
	apiKey       ApiKey            `yaml:"-"`
}

// :grin:
func inSlice(what string, slice []string) bool {
	for _, i := range slice {
		if what == i {
			return true
		}
	}
	return false
}

func (c *Config) load() error {
	b, err := os.ReadFile(path.Join(ConfigDir, "config.yaml"))
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(b, c)
	if err != nil {
		return err
	}
	return nil
}

func LoadConfig(optionsOverride []string) *Config {
	c := Config{}
	c.load()
	if inSlice("--md", optionsOverride) || c.RendererRef == "md" {
		c.Renderer = renderer.NewMarkdownRenderer(os.Stdout, "ChatGPT: ")
		c.Format = FormatMarkdown
		c.RendererRef = "md"
	} else if inSlice("--md2", optionsOverride) || c.RendererRef == "md2" {
		c.Renderer = renderer.NewMarkdown2Renderer(os.Stdout, "ChatGPT: ")
		c.Format = FormatMarkdown
		c.RendererRef = "md2"
	} else if inSlice("--token", optionsOverride) || c.RendererRef == "token" {
		c.Renderer = renderer.NewTokenRenderer(os.Stdout, "ChatGPT: ")
		c.Format = FormatPlain
		c.RendererRef = "token"
	} else {
		c.Renderer = renderer.NewPassthruRenderer(os.Stdout, "ChatGPT: ")
		c.Format = FormatPlain
		c.RendererRef = "plain"
	}

	switch c.Format {
	case FormatMarkdown:
		c.PromptPrefix = "Return response in markdown format. Prompt on a new line:\n"
	default:
		c.PromptPrefix = ""
	}

	apiKey, _ := LoadApiKey()

	c.apiKey = apiKey

	// TODO later we want to provide this as configurable option.
	if c.Model == "" {
		c.Model = gogpt.GPT3TextDavinci003
	}

	if c.MaxTokens == 0 {
		c.MaxTokens = 1000
	}

	if c.Temperature == 0.0 {
		c.Temperature = 0.5
	}

	return &c
}

func (c *Config) UpdateApiKey(apiKey ApiKey) {
	c.apiKey = apiKey
}

func (c *Config) ApiKey() ApiKey {
	return c.apiKey
}

// Returns api key, source it's been pulled from, and error (if there was one)
func LoadApiKey() (ApiKey, error) {
	apiKey := os.Getenv(API_KEY_ENV_VAR_NAME)
	if apiKey != "" {
		return ApiKey{
			apiKey,
			SourceEnv,
		}, nil
	}
	_, err := os.Stat(ApiKeyFilePath)
	if err != nil {
		return ApiKey{}, ErrApiTokenFileDoesntExist
	}
	file, err := os.Open(ApiKeyFilePath)
	if err != nil {
		return ApiKey{}, ErrCantOpenApiTokenFileForReading
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)
	scanner.Scan()
	return ApiKey{
		scanner.Text(),
		SourceFile,
	}, nil
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

func StoreApiKey(token string) (ApiKey, error) {
	err := os.WriteFile(ApiKeyFilePath, []byte(token), 0755)
	if err == nil {
		return ApiKey{token, SourceFile}, nil
	}
	return ApiKey{}, err
}
