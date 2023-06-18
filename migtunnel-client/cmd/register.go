package cmd

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"golang/migtunnel-client/data"
	"golang/migtunnel-client/util"
)

var hostName string
var localServerPort int16
var adminServerPort int16
var localTlsServer string

var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "Register a new tunnel",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		fmt.Println(args)
		if hostName == "" || localServerPort == 0 || adminServerPort == 0 {
			cmd.PrintErrln("Please provide a valid data for hostName,localServerPort and adminServerPort")
			cmd.Println(cmd.UsageString())
			return
		}
		tunnelRegisterRequest := data.TunnelCreateRequest{
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
	registerCmd.Flags().Int16Var(&adminServerPort, "adminPort", 0, "Port on which the ui server is running")

	tlsSubCommand := &cobra.Command{
		Use:   "tls",
		Short: "Create a new tunnel to a local server running TLS ",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			if localTlsServer == "" {
				cmd.PrintErrln("Please provide the fully qualified domain for the tls server")
				cmd.Println(cmd.UsageString())
				return
			}
			if localServerPort == 0 {
				cmd.PrintErrln("Please prove the port on which the server is running")
				cmd.Println(cmd.UsageString())
				return
			}
			if adminServerPort == 0 {
				cmd.PrintErrln("Please provide a valid adminServerPort")
				cmd.Println(cmd.UsageString())
				return
			}
			tunnelRegisterRequest := data.TunnelCreateRequest{
				HostName:        hostName,
				TunnelName:      uuid.NewString(),
				LocalServerPort: localServerPort,
				Tls:             true,
				TlsServerFQDN:   localTlsServer,
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
	tlsSubCommand.Flags().StringVar(&localTlsServer, "server", "", "Fully qualified domain of the local server")
	tlsSubCommand.Flags().Int16Var(&localServerPort, "port", 0, "Port on which the local server is running")
	tlsSubCommand.Flags().Int16Var(&adminServerPort, "adminPort", 0, "Port on which the ui server is running")
	tlsSubCommand.Flags().StringVar(&hostName, "host", "", "Host name for the tunnel ")
	tlsSubCommand.MarkFlagRequired("server")
	tlsSubCommand.MarkFlagRequired("port")
	tlsSubCommand.MarkFlagRequired("adminPort")
	tlsSubCommand.MarkFlagRequired("host")
	registerCmd.AddCommand(tlsSubCommand)

}
