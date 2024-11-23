package cmd

import (
	"fmt"
	"time"

	"github.com/AmirMirzayi/relay/host"
	"github.com/spf13/cobra"
)

var sendCmd = &cobra.Command{
	Use:     "send",
	Short:   "send file(s)",
	Example: "relay send [-p 12345 | -i 127.0.0.1 | -t 120 | -w 25] some_file.ext other_file2.ext some_directory_within_subdirectories",
	Args:    cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if 100%progressbarWidth != 0 {
			return fmt.Errorf("%d is not divisible to 100", progressbarWidth)
		}

		return host.Serve(ip,
			port,
			progressbarWidth,
			time.Duration(timeoutInSecond)*time.Second,
			silentTransfer,
			args...)
	},
}
