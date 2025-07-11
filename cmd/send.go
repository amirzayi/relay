package cmd

import (
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/amirzayi/relay/config"
	"github.com/amirzayi/relay/pkg/fileutil"
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
		return Serve(setting, args...)
	},
}

// Serve starts a file transfer server that listens on the specified IP and port.
// It serves files located at the provided paths to connected clients.
func Serve(setting config.Setting, paths ...string) error {
	files, err := fileutil.GetFilesByPaths(paths...)
	if err != nil {
		return err
	}

	fmt.Printf("Preparing to send %d files with %s\n", len(files), files.HumanReadableTotalSize())

	listener, err := net.ListenTCP("tcp", &net.TCPAddr{IP: setting.IP, Port: setting.Port})
	if err != nil {
		return err
	}
	defer listener.Close()
	listener.SetDeadline(time.Now().Add(setting.Timeout))

	conn, err := listener.AcceptTCP()
	if err != nil {
		return err
	}
	defer conn.Close()

	if err = files.SendDetails(conn); err != nil {
		return err
	}

	for i, file := range files {
		_, err = file.ProgressiveRead(conn, func(transferred int64, percent int) {
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
