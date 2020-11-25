package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"net"
)

func init() {
	rootCmd.AddCommand(devCmd)
}

var devCmd = &cobra.Command{
	Use:   "dev",
	Short: "For testing on dev various things",
	Long:  "",
	//Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
		if err != nil {
			//
		}

		listener, err := net.ListenTCP("tcp", addr)
		if err != nil {
			//
		}
		defer listener.Close()
		logrus.Info(listener.Addr().(*net.TCPAddr).Port)

	},
}
