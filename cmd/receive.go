package cmd

import (
	"fmt"
	"net"
	"strings"

	"github.com/amirzayi/relay/config"
	"github.com/amirzayi/relay/pkg/fileutil"
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
		return Receive(setting)
	},
}

// Receive connects to a file transfer server at the specified IP and port to receive files.
// It handles file reception and displays a progress bar.
func Receive(setting config.Setting) error {
	serverAddress := fmt.Sprintf("%s:%d", setting.IP, setting.Port)
	fmt.Printf("Connecting to %s...", serverAddress)
	conn, err := net.DialTimeout("tcp", serverAddress, setting.Timeout)
	if err != nil {
		return err
	}
	defer conn.Close()
	fmt.Printf("\rSuccessfully Connected to %s ✓\n", serverAddress)

	files, err := fileutil.GetDetails(conn)
	if err != nil {
		return err
	}

	fmt.Printf("Preparing to receive %d files with %s\n", len(files), files.HumanReadableTotalSize())

	for i, file := range files {
		_, err = file.ProgressiveWrite(conn, func(transferred int64, percent int) {
			barTitle := fmt.Sprintf("<%s ^ %s>", file.HumanReadableSize(), file.ShortedName(10, 8, 3))
			width := setting.ProgressbarWidth
			progress := strings.Repeat("", 100-percent/(100/width))
			progress += strings.Repeat("█", percent/(100/width))

			fmt.Printf("\r%d%% |%-*s| %s", percent, width, progress, barTitle)
		})
		if err != nil {
			return err
		}
		fmt.Printf("\r[%d] %s ✓\033[K\n", i, file.Name)
	}

	return nil
}
