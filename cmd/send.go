package cmd

import (
	"fmt"

	"github.com/amirzayi/relay/host"
	"github.com/spf13/cobra"
)

var sendCmd = &cobra.Command{
	Use:     "send",
	Aliases: []string{"s"},
	Short:   "send file(s)",
	Example: "relay send [-p 12345 | -i 127.0.0.1 | -b 1024 | -t 120s | -w 25] some_file.ext other_file2.ext some_directory_within_subdirectories",
	Args:    cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if 100%setting.ProgressbarWidth != 0 {
			return fmt.Errorf("100 is not divisible by %d", setting.ProgressbarWidth)
		}
		return host.Serve(setting, args...)
	},
}
