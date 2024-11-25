// Package cmd implements command line arguments to use application
package cmd

import (
	"os"

	"github.com/AmirMirzayi/relay/config"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "relay",
	Short: "relay is an CLI application wich provide transfer files over network.",
	Long: `relay is a CLI application for transferring files over network(local or global).
This application will transfer folders within files as exists in source machine.
credit: Amir Mirzaei
mirzayi994@gmail.com
https://github.com/AmirMirzayi
https://www.linkedin.com/in/amir-mirzaei
	`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var setting config.Setting

func init() {
	rootCmd.PersistentFlags().IntVarP(&setting.Port, "port", "p", config.DefaultPort, "application running port")
	rootCmd.PersistentFlags().DurationVarP(&setting.Timeout, "timeout", "t", config.DefaultTimeout, "connection timeout")
	rootCmd.PersistentFlags().IntVarP(&setting.ProgressbarWidth, "width", "w", config.DefaultProgressbarWidth, "progress bar width[must divisible to 100]")
	rootCmd.PersistentFlags().BoolVarP(&setting.SilentTransfer, "silent", "l", config.DefaultSilent, "silent transfer")
	rootCmd.PersistentFlags().IntVarP(&setting.BufferSize, "buffer", "b", config.DefaultBufferSize, "buffer size in byte")

	sendCmd.PersistentFlags().IPVarP(&setting.IP, "ip", "i", config.DefaultIP(), "sender machine binding ip address")
	receiveCmd.PersistentFlags().IPVarP(&setting.IP, "ip", "i", nil, "sender machine ip address")
	receiveCmd.MarkPersistentFlagRequired("ip")
	receiveCmd.PersistentFlags().StringVarP(&setting.SavePath, "save", "s", config.DefaultDirectory(), "files save path")

	rootCmd.AddCommand(receiveCmd)
	rootCmd.AddCommand(sendCmd)
}
