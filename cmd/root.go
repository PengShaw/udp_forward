package cmd

import (
	"errors"
	"net/netip"

	"github.com/spf13/cobra"

	"github.com/PengShaw/udp_forward/forward"
)

var source string
var destinations []string
var showTest bool

var rootCmd = &cobra.Command{
	Use:   "udp_forward",
	Short: "udp_forward can send udp data from one source to mutli destinations",
	RunE: func(cmd *cobra.Command, args []string) error {
		sourceAddr, err := netip.ParseAddrPort(source)
		if err != nil {
			return errors.New("source address format should be like a.b.c.d:p")
		}
		destinationAddrs := []netip.AddrPort{}
		for _, destination := range destinations {
			destinationAddr, err := netip.ParseAddrPort(destination)
			if err != nil {
				return errors.New("destination address format should be like a.b.c.d:p")
			}
			destinationAddrs = append(destinationAddrs, destinationAddr)
		}

		return forward.Run(sourceAddr, destinationAddrs, showTest)
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.Flags().StringVarP(&source, "listen", "l", "", "listen address & port for udp source, e.g., 0.0.0.0:514")
	rootCmd.Flags().StringArrayVarP(&destinations, "destinations", "d", []string{""}, "destinations for udp data, e.g., 192.168.1.2:9000")
	rootCmd.Flags().BoolVarP(&showTest, "test", "t", false, "print data to std.out for testing")
	rootCmd.MarkFlagRequired("listen")
	rootCmd.MarkFlagRequired("destinations")
}
