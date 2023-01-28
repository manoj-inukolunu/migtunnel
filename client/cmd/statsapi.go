/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"github.com/spf13/cobra"
	"golang/client/admin/http"
)

// statsapiCmd represents the statsapi command
var statsapiCmd = &cobra.Command{
	Use:   "statsapi",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		data, _ := json.Marshal(map[string]interface{}{
			"config":   http.GetClientConfig(),
			"statsUrl": "/stats",
		})
		cmd.Println(string(data))
	},
}

func init() {
	rootCmd.AddCommand(statsapiCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// statsapiCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// statsapiCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
