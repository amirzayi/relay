package cmd

import (
	"fmt"
	"time"

	"github.com/AmirMirzayi/relay/client"
	"github.com/spf13/cobra"
)

var receiveCmd = &cobra.Command{
	Use:     "receive",
	Short:   "receive file(s)",
	Example: "relay receive -i 127.0.0.1 [-p 12345 | -t 120 | -w 25]",
	RunE: func(cmd *cobra.Command, args []string) error {
		if 100%progressbarWidth != 0 {
			return fmt.Errorf("%d is not divisible to 100", progressbarWidth)
		}

		return client.Receive(ip,
			port,
			bufferSize,
			progressbarWidth,
			time.Duration(timeoutInSecond)*time.Second,
			savePath,
			silentTransfer,
		)
	},
}
