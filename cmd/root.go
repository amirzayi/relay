// Package cmd implements command line arguments to use application
package cmd

import (
	"net"
	"os"
	"time"

	"github.com/AmirMirzayi/relay/config"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "relay",
	Short: "relay is an CLI application wich provide transfer files over network.",
	Long: `relay is a CLI application for transferring files over network(local or global).
This application will transfer folders within files as exists in source machine.
	`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var (
	ip net.IP
	port,
	timeoutInSecond,
	progressbarWidth int
	savePath       string
	silentTransfer bool
)

func init() {
	rootCmd.PersistentFlags().IntVarP(&port, "port", "p", config.DefaultPort, "application running port")
	rootCmd.PersistentFlags().IntVarP(&timeoutInSecond, "timeout", "t", int(config.DefaultTimeout/time.Second), "connection timeout in second")
	rootCmd.PersistentFlags().IntVarP(&progressbarWidth, "width", "w", config.DefaultProgressbarWidth, "progress bar width[must divisible to 100]")
	rootCmd.PersistentFlags().BoolVarP(&silentTransfer, "silent", "l", config.DefaultSilent, "silent transfer")

	sendCmd.PersistentFlags().IPVarP(&ip, "ip", "i", config.DefaultIP(), "sender machine binding ip address")
	receiveCmd.PersistentFlags().IPVarP(&ip, "ip", "i", nil, "sender machine ip address")
	receiveCmd.MarkPersistentFlagRequired("ip")
	receiveCmd.PersistentFlags().StringVarP(&savePath, "save", "s", config.DefaultDirectory(), "files save path")

	rootCmd.AddCommand(receiveCmd)
	rootCmd.AddCommand(sendCmd)
}
