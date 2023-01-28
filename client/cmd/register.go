package cmd

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"golang/client/admin/data"
	tunnels2 "golang/client/admin/tunnels"
)

var hostName string
var port int16

var logger, _ = zap.NewProduction()
var sugar = logger.Sugar()

// registerCmd represents the register command
var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "Register a new tunnel",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("register called ", hostName)

		tunnelRegisterRequest := data.TunnelCreateRequest{
			HostName:        hostName,
			TunnelName:      uuid.NewString(),
			LocalServerPort: port,
		}

		err := tunnels2.RegisterTunnel(tunnelRegisterRequest)
		if err != nil {
			sugar.Error("Failed to register tunnel for request ",
				tunnelRegisterRequest, err.Error())
		}
	},
}

func init() {
	rootCmd.CompletionOptions.HiddenDefaultCmd = true
	rootCmd.AddCommand(registerCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// registerCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// registerCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	registerCmd.Flags().StringVar(&hostName,
		"host", "", "Host name for the tunnel ")
	registerCmd.Flags().Int16Var(&port,
		"port", 0, "Port on which the local server is running")

}
