/*
Copyright © 2023 Patrick Hermann patrick.hermann@sva.de
*/

package cmd

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"

	http "github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/spf13/viper"
	sthingsBase "github.com/stuttgart-things/sthingsBase"

	"github.com/spf13/cobra"
	"github.com/stuttgart-things/yacht-application-client/internal"
)

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "get command",
	Long:  `get command`,
	Run: func(cmd *cobra.Command, args []string) {

		// GET FLAGS
		repoUrl, _ := cmd.LocalFlags().GetString("repo")
		local, _ := cmd.LocalFlags().GetString("local")
		remote, _ := cmd.LocalFlags().GetString("remote")

		// GET COMMIT INFORMATION
		commit := internal.GetGitRevision(repoUrl)

		// GET YACHT CONFIG/DEFAULTS FROM GIT
		revisionRunconfig := internal.GetYachtConfig(repoUrl, remote, GetGitAuth("", ""))
		allRevisionRuns := RevisionRunConfig{}
		if err := yaml.Unmarshal([]byte(revisionRunconfig), &allRevisionRuns); err != nil {
			log.Fatal(err)
		}
		fmt.Println("allRevisionRuns:", allRevisionRuns)

		// READ PIPELINERUNVALUES FROM LOCAL FILE
		var allPipelineRuns PipelineRunConfig
		allPipelineRuns = ReadYamlToObject(local, ".yaml", allPipelineRuns).(PipelineRunConfig)

		fmt.Println(allPipelineRuns)
		// for _, pipelineRuns := range allPipelineRuns.PipelineRunProfile {
		// 	for name, pipelineRun := range pipelineRuns {
		// 		fmt.Println("PIPELINE", name)
		// 		fmt.Println("WORKSPACES", pipelineRun.Workspaces)
		// 		fmt.Println("PARAMS", pipelineRun.Params)
		// 	}
		// }

		// ITERATE OVER STAGES
		for _, revisionRun := range allRevisionRuns.RevisionRunProfile {
			for name, revisionRun := range revisionRun {

				fmt.Println("NAME", name)
				fmt.Println("STAGE", revisionRun.Stage)
				ValidateGetLocalValues(revisionRun.Pipeline, allPipelineRuns)
			}
		}

		localValues, _ := ReadPipelineRunValues(local, "build-kaniko-image")
		gitValues, _ := ReadPipelineRunValues(remote, "build-kaniko-image")

		fmt.Println(localValues, gitValues)
		// templatePath := "yacht-values.yaml"

		// RENDER YACHT JSON
		// + READ ENV VARS/WORKSPACES
		// vars := map[string]interface{}{"author": "patrick"}
		renderedModuleCall, _ := sthingsBase.RenderTemplateInline(YachtRevisionRunJson, "missingkey=zero", "{{", "}}", commit)
		fmt.Println(string(renderedModuleCall))
		// + OUTPUT TO FILE

	},
}

func init() {
	rootCmd.AddCommand(getCmd)
	getCmd.Flags().String("repo", "", "git repository url")
	getCmd.Flags().String("remote", ".yacht.yaml", "remote/repo values")
	getCmd.Flags().String("local", "yacht-values.yaml", "local/default values")
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
	var allPipelineRuns PipelineRunConfig

	allPipelineRuns = ReadYamlToObject(templatePath, ".yaml", allPipelineRuns).(PipelineRunConfig)

	for _, pipelineRuns := range allPipelineRuns.PipelineRunProfile {

		for name, pipelineRun := range pipelineRuns {
			fmt.Println("PIPELINE", name)

			if strings.Contains(name, pipelineName) {

				fmt.Println("WORKSPACES", pipelineRun.Workspaces)
				fmt.Println("PARAMS", pipelineRun.Params)
				// pipelineRunValues["PARAMS"] = pipelineRun.Params
				// pipelineRunValues["WORKSPACES"] = pipelineRun.Workspaces
				pipelineFound = true
			}

		}
	}

	return

}

func ValidateGetLocalValues(pipelineName string, allPipelineRuns PipelineRunConfig) (params, workspaces string) {

	for _, pipelineRuns := range allPipelineRuns.PipelineRunProfile {
		for name, pipelineRun := range pipelineRuns {
			fmt.Println("1", name)
			fmt.Println("2", pipelineName)

			if strings.Contains(name, pipelineName) {
				fmt.Println("FOUND!")
				fmt.Println("WORKSPACES", pipelineRun.Workspaces)
				fmt.Println("PARAMS", pipelineRun.Params)
				params = pipelineRun.Params
				workspaces = pipelineRun.Workspaces
			}
		}
	}

	return
}
