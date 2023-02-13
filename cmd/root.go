/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"strings"

	"os"

	"github.com/nemoden/chat/chat"
	"github.com/nemoden/chat/config"
	gogpt "github.com/sashabaranov/go-gpt3"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "chat",
	Short: "ChatGPT CLI and REPL",
	Long:  `Use ChatGPT not leaving your terminal`,
	FParseErrWhitelist: cobra.FParseErrWhitelist{
		UnknownFlags: true,
	},
	Run: func(cmd *cobra.Command, args []string) {
		reader := bufio.NewReader(os.Stdin)

		options := getOptions(os.Args[1:])
		c := config.LoadConfig(options)
		client, err := chat.NewClient(c)

		if errors.Is(err, config.ErrApiTokenFileDoesntExist) {
			fmt.Printf(config.TokenFileDoesntExistPrompt)
			os.Exit(0)
		}
		if errors.Is(err, config.ErrCantOpenApiTokenFileForReading) {
			fmt.Printf(config.TokenFileNotReadablePrompt)
			os.Exit(0)
		}

		ctx := context.Background()

		for {
			fmt.Printf("You: ")
			input, _ := reader.ReadString('\n')
			input = strings.TrimSpace(input)

			if input == "exit" {
				break
			}

			prompt := ""
			if c.Format == "markdown" {
				prompt += "Return response in markdown format. Prompt on a new line:\n"
			}
			prompt += input
			req := gogpt.CompletionRequest{
				Model:       gogpt.GPT3TextDavinci003,
				Prompt:      prompt,
				MaxTokens:   1000,
				Temperature: 0.5,
			}

			fmt.Printf("ChatGPT: ")

			stream, _ := client.GptClient().CreateCompletionStream(ctx, req)

			defer stream.Close()

			c.Renderer.Render(stream)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.chat.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// @TODO DRY, duplication in ./main.go
func getOptions(args []string) []string {
	var opts []string
	for _, arg := range args {
		if strings.HasPrefix(arg, "--") {
			opts = append(opts, arg)
		}
	}
	return opts
}
