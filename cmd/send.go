/*
Copyright © 2023 Patrick Hermann patrick.hermann@sva.de
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

// sendCmd represents the send command
var sendCmd = &cobra.Command{
	Use:   "send",
	Short: "send command",
	Long:  `send command`,
	Run: func(cmd *cobra.Command, args []string) {

		// GET FLAGS
		endpoint, _ := cmd.LocalFlags().GetString("endpoint")
		port, _ := cmd.LocalFlags().GetString("port")
		file, _ := cmd.LocalFlags().GetString("file")

		// SET ADDRESS + SEND TO YAS
		address := endpoint + ":" + port

		if strings.Contains(address, "localhost") {
			internal.ConnectInsecure(address, file)
		} else {
			internal.ConnectSecure(address, file)
		}

	},
}

func init() {
	rootCmd.AddCommand(sendCmd)
	sendCmd.Flags().String("endpoint", "", "grpc fqdn/ip of yas-endpoint")
	sendCmd.Flags().String("port", "443", "grpc port of yas-endpoint")
	sendCmd.Flags().String("file", "", "file containing revisionRuns")
	sendCmd.MarkFlagRequired("endpoint")
	sendCmd.MarkFlagRequired("file")
}
