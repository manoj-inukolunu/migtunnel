package cmd

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"golang/jtunnel-client/admin/data"
	"golang/jtunnel-client/util"
)

var hostName string
var localServerPort int16
var adminServerPort int16

var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "Register a new tunnel",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if hostName == "" || localServerPort == 0 || adminServerPort == 0 {
			cmd.PrintErrln("Please provide a valid data for hostName,localServerPort and adminServerPort")
			cmd.Println(cmd.UsageString())
			return
		}
		tunnelRegisterRequest := data.TunnelData{
			HostName:        hostName,
			TunnelName:      uuid.NewString(),
			LocalServerPort: localServerPort,
		}

		err := util.RegisterTunnel(fmt.Sprintf("http://localhost:%d/register", adminServerPort), tunnelRegisterRequest)
		if err != nil {
			cmd.PrintErrln("Could not register tunnel error is ", err.Error())
			return
		}
		cmd.Printf("Tunnel Successfully created. "+
			"All Requests to %s will now be routed to server running on %d", hostName+".migtunnel.net", localServerPort)

	},
}

func init() {
	rootCmd.CompletionOptions.HiddenDefaultCmd = true
	rootCmd.AddCommand(registerCmd)
	registerCmd.Flags().StringVar(&hostName,
		"host", "", "Host name for the tunnel ")
	registerCmd.Flags().Int16Var(&localServerPort,
		"port", 0, "Port on which the local server is running")
	registerCmd.Flags().Int16Var(&adminServerPort, "adminPort", 0, "Port on which the admin server is running")

}
