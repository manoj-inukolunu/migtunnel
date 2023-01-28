package cmd

import (
	"context"
	markdown "github.com/MichaelMure/go-term-markdown"
	"github.com/spf13/cobra"
	"github.com/thejerf/suture/v4"
	"golang/jtunnel-client/admin/client"
	"golang/jtunnel-client/admin/data"
	"log"
)

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
	c := client.Client{ClientConfig: data.ClientConfig{AdminPort: 1234}}
	main.cmd.Printf("Starting Admin Server on %d \n", 1234)
	go c.StartAdminServer()
	c.StartControlConnection()
	return nil
}

func (main *Main) Stop() {
	main.cmd.Println("Stopping Client")
}

func init() {
	rootCmd.AddCommand(startCmd)
}
