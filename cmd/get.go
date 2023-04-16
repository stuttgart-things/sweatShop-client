/*
Copyright © 2023 Patrick Hermann patrick.hermann@sva.de
*/

package cmd

import (
	"fmt"
	"path/filepath"
	"strings"

	http "github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/spf13/viper"
	sthingsBase "github.com/stuttgart-things/sthingsBase"

	"github.com/spf13/cobra"
	"github.com/stuttgart-things/yacht-application-client/internal"
)

type PipelineRun struct {
	Params     string `mapstructure:"params"`
	Workspaces string `mapstructure:"workspaces"`
}

type Config struct {
	PipelineRunProfile []map[string]PipelineRun `mapstructure:"pipelineruns"`
}

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "get command",
	Long:  `get command`,
	Run: func(cmd *cobra.Command, args []string) {

		// GET FLAGS
		repoUrl, _ := cmd.LocalFlags().GetString("repo")

		// GET COMMIT INFORMATION
		commit := internal.GetGitRevision(repoUrl)
		fmt.Println("COMMIT", commit)
		fmt.Println(commit["author"])
		fmt.Println(commit["ids"])

		// GET YACHT CONFIG/DEFAULTS
		revisionRunconfig := internal.GetYachtConfig("https://github.com/stuttgart-things/yacht-application-server.git", ".yacht.yaml", GetGitAuth("", ""))
		fmt.Println(revisionRunconfig)

		// RENDER YACHT JSON
		// + READ ENV VARS/WORKSPACES
		// vars := map[string]interface{}{"author": "patrick"}
		renderedModuleCall, _ := sthingsBase.RenderTemplateInline(YachtRevisionRunJson, "missingkey=zero", "{{", "}}", commit)
		fmt.Println(string(renderedModuleCall))
		// + OUTPUT TO FILE

		// READ PIPELINERUNVALUES
		// templatePath := "yacht-values.yaml"
		// var allPipelineRunValues Config

		// allPipelineRunValues = ReadYamlToObject(templatePath, ".yaml", allPipelineRunValues).(Config)

		// for _, pipelineRuns := range allPipelineRunValues.PipelineRunProfile {

		// 	for name, pipelineRun := range pipelineRuns {
		// 		fmt.Println("PIPELINE", name)
		// 		fmt.Println("WORKSPACES", pipelineRun.Workspaces)
		// 		fmt.Println("PARAMS", pipelineRun.Params)
		// 	}
		// }

		localValues, _ := ReadPipelineRunValues("yacht-values.yaml", "build-kaniko-image")
		gitValues, _ := ReadPipelineRunValues("yacht-values.yaml", "build-kaniko-image")

		fmt.Println(localValues, gitValues)
		// templatePath := "yacht-values.yaml"

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

func ReadYamlToObject(pathToConfig, extension string, object interface{}) interface{} {

	// read yaml
	viper.SetConfigName(strings.ReplaceAll(filepath.Base(pathToConfig), extension, ""))
	viper.AddConfigPath(filepath.Dir(pathToConfig))

	if err := viper.ReadInConfig(); err != nil {
		return object
	}

	if err := viper.Unmarshal(&object); err != nil {
		fmt.Println(err)
		return object
	}

	return object

}

func ReadPipelineRunValues(templatePath, pipelineName string) (pipelineRunValues map[string]string, pipelineFound bool) {

	pipelineRunValues = make(map[string]string)
	// templatePath := "yacht-values.yaml"
	var allPipelineRunValues Config

	allPipelineRunValues = ReadYamlToObject(templatePath, ".yaml", allPipelineRunValues).(Config)

	for _, pipelineRuns := range allPipelineRunValues.PipelineRunProfile {

		for name, pipelineRun := range pipelineRuns {
			fmt.Println("PIPELINE", name)

			if name == pipelineName {

				fmt.Println("WORKSPACES", pipelineRun.Workspaces)
				fmt.Println("PARAMS", pipelineRun.Params)
				pipelineRunValues["PARAMS"] = pipelineRun.Params
				pipelineRunValues["WORKSPACES"] = pipelineRun.Workspaces
				pipelineFound = true
			}

		}
	}

	return

}
