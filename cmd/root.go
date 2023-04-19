/*
Copyright © 2023 Patrick Hermann patrick.hermann@sva.de
*/

package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

type RevisionRun struct {
	Stage    int    `mapstructure:"stage"`
	Pipeline string `mapstructure:"pipeline"`
	Params   string `mapstructure:"params"`
}

type RevisionRunConfig struct {
	RevisionRunProfile []map[string]RevisionRun `yaml:"revisionrun"`
}

type PipelineRun struct {
	Params     string `mapstructure:"params"`
	Workspaces string `mapstructure:"workspaces"`
}

type PipelineRunConfig struct {
	PipelineRunProfile []map[string]PipelineRun `mapstructure:"pipelineruns"`
}

type PipelineRunJson struct {
	Name       string
	Stage      int
	Params     string
	Workspaces string
}

type RevisionRunJson struct {
	RepoName            string
	PushedAt            string
	Author              string
	RepoUrl             string
	CommitId            string
	RevisionRunCreation string
	PipelineRuns        []PipelineRun
}

const YachtRevisionRunJson = `
{
	"repo_name": "{{ .name }}",
	"pushed_at": "{{ .date }}",
	"author": "{{ .author }}",
	"repo_url": "{{ .url }}",
	"commit_id": "{{ .id }}",
	"pipelineruns": [
	{{ range .halo2 }}
	{{ .Stage }}
		{
		"name": "build-kaniko-image",
		"stage": 0,
		"params": "context=/kaniko/decksman, dockerfile=./Dockerfile, git-revision=main, gitRepoUrl='git@codehub.sva.de:Lab/stuttgart-things/yacht/yacht-application-server.git', gitWorkspaceSubdirectory=/kaniko/decksman, image=scr.tiab.labda.sva.de/yacht-application-server/yacht-application-server, registry=scr.tiab.labda.sva.de, tag=0.0.90",
		"workspaces": "ssh-credentials=secret;codehub-ssh;secretName, shared-workspace=persistentVolumeClaim;sthings-kaniko-workspace;claimName, dockerconfig=secret;scr-labda;secretName"
		}
	{{ end }}
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
