/*
Copyright © 2022 Flagship Team flagship@abtasty.com

*/

package resource

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/d5/tengo/v2"
	"github.com/flagship-io/flagship/models"
	"github.com/flagship-io/flagship/utils"
	httprequest "github.com/flagship-io/flagship/utils/httpRequest"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Data interface {
	Save(data string) ([]byte, error)
}

type ProjectData struct {
	*models.Project
}

type ResourceData struct {
	Id string `json:"id"`
}

func (f ProjectData) Save(data string) ([]byte, error) {
	return httprequest.HTTPCreateProject(data)
}

type CampaignData struct {
	Id              string               `json:"id,omitempty"`
	ProjectId       string               `json:"project_id"`
	Name            string               `json:"name"`
	Description     string               `json:"description"`
	Type            string               `json:"type"`
	VariationGroups []VariationGroupData `json:"variation_groups"`
}

func (f CampaignData) Save(data string) ([]byte, error) {
	return httprequest.HTTPCreateCampaign(data)
}

type FlagData struct {
	*models.Flag
}

func (f FlagData) Save(data string) ([]byte, error) {
	return httprequest.HTTPCreateFlag(data)
}

type GoalData struct {
	*models.Goal
}

func (f GoalData) Save(data string) ([]byte, error) {
	return httprequest.HTTPCreateGoal(data)
}

type TargetingKeysData struct {
	*models.TargetingKey
}

func (f TargetingKeysData) Save(data string) ([]byte, error) {
	return httprequest.HTTPCreateTargetingKey(data)
}

type VariationGroupData struct {
	*models.VariationGroup
}

/* func (f VariationGroupData) Save(data string) ([]byte, error) {
	return httprequest.HTTPCreateVariationGroup(campaignID, data)
} */

type VariationData struct {
	*models.Variation
}

/* func (f VariationData) Save(data string) ([]byte, error) {
	return httprequest.HTTPCreateVariation(campaignID, variationGroupID, data)
} */

// define structs for other resource types

type ResourceType int

const (
	Project ResourceType = iota
	Flag
	TargetingKey
	Goal
	Campaign
	VariationGroup
	Variation
)

var resourceTypeMap = map[string]ResourceType{
	"project":         Project,
	"flag":            Flag,
	"targeting_key":   TargetingKey,
	"goal":            Goal,
	"campaign":        Campaign,
	"variation_group": VariationGroup,
	"variation":       Variation,
}

type Resource struct {
	Name             ResourceType
	Data             Data
	ResourceVariable string
}

func UnmarshalConfig(filePath string) ([]Resource, error) {
	var config struct {
		Resources []struct {
			Name             string
			Data             json.RawMessage
			ResourceVariable string
		}
	}

	bytes, err := os.ReadFile(resourceFile)

	if err != nil {
		log.Fatalf("error occurred: %v", err)
	}

	if err := json.Unmarshal(bytes, &config); err != nil {
		return nil, err
	}

	var resources []Resource
	for _, r := range config.Resources {
		name, ok := resourceTypeMap[r.Name]
		if !ok {
			return nil, fmt.Errorf("invalid resource name: %s", r.Name)
		}

		var data Data = nil
		var err error = nil

		switch name {

		case Project:
			projectData := ProjectData{}
			err = json.Unmarshal(r.Data, &projectData)
			data = projectData
			//fmt.Println(data)

		//data = &ProjectData{}
		case Flag:
			flagData := FlagData{}
			err = json.Unmarshal(r.Data, &flagData)
			data = flagData
			//fmt.Println(data)

		case TargetingKey:
			targetingKeyData := TargetingKeysData{}
			err = json.Unmarshal(r.Data, &targetingKeyData)
			data = targetingKeyData
			//fmt.Println(data)

		case Campaign:
			campaignData := CampaignData{}
			err = json.Unmarshal(r.Data, &campaignData)
			data = campaignData

		case Goal:
			goalData := GoalData{}
			err = json.Unmarshal(r.Data, &goalData)
			data = goalData
			//fmt.Println(data)

			/* 		case VariationGroup:
			variationGroupData := VariationGroupData{}
			err = json.Unmarshal(r.Data, &variationGroupData)
			data = variationGroupData
			fmt.Println(data) */

		}

		if err != nil {
			return nil, err
		}

		resources = append(resources, Resource{Name: name, Data: data, ResourceVariable: r.ResourceVariable})
	}

	//flag := resources[1].Data.(ProjectData).Name
	//fmt.Println(flag)
	return resources, nil
}

func loadResources(resources []Resource) (string, error) {

	for _, resource := range resources {
		var url = ""
		var resp []byte
		data, err := json.Marshal(resource.Data)
		if err != nil {
			return "", err
		}

		switch resource.Name {
		case Project:
			url = "/projects"
		case Flag:
			url = "/flags"
		case TargetingKey:
			url = "/targeting_keys"
		case Goal:
			url = "/goals"
		case VariationGroup:
			url = "/variable_groups"
		case Variation:
			url = "/variations"
		case Campaign:
			url = "/campaigns"
		}

		if resource.Name == Project || resource.Name == TargetingKey || resource.Name == Flag {
			resp, err = httprequest.HTTPRequest(http.MethodPost, utils.GetHost()+"/v1/accounts/"+viper.GetString("account_id")+url, data)
		}

		if resource.Name == Goal || resource.Name == Campaign {
			resp, err = httprequest.HTTPRequest(http.MethodPost, utils.GetHost()+"/v1/accounts/"+viper.GetString("account_id")+"/account_environments/"+viper.GetString("account_environment_id")+url, data)
		}

		if err != nil {
			return "", err
		}

		log.Println(string(resp))

	}
	return "done", nil
}

var gResources []Resource

// LoadCmd represents the load command
var loadCmd = &cobra.Command{
	Use:   "load [--file=<file>]",
	Short: "Load your resources",
	Long:  `Load your resources`,
	Run: func(cmd *cobra.Command, args []string) {
		/* 		res, err := loadResources(gResources)
		   		if err != nil {
		   			log.Fatalf("error occurred: %v", err)
		   		}
		   		fmt.Fprintf(cmd.OutOrStdout(), "%s\n", res) */
		ScriptResource(gResources)
	},
}

func init() {
	cobra.OnInitialize(initResource)

	loadCmd.Flags().StringVarP(&resourceFile, "file", "", "", "resource file that contains your resource")

	if err := loadCmd.MarkFlagRequired("file"); err != nil {
		log.Fatalf("error occurred: %v", err)
	}

	ResourceCmd.AddCommand(loadCmd)
}

func initResource() {

	// Use config file from the flag.
	var err error
	if resourceFile != "" {
		gResources, err = UnmarshalConfig(resourceFile)
		if err != nil {
			log.Fatalf("error occurred: %v", err)
		}
	}
}

func ScriptResource(resources []Resource) {

	//var resourceVariables map[string]interface{}
	resourceVariables := make(map[string]interface{})

	for _, resource := range resources {
		var resourceData map[string]interface{}
		var response []byte
		var responseData interface{}
		var url = ""

		data, err := json.Marshal(resource.Data)
		if err != nil {
			fmt.Printf("error occurred marshal data: %v\n", err)
		}

		switch resource.Name {
		case Project:
			url = "/projects"
		case Flag:
			url = "/flags"
		case TargetingKey:
			url = "/targeting_keys"
		case Goal:
			url = "/goals"
		case VariationGroup:
			url = "/variable_groups"
		case Variation:
			url = "/variations"
		case Campaign:
			url = "/campaigns"
		}

		err = json.Unmarshal(data, &resourceData)

		if err != nil {
			fmt.Printf("error occurred unmarshall resourceData: %v\n", err)
		}

		for k, vInterface := range resourceData {
			v, ok := vInterface.(string)
			if ok {
				if strings.Contains(v, "$") {
					vTrim := strings.Trim(v, "$")
					for k_, variable := range resourceVariables {
						script, err := tengo.Eval(context.Background(), vTrim, map[string]interface{}{
							k_: variable,
						})

						if err != nil {
							log.Fatalf("error compiled: %v", err)
						}
						resourceData[k] = script.(string)
					}
				}

			}

		}

		dataResource, err := json.Marshal(resourceData)
		if err != nil {
			log.Fatalf("error occurred http call: %v\n", err)
		}

		if resource.Name == Project || resource.Name == TargetingKey || resource.Name == Flag {
			response, err = httprequest.HTTPRequest(http.MethodPost, utils.GetHost()+"/v1/accounts/"+viper.GetString("account_id")+url, dataResource)
		}

		if resource.Name == Goal || resource.Name == Campaign {
			response, err = httprequest.HTTPRequest(http.MethodPost, utils.GetHost()+"/v1/accounts/"+viper.GetString("account_id")+"/account_environments/"+viper.GetString("account_environment_id")+url, dataResource)
		}

		if err != nil {
			log.Fatalf("error occurred http call: %v\n", err)
		}

		fmt.Println(string(response))

		err = json.Unmarshal(response, &responseData)

		if err != nil {
			fmt.Printf("error occurred unmarshal responseData: %v\n", err)
		}

		if responseData == nil {
			fmt.Println("error occurred not response data: " + string(response))
			continue
		}

		resourceVariables[resource.ResourceVariable] = responseData
	}
}
