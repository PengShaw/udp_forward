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
var version bool

var rootCmd = &cobra.Command{
	Use:   "udp_forward",
	Short: "udp_forward can send udp data from one source to mutli destinations",
	RunE: func(cmd *cobra.Command, args []string) error {
		if version {
			fmt.Println("udp_forward version: v0.0.1")
			return nil
		}
		return forward.Run(source, destinations, v, vv)
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.Flags().StringVarP(&source, "listen", "l", "", "listen for udp data, e.g., udp:0.0.0.0:514 or unix:/path/to/unix.sock")
	rootCmd.Flags().StringArrayVarP(&destinations, "destinations", "d", []string{""}, "destinations for udp data, e.g., udp:192.168.1.2:9000 or unix:/path/to/unix.sock")
	rootCmd.Flags().BoolVarP(&v, "verbose", "v", false, "print info log")
	rootCmd.Flags().BoolVar(&vv, "vv", false, "more verbose, print debug log")
	rootCmd.Flags().BoolVar(&version, "version", false, "show version")

	rootCmd.MarkFlagRequired("listen")
	rootCmd.MarkFlagRequired("destinations")
}
