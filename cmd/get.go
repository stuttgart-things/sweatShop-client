/*
Copyright © 2023 Patrick Hermann patrick.hermann@sva.de
*/

package cmd

import (
	"fmt"

	http "github.com/go-git/go-git/v5/plumbing/transport/http"

	"github.com/spf13/cobra"
	"github.com/stuttgart-things/yacht-application-client/internal"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "get command",
	Long:  `get command`,
	Run: func(cmd *cobra.Command, args []string) {

		// GET FLAGS
		repoUrl, _ := cmd.LocalFlags().GetString("repo")

		commit := internal.GetGitRevision(repoUrl)
		fmt.Println(commit)

		internal.CloneRepository("https://github.com/stuttgart-things/yacht-application-server.git", GetGitAuth("", ""))

	},
}

func init() {
	rootCmd.AddCommand(getCmd)
	getCmd.Flags().String("repo", "", "git repository url")
	getCmd.MarkFlagRequired("repo")
}

func GetGitAuth(gitUser, gitToken string) *http.BasicAuth {
	return &http.BasicAuth{
		Username: gitUser,
		Password: gitToken,
	}
}
