package cli

import (
	"dns-server/cli/focus"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "dns-server",
	Short: "dns-server is a CLI for start a focus mode, start your own dns server with custom blocking ips ...",
	
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
