/*
Copyright © 2023 Patrick Hermann patrick.hermann@sva.de
*/

package cmd

import (
	"fmt"
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

		// GET FLAGS
		endpoint, _ := cmd.LocalFlags().GetString("endpoint")
		port, _ := cmd.LocalFlags().GetString("port")
		file, _ := cmd.LocalFlags().GetString("file")

		// SET ADDRESS
		address := endpoint + ":" + port

		if strings.Contains(address, "localhost") {
			internal.ConnectInsecure(address, file)
		} else {
			internal.ConnectSecure(address, file)
		}

		commit := internal.GetGitRevision("https://github.com/stuttgart-things/yacht-application-server.git")
		fmt.Println(commit)

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
