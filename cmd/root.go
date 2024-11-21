package cmd

import (
	"fmt"
	"net"
	"os"
	"time"

	"github.com/AmirMirzayi/relay/config"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "relay",
	Short: "relay is an cli tool wich provide transfer files over local network.",
	// 	Long: `A longer description that spans multiple lines and likely contains
	// examples and usage of using your application. For example:

	// Cobra is a CLI library for Go that empowers applications.
	// This application is a tool to generate the needed files
	// to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
	defer func() {
		r := recover()
		if r != nil {
			fmt.Printf("error occured: %v", r)
		}
	}()
}

var (
	ip              net.IP
	port            int
	timeoutInSecond int
)

func init() {

	rootCmd.PersistentFlags().IntVarP(&port, "port", "p", config.DefaultPort, "application running port")
	rootCmd.PersistentFlags().IntVarP(&timeoutInSecond, "timeout", "t", int(config.DefaultTimeout/time.Second), "connection timeout in second")

	sendCmd.PersistentFlags().IPVarP(&ip, "ip", "i", net.ParseIP(config.DefaultIP), "receiver machine ip address")
	sendCmd.MarkPersistentFlagRequired("ip")

	rootCmd.AddCommand(receiveCmd)
	rootCmd.AddCommand(sendCmd)

}
