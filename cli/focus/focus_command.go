package focus

import (
	"dns-server/server"
	"dns-server/server/dto"
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
				return
			}

			flag := cmd.Flag("file")
			if flag.Value.String() == "" {
				cmd.Usage()
				return
			}

			dnsConfig, err = fileToDnsConfig(flag.Value.String())
			if err != nil {
				log.Printf("Error %s", err.Error())
				cmd.Usage()
				return
			}

			server.Start(dnsConfig)
		},
	}
	var Path string
	commandAddCmd.Flags().StringVarP(&Path, "file", "f", "", "file to read")
	
	return commandAddCmd
}

func fileToDnsConfig(path string) (*dto.DnsConfig, error) {
	file, err := os.Open(path)
	if err != nil {
		fmt.Printf("Error opening the file: %v\n", err)
		return &dto.DnsConfig{}, err
	}
	defer file.Close()
	
	content, err := io.ReadAll(file)
	if err != nil {
		fmt.Printf("Error reading the file: %v\n", err)
		return &dto.DnsConfig{}, err
	}

	var dnsConfig dto.DnsConfig
	if err := json.Unmarshal(content, &dnsConfig); err != nil {
		return &dto.DnsConfig{}, errors.New("json_format_error")
	}

	return &dnsConfig, nil
}