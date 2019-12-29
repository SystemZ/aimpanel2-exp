package cmd

import (
	"github.com/spf13/cobra"
	"gitlab.com/systemz/aimpanel2/slave/model"
	"gitlab.com/systemz/aimpanel2/slave/tasks"
)

func init() {
	rootCmd.AddCommand(cmdCmd)
	cmdCmd.Flags().StringVarP(&cmdBody, "cmd", "c", "", "-c \"say I like\"")
	cmdCmd.MarkFlagRequired("cmd")
}

var (
	cmdBody string
)

var cmdCmd = &cobra.Command{
	Use:   "cmd -c <CMD> <GS ID1> <GS ID2> ...",
	Short: "Run console command on game server",
	Long:  "Execute cmd on game server console\nUse \"all\" to send this to all GS on this host",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		model.InitRedis()
		for _, gsId := range args {
			tasks.GsCmd(gsId, cmdBody)
		}
	},
}
