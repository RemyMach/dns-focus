package focus

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/spf13/cobra"
)


func InitFocusCommand() *cobra.Command {
	var commandAddCmd = &cobra.Command{
		Use:   "focus",
		Short: "Start your dns server in focus mode",
		Long: `no need log description`,
		DisableFlagsInUseLine: true,
		Run: func(cmd *cobra.Command, args []string) {

			dnsConfig, err := fileToDnsConfig("config/config.json")

			if err != nil {
				log.Printf("Error %s", err.Error())
				cmd.Usage()
			}

			log.Println(dnsConfig)
			
			//nameTask := strings.Join(args[:], " ")
			//usecases.AddCommand(nameTask)
		},
	}
	
	return commandAddCmd
}

func fileToDnsConfig(path string) (*DnsConfig, error) {
	file, err := os.Open(path)
	if err != nil {
		fmt.Printf("Error opening the file: %v\n", err)
		return &DnsConfig{}, err
	}
	defer file.Close()
	
	content, err := io.ReadAll(file)
	if err != nil {
		fmt.Printf("Error reading the file: %v\n", err)
		return &DnsConfig{}, err
	}

	var dnsConfig DnsConfig
	if err := json.Unmarshal(content, &dnsConfig); err != nil {
		return &DnsConfig{}, errors.New("json_format_error")
	}

	return &dnsConfig, nil
}