package cmd

import "github.com/spf13/cobra"

func GetRootCmd() *cobra.Command {
	return rootCmd
}

func GetSubcommandsMap() map[string]bool {
	commands := rootCmd.Commands()
	commandsLen := len(commands)
	subcommands := make(map[string]bool, commandsLen)
	for _, c := range commands {
		subcommands[c.Use] = true
	}
	return subcommands
}
