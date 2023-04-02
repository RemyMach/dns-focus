package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"dns-server/cli/focus"
)

var rootCmd = &cobra.Command{
	Use:   "task",
	Short: "task is a CLI for managing your TODOs.",
	
}

func InitCli() {
	focusCommand := focus.InitFocusCommand()
	rootCmd.AddCommand(focusCommand)
	rootCmd.CompletionOptions.DisableDefaultCmd = true

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
