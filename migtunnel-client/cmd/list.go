package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"golang/migtunnel-client/util"
)

var adminPort int16

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all tunnels registered",
	Run: func(cmd *cobra.Command, args []string) {
		if adminPort == -1 {
			cmd.Println("Please provide a valid port")
			cmd.Println(cmd.UsageString())
			return
		}
		adminUrl := fmt.Sprintf("http://localhost:%d/list", adminPort)
		allTunnels, err := util.GetTunnels(adminUrl)
		if err != nil {
			cmd.Println("Unable to list tunnels , please check if jtunnel is up and "+
				"ui server is running on port ", adminPort)
		}
		cmd.Println(allTunnels)
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().Int16VarP(&adminPort, "adminPort", "p", -1, "Admin Server port")
}
