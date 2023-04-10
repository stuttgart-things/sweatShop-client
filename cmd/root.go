/*
Copyright © 2023 Patrick Hermann patrick.hermann@sva.de
*/

package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

const YachtRevisionRunJson = `
{
	"repo_name": "stuttgart-things",
	"pushed_at": "2023-02-200T22:40:36Z",
	"author": "{{ .author }}",
	"repo_url": "https://codehub.sva.de/Lab/stuttgart-things/stuttgart-things.git",
	"commit_id": "03eee7226417",
	"pipelineruns": [
		{
		"name": "build-kaniko-image",
		"stage": 0,
		"params": "context=/kaniko/decksman, dockerfile=./Dockerfile, git-revision=main, gitRepoUrl='git@codehub.sva.de:Lab/stuttgart-things/yacht/yacht-application-server.git', gitWorkspaceSubdirectory=/kaniko/decksman, image=scr.tiab.labda.sva.de/yacht-application-server/yacht-application-server, registry=scr.tiab.labda.sva.de, tag=0.0.90",
		"workspaces": "ssh-credentials=secret;codehub-ssh;secretName, shared-workspace=persistentVolumeClaim;sthings-kaniko-workspace;claimName, dockerconfig=secret;scr-labda;secretName"
		}
	]
}
`

var rootCmd = &cobra.Command{
	Use:   "yacht-application-client",
	Short: "yacht-application-client",
	Long:  `yacht-application-client.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
