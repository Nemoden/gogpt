/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"context"
	"fmt"
	"strings"

	"os"

	"github.com/nemoden/gogpt/chat"
	"github.com/nemoden/gogpt/config"
	"github.com/nemoden/gogpt/util"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gogpt",
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

		if err != nil {
			util.Hangup(err)
		}

		ctx := context.Background()

		for {
			fmt.Printf("You: ")
			input, _ := reader.ReadString('\n')
			input = strings.TrimSpace(input)

			if input == "exit" {
				break
			}

			req := chat.CompletionRequest(input, c)

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
