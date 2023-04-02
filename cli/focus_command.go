package cli

import (
	"log"

	"github.com/spf13/cobra"
)


func initFocusCommand() *cobra.Command {
	var commandAddCmd = &cobra.Command{
		Use:   "focus",
		Short: "Start your dns server in focus mode",
		Long: `no need log description`,
		DisableFlagsInUseLine: true,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				cmd.Usage()
				return
			}

			jsonHandler, err := urlshort.JSONHandler(fileContent, yamlHandler)
			if err != nil {
				log.Println("[Error]: the file is not valid in the json")
			}
			
			//nameTask := strings.Join(args[:], " ")
			//usecases.AddCommand(nameTask)
		},
	}
	
	return commandAddCmd
}

func fileTo