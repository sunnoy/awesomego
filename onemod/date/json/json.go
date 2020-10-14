/*
 *@Description
 *@author          lirui
 *@create          2020-10-13 11:50
 */
package main

import (
	"encoding/json"
	"fmt"
)

type pipelineRunParam struct {
	// 结构体内的字段必须还要大写才可以进行json转换
	AppId                   string `json:"app_id"`
	BuildScript             string `json:"build_script"`
	DevopsUrl               string `json:"devops_url"`
	ArtifactspPath          string `json:"artifactsp_path"`
	BuildImageEnv           string `json:"build_image_env"`
	MarjorVer               string `json:"marjor_ver"`
	MinorVer                string `json:"minor_ver"`
	ReVer                   string `json:"re_ver"`
	RegistryProject         string `json:"registry_project"`
	DeployEnv               string `json:"deploy_env"`
	BuildProfile            string `json:"build_profile"`
	ImagesRepository        string `json:"images_repository"`
	ImagesRepositoryProject string `json:"images_repository_project"`
}

func main() {
	user := pipelineRunParam{}

	//user := PipelineRunParam{
	//	registryProject: "sssss",
	//}

	b, err := json.Marshal(user)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(b))
}
