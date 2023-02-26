package cmd

import (
	"context"
	markdown "github.com/MichaelMure/go-term-markdown"
	"github.com/spf13/cobra"
	"github.com/thejerf/suture/v4"
	"golang/jtunnel-client/client"
	"golang/jtunnel-client/data"
	"golang/jtunnel-client/db"
	"log"
)

var dbFilePath string

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Starts jtunnel",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		supervisor := suture.NewSimple("Client")
		service := &Main{cmd: cmd}
		ctx, cancel := context.WithCancel(context.Background())
		supervisor.Add(service)
		errors := supervisor.ServeBackground(ctx)
		cmd.PrintErrln(<-errors)
		cancel()
	},
}

const usage = "Welcome to JTunnel .\n\nSource code is at `https://github.com/manoj-inukolunu/jtunnel-go`\n\nTo create a new tunnel\n\nMake a `POST` request to `client://127.0.0.1:1234/create`\nwith the payload\n\n```\n{\n    \"HostName\":\"myhost\",\n    \"TunnelName\":\"Tunnel Name\",\n    \"localServerPort\":\"3131\"\n}\n\n```\n\nThe endpoint you get is `https://myhost.migtunnel.net`\n\nAll the requests to `https://myhost.migtunnel.net` will now\n\nbe routed to your server running on port `3131`\n\n"

type Main struct {
	cmd *cobra.Command
}

func (main *Main) Serve(ctx context.Context) error {
	result := markdown.Render(usage, 80, 6)
	log.Println(string(result))
	localDb := db.NewLocalDb(dbFilePath)
	c := client.NewClient(data.ClientConfig{AdminPort: 1234}, localDb)
	main.cmd.Printf("Starting Admin Server on %d \n", 1234)
	go c.StartAdminServer()
	c.StartControlConnection(localDb)
	return nil
}

func (main *Main) Stop() {
	main.cmd.Println("Stopping Client")
}

func init() {
	rootCmd.AddCommand(startCmd)
	startCmd.Flags().StringVar(&dbFilePath,
		"file", "", "Optional Full File path where db is stored."+
			"If given migtunnel will save requests and responses in sqlite db located at `file`")
}
