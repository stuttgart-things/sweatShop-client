/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/stuttgart-things/yacht-application-client/internal"

	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "create command",
	Long:  `create command`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("create called")
	},
}

var (
	//address = "localhost:50051"
	address = "yas.dev.sthings.tiab.ssc.sva.de:443"
	//address = "sleeveless-dev.labul.sva.de:8011"
)

func init() {
	rootCmd.AddCommand(createCmd)

	// Check env vor given server port
	if os.Getenv("YAS_URL") != "" {
		address = os.Getenv("YAS_URL")
	}

	if strings.Contains(address, "localhost") {
		internal.ConnectInsecure(address)
	} else {
		internal.ConnectSecure(address)
	}

}
