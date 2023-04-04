/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"strings"

	"github.com/stuttgart-things/yacht-application-client/internal"

	"github.com/spf13/cobra"
)

//var (
// address = "localhost:50051"
// address = "yas.dev.sthings.tiab.ssc.sva.de:443"
// address = "sleeveless-dev.labul.sva.de:8011"
//)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "create command",
	Long:  `create command`,
	Run: func(cmd *cobra.Command, args []string) {

		// input: get flags
		endpoint, _ := cmd.LocalFlags().GetString("endpoint")
		port, _ := cmd.LocalFlags().GetString("port")

		// set adress
		address := endpoint + ":" + port

		if strings.Contains(address, "localhost") {
			internal.ConnectInsecure(address)
		} else {
			internal.ConnectSecure(address)
		}

	},
}

func init() {
	rootCmd.AddCommand(createCmd)
	createCmd.Flags().String("endpoint", "", "grpc fqdn/ip of yas-endpoint")
	createCmd.Flags().String("port", "443", "grpc port of yas-endpoint")
	createCmd.Flags().String("file", "", "file containing revisionRuns")
	createCmd.MarkFlagRequired("endpoint")
	createCmd.MarkFlagRequired("file")
}
