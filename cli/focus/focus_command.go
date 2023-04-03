package focus

import (
	"dns-focus/common"
	"dns-focus/config"
	"dns-focus/server"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/spf13/cobra"
)


func InitFocusCommand() *cobra.Command {
	var proxyFlag bool
	var focusServerMode common.ServerMode
	var focusCommand = &cobra.Command{
		Use:   "focus",
		Short: "Start your dns server in focus mode",
		DisableFlagsInUseLine: true,
		Run: func(cmd *cobra.Command, args []string) {


			var filePath string
			flag := cmd.Flag("file")
			if flag.Value.String() == "" {
				filePath = "config/config.json"
			} else {
				filePath = flag.Value.String()
			}

			dnsConfig, err := fileToDnsConfig(filePath)
			if err != nil {
				log.Printf("Error %s", err.Error())
				cmd.Usage()
				return
			}

			if proxyFlag {
				focusServerMode = common.Proxy
			} else {
				focusServerMode = common.Server
			}

			server.Start(dnsConfig, focusServerMode)
		},
	}
	var Path string
	focusCommand.Flags().StringVarP(&Path, "file", "f", "", "file to read")
	focusCommand.Flags().BoolVar(&proxyFlag, "proxy", false, "Use proxy")
	
	return focusCommand
}

func fileToDnsConfig(path string) (*config.DnsConfig, error) {
	file, err := os.Open(path)
	if err != nil {
		fmt.Printf("Error opening the file: %v\n", err)
		return &config.DnsConfig{}, err
	}
	defer file.Close()
	
	content, err := io.ReadAll(file)
	if err != nil {
		fmt.Printf("Error reading the file: %v\n", err)
		return &config.DnsConfig{}, err
	}

	var dnsConfig config.DnsConfig
	if err := json.Unmarshal(content, &dnsConfig); err != nil {
		return &config.DnsConfig{}, errors.New("json_format_error")
	}

	return &dnsConfig, nil
}