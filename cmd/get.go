/*
Copyright © 2023 Patrick Hermann patrick.hermann@sva.de
*/

package cmd

import (
	"fmt"

	http "github.com/go-git/go-git/v5/plumbing/transport/http"
	sthingsBase "github.com/stuttgart-things/sthingsBase"

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

		// GET COMMIT INFORMATION
		commit := internal.GetGitRevision(repoUrl)
		fmt.Println(commit)

		// GET YACHT CONFIG/DEFAULTS
		config := internal.GetYachtConfig("https://github.com/stuttgart-things/yacht-application-server.git", "yas.log", GetGitAuth("", ""))
		fmt.Println(config)

		// RENDER YACHT JSON
		// + READ ENV VARS/WORKSPACES
		vars := map[string]interface{}{"author": "patrick"}
		renderedModuleCall, _ := sthingsBase.RenderTemplateInline(YachtRevisionRunJson, "missingkey=zero", "{{", "}}", vars)
		fmt.Println(string(renderedModuleCall))
		// + OUTPUT TO FILE

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
