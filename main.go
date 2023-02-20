/*
Copyright © 2023 Kirill Kovalchuk <the.nemoden@gmail.com>
*/
package main

import (
	"context"
	"os"
	"strings"

	"github.com/nemoden/chat/chat"
	"github.com/nemoden/chat/cmd"
	"github.com/nemoden/chat/config"
	"github.com/nemoden/chat/util"
)

func main() {
	args := os.Args[1:]
	options := getOptions(args)
	positionalArgs := getPotisionalArgs(args)
	if len(positionalArgs) > 0 {
		subcommands := cmd.GetSubcommandsMap()
		// no positional args
		_, ok := subcommands[positionalArgs[0]]
		// we have no such subcommand, chat implied.
		if !ok {
			c := config.LoadConfig(options)
			input := strings.Join(positionalArgs, " ")
			client, err := chat.NewClient(c)
			if err != nil {
				util.Hangup(err)
			}
			ctx := context.Background()
			request := chat.CompletionRequest(input, c)
			response, _ := client.GptClient().CreateCompletionStream(ctx, request)
			c.Renderer.Render(response)
			os.Exit(0)
		}
	}
	cmd.Execute()
}

func getOptions(args []string) []string {
	var opts []string
	for _, arg := range args {
		if strings.HasPrefix(arg, "--") {
			opts = append(opts, arg)
		}
	}
	return opts
}

func getPotisionalArgs(args []string) []string {
	i := 0
	for idx, arg := range args {
		if strings.HasPrefix(arg, "--") {
			i = idx + 1
		} else {
			break
		}
	}
	if len(args) >= i {
		return args[i:]
	}
	return []string{}
}
