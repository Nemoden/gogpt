package main

import (
	"os"
	"strings"

	"github.com/nemoden/gogpt/chat"
	"github.com/nemoden/gogpt/cmd"
	"github.com/nemoden/gogpt/config"
	"github.com/nemoden/gogpt/util"
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
			ch, err := chat.New(c)
			if err != nil {
				util.Hangup(err)
			}
			err = ch.AskAndRender(input)
			if err != nil {
				util.Hangup(err)
			}
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
