package command

import (
	hostDns "dns-focus/src/host-dns"
	"log"

	"github.com/spf13/cobra"
)


func InitHostCommand() *cobra.Command {
	var setDnsFlag bool
	var resetDnsFlag bool
	var hostCommand = &cobra.Command{
		Use:   "host",
		Short: "(only for mac) set/reset your dns",
		DisableFlagsInUseLine: true,
		Run: func(cmd *cobra.Command, args []string) {


			if setDnsFlag && resetDnsFlag {
				log.Printf("Error you can use only set or reset not both")
				cmd.Usage()
				return
			} else if setDnsFlag {
				hostDns.SetDnsForMac()
				return
			} else if resetDnsFlag {
				hostDns.ResetDnsForMac()
				return
			}
			cmd.Usage()
		},
	}
	hostCommand.Flags().BoolVar(&setDnsFlag, "set", false, "Set your dns to 127.0.0.1")
	hostCommand.Flags().BoolVar(&resetDnsFlag, "reset", false, "Reset your dns to the last backup file in host-dns/backup")
	
	return hostCommand
}
