package cmd

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/nemoden/gogpt/config"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		c := config.LoadConfig([]string{})
		askUserToUpdateConfig(c)
	},
}

// Fill the gaps in config, only ask a user what's not been configured, or ask to override.
func askUserToUpdateConfig(c *config.Config) {
	var newApiKey string
	userProvidedApiKey := false
	hasApiKey := false
	if c.ApiKey().IsEmpty() {
		if askIfUserWantsKeyStoredInFile() {
			newApiKey = askApiKey()
			userProvidedApiKey = true
		}
	} else {
		switch c.ApiKey().Source {
		case config.SourceEnv:
			notifyUserTheirKeyStoreInEnv(c.ApiKey())
		case config.SourceFile:
			if askIfUserWantsToReplaceApiKey(c.ApiKey()) {
				newApiKey = askApiKey()
				userProvidedApiKey = true
			}
		}
	}

	if userProvidedApiKey {
		_, err := config.StoreApiKey(newApiKey)
		if err != nil {
			fmt.Printf("We were unable to store the API key you have provided: %s", err)
		}
	}

	askToUpdateConfigValues(c)

	hasApiKey = !c.ApiKey().IsEmpty() || userProvidedApiKey
	if !hasApiKey {
		fmt.Printf("API key is not configured. To configure the API key run this command again.")
	}
}

func notifyUserTheirKeyStoreInEnv(key config.ApiKey) {
	fmt.Printf(`Your API key %s is stored in the environment variable %s.

If you wish to change it, just set it to a different value.

If you wish to store the API key in %s instead,
delete the %s environment variable first,
and re-run the configuration command.

Alternatively, you can always either set the environment variable or edit the file manually
`,
		key.Mask(),
		config.API_KEY_ENV_VAR_NAME,
		config.ApiKeyFilePath,
		config.API_KEY_ENV_VAR_NAME)
}

func askIfUserWantsKeyStoredInFile() bool {
	var option string
	survey.AskOne(&survey.Select{
		Message: fmt.Sprintf("Your API key can be stored permanently in %s file or in the %s environment variable. Would you like to store it in file or provide via env var later?", config.ApiKeyFilePath, config.API_KEY_ENV_VAR_NAME),
		Options: []string{"File", "Environment Variable"},
		Default: "File",
	}, &option)
	return option == "File"
}

func askIfUserWantsToReplaceApiKey(key config.ApiKey) bool {
	var replaceApiKey bool
	survey.AskOne(&survey.Confirm{Message: fmt.Sprintf("We have your API key configured as %s. Would you like to replace it?", key.Mask()), Default: false}, &replaceApiKey)
	return replaceApiKey
}

func askApiKey() string {
	var apiKey string
	survey.AskOne(&survey.Password{Message: fmt.Sprintf("An API key (that will be stored in %s): ", config.ApiKeyFilePath)}, &apiKey)
	return apiKey
}

func maskApiToken(t string) string {
	head := 5
	tail := 4
	mask := 5
	if len(t) < (head + tail + mask + 1) {
		if len(t) < 4 {
			return strings.Repeat("*", len(t))
		}
		return strings.Join([]string{t[:1], strings.Repeat("*", len(t)-1)}, "")
	}
	return strings.Join([]string{
		t[:head],
		strings.Repeat("*", mask),
		t[len(t)-tail+1:],
	}, "")
}

func askToUpdateConfigValues(c *config.Config) {
	askToUpdateRenderer(c)
	d, err := yaml.Marshal(c)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = os.MkdirAll(config.ConfigDir, os.ModePerm)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = os.WriteFile(path.Join(config.ConfigDir, "config.yaml"), d, 0755)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("All good, your config has been saved up. Happy chatting!")
}

type renderer struct {
	Title   string
	Comment string
}

func askToUpdateRenderer(c *config.Config) {
	options := map[string]renderer{
		"plain": {Title: "Plain text", Comment: "No formatting, output passed through as-is"},
		"md":    {Title: "Markdown", Comment: "Experimental. Renders responses as markdown"},
		"md2":   {Title: "Markdown 2", Comment: "Experimental, enhanced. Renders responses as markdown"},
	}
	var current string

	titles := make([]string, len(options))
	idxToRef := make(map[int]string, len(options))
	i := 0
	for ref, r := range options {
		titles[i] = r.Title
		idxToRef[i] = ref
		if ref == c.RendererRef {
			current = r.Title
		}
		i++
	}
	answerIndex := 0
	err := survey.AskOne(&survey.Select{
		Message: fmt.Sprintf("Choose formatter (you can change it later)"),
		Options: titles,
		Description: func(value string, index int) string {
			return options[idxToRef[index]].Comment
		},
		Default: current,
	}, &answerIndex)
	if err != nil {
		fmt.Println(err)
		return
	}
	c.RendererRef = idxToRef[answerIndex]
}

func init() {
	rootCmd.AddCommand(configCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// configCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	configCmd.Flags().BoolP("print", "p", false, "Prints current config")
}
