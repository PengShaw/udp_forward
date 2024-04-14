package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/PengShaw/udp_forward/forward"
)

var source string
var destinations []string
var v bool
var vv bool

var rootCmd = &cobra.Command{
	Use:   "udp_forward",
	Short: "udp_forward received data from one source socket (udp/tcp/unix), and send to multi destination sockets (udp/tcp/unix)",
	Run: func(cmd *cobra.Command, args []string) {
		forward.Run(source, destinations, v, vv)
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of udp_forward",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("udp_forward version: v0.0.2")
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.Flags().StringVarP(&source, "listen", "l", "", "listen, e.g., udp:0.0.0.0:514 or tcp:0.0.0.0:514 or unix:/path/to/unix.sock")
	rootCmd.Flags().StringArrayVarP(&destinations, "destinations", "d", []string{""}, "destinations, e.g., udp:192.168.1.2:9000 or tcp:0.0.0.0:514 or unix:/path/to/unix.sock")
	rootCmd.Flags().BoolVarP(&v, "verbose", "v", false, "print info log")
	rootCmd.Flags().BoolVar(&vv, "vv", false, "more verbose, print debug log")

	rootCmd.MarkFlagRequired("listen")
	rootCmd.MarkFlagRequired("destinations")

	rootCmd.AddCommand(versionCmd)
}
