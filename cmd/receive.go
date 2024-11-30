package cmd

import (
	"fmt"

	"github.com/AmirMirzayi/relay/client"
	"github.com/spf13/cobra"
)

var receiveCmd = &cobra.Command{
	Use:     "receive",
	Aliases: []string{"r"},
	Short:   "receive file(s)",
	Example: "relay receive -i 127.0.0.1 [-p 12345 | -b 1024 | -t 120s | -w 25]",
	RunE: func(cmd *cobra.Command, args []string) error {
		if 100%setting.ProgressbarWidth != 0 {
			return fmt.Errorf("100 is not divisible by %d", setting.ProgressbarWidth)
		}
		return client.Receive(setting)
	},
}
