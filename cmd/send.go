package cmd

import (
	"time"

	"github.com/AmirMirzayi/relay/host"
	"github.com/spf13/cobra"
)

var sendCmd = &cobra.Command{
	Use:   "send",
	Short: "send file",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return host.Serve(ip,
			port,
			progressbarWidth,
			time.Duration(timeoutInSecond)*time.Second,
			args...)
	},
}
