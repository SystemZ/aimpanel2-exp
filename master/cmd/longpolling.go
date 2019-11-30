package cmd

import (
	"fmt"
	"github.com/jcuga/golongpoll"
	"github.com/spf13/cobra"
	"log"
	"net/http"
)

func init() {
	rootCmd.AddCommand(longpollingCmd)
}

var longpollingCmd = &cobra.Command{
	Use:   "longpolling",
	Short: "For testing longpoll communication",
	Long:  "",
	//Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		manager, err := golongpoll.StartLongpoll(golongpoll.Options{
			LoggingEnabled: true,
		})
		if err != nil {
			log.Fatalf("Failed to create manager: %q", err)
		}

		// Serve our event subscription web handler
		http.HandleFunc("/events", manager.SubscriptionHandler)

		fmt.Println("Serving handler on 127.0.0.1:8081/events")
		http.ListenAndServe("127.0.0.1:8081", nil)
	},
}
