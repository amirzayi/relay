package cmd

import (
	"time"

	"github.com/AmirMirzayi/relay/client"
	"github.com/spf13/cobra"
)

var receiveCmd = &cobra.Command{
	Use:   "receive",
	Short: "receive file",
	RunE: func(cmd *cobra.Command, args []string) error {
		return client.Receive(ip, port, time.Duration(timeoutInSecond)*time.Second)
	},
}
