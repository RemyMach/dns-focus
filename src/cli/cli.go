package cli

import (
	"dns-focus/src/cli/command"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "dns-focus",
	Short: "dns-focus is a CLI for start a focus mode, start your own dns server with custom blocking domains ...",
	
}

func InitCli() {
	focusCommand := command.InitFocusCommand()
	hostCommand := command.InitHostCommand()
	rootCmd.AddCommand(focusCommand)
	rootCmd.AddCommand(hostCommand)
	rootCmd.CompletionOptions.DisableDefaultCmd = true

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
