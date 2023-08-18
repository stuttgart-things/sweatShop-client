/*
Copyright © 2023 Patrick Hermann patrick.hermann@sva.de
*/

package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"path/filepath"
	"strings"

	sthingsBase "github.com/stuttgart-things/sthingsBase"
	"gopkg.in/yaml.v2"

	http "github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stuttgart-things/sweatShop-client/internal"
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
		// fmt.Println("allRevisionRuns:", allRevisionRuns)

		// READ PIPELINERUNVALUES FROM LOCAL FILE
		var allPipelineRuns PipelineRunConfig
		allPipelineRuns = ReadYamlToObject(local, ".yaml", allPipelineRuns).(PipelineRunConfig)

		// ITERATE OVER STAGES

		var prs1 []PipelineRunJson

		for _, revisionRun := range allPipelineRuns.PipelineRunProfile {

			for name := range revisionRun {

				stage, params := ValidatePipelineRunValues(name, allRevisionRuns)

				renderedParameters, _ := sthingsBase.RenderTemplateInline(params+", "+revisionRun[name].Params, "missingkey=zero", "{{", "}}", commit)

				pr := PipelineRunJson{name, stage, string(renderedParameters), revisionRun[name].Workspaces}

				prs1 = append(prs1, pr)

			}
		}

		// fmt.Println(prs1)

		// localValues, _ := ReadPipelineRunValues(local, "build-kaniko-image")
		// gitValues, _ := ReadPipelineRunValues(remote, "build-kaniko-image")

		// fmt.Println(localValues, gitValues)
		// templatePath := "yacht-values.yaml"

		// RENDER YACHT JSON
		// + READ ENV VARS/WORKSPACES
		// vars := map[string]interface{}{"author": "patrick"}
		// vars := map[string]interface{}{"author": "patrick"}

		// renderedModuleCall, _ := sthingsBase.RenderTemplateInline(YachtRevisionRunJson, "missingkey=zero", "{{", "}}", commit)
		// fmt.Println(string(renderedModuleCall))
		// + OUTPUT TO FILE

		// hello := yas.Workspace{"", "", "", ""}
		// bla := []yas.Workspace{hello}
		// bla = append(bla, hello)
		// pipelineParams := make(map[string]string)
		// pipelineParams["hello"] = "hello"
		// hello2 := yas.PipelineRun{commit["name"].(string), commit["author"].(string), commit["name"].(string), commit["url"].(string), commit["id"].(string), commit["date"].(string), "", "", "", "", pipelineParams, bla, "", "", ""}

		hello3 := RevisionRunJson{commit["name"].(string), commit["date"].(string), commit["author"].(string), commit["url"].(string), commit["id"].(string), commit["date"].(string), prs1}
		k, _ := json.MarshalIndent(hello3, "", "  ")
		log.Println(string(k))

		// fmt.Println(hello)
		// fmt.Println(hello2)

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
	var allPipelineRuns PipelineRunConfig

	allPipelineRuns = ReadYamlToObject(templatePath, ".yaml", allPipelineRuns).(PipelineRunConfig)

	for _, pipelineRuns := range allPipelineRuns.PipelineRunProfile {

		for name, pipelineRun := range pipelineRuns {
			fmt.Println("PIPELINE", name)

			if strings.Contains(name, pipelineName) {

				pipelineRunValues["PARAMS"] = pipelineRun.Params
				pipelineRunValues["WORKSPACES"] = pipelineRun.Workspaces
				pipelineFound = true

			}
		}
	}

	return

}

func ValidatePipelineRunValues(pipelineName string, allRevisionRuns RevisionRunConfig) (stage int, params string) {

	for _, revisionRun := range allRevisionRuns.RevisionRunProfile {

		for _, revisionRun := range revisionRun {

			if strings.Contains(pipelineName, revisionRun.Pipeline) {
				stage = revisionRun.Stage
				params = revisionRun.Params

			}

		}
	}

	return
}
