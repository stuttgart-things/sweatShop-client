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
	PipelineRuns        []PipelineRunJson
}

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
