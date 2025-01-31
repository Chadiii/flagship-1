package mockfunction

import (
	"net/http"

	"github.com/flagship-io/flagship/models"
	"github.com/flagship-io/flagship/utils"
	"github.com/flagship-io/flagship/utils/config"
	"github.com/jarcoal/httpmock"
	"github.com/spf13/viper"
)

var TestFlag = models.Flag{
	Id:          "testFlagID",
	Name:        "testFlagName",
	Type:        "string",
	Description: "testFlagDescription",
	Source:      "cli",
}

var TestFlag1 = models.Flag{
	Id:          "testFlagID1",
	Name:        "testFlagName1",
	Type:        "string",
	Description: "testFlagDescription1",
	Source:      "cli",
}

var TestFlagEdit = models.Flag{
	Id:          "testFlagID",
	Name:        "testFlagName1",
	Type:        "string",
	Description: "testFlagDescription1",
	Source:      "cli",
}

var TestFlagList = []models.Flag{
	TestFlag,
	TestFlag1,
}

var TestFlagUsageList = []models.FlagUsage{
	{
		Id:                "testFlagUsageID",
		FlagKey:           "isVIP",
		Repository:        "flagship",
		FilePath:          "https://github.com/flagship-io/flagship",
		Branch:            "main",
		Line:              "Line116",
		CodeLineHighlight: "codeLineHighlight",
		Code:              "code",
	},
}

func APIFlag() {
	config.SetViper()

	resp := utils.HTTPListResponse[models.Flag]{
		Items:             TestFlagList,
		CurrentItemsCount: 2,
		CurrentPage:       1,
		TotalCount:        2,
		ItemsPerPage:      10,
		LastPage:          1,
	}

	respUsage := utils.HTTPListResponse[models.FlagUsage]{
		Items:             TestFlagUsageList,
		CurrentItemsCount: 2,
		CurrentPage:       1,
		TotalCount:        1,
		ItemsPerPage:      10,
		LastPage:          1,
	}

	httpmock.RegisterResponder("GET", utils.GetHost()+"/v1/accounts/"+viper.GetString("account_id")+"/flags/"+TestFlag.Id,
		func(req *http.Request) (*http.Response, error) {
			resp, _ := httpmock.NewJsonResponse(200, TestFlag)
			return resp, nil
		},
	)

	httpmock.RegisterResponder("GET", utils.GetHost()+"/v1/accounts/"+viper.GetString("account_id")+"/flags",
		func(req *http.Request) (*http.Response, error) {
			resp, _ := httpmock.NewJsonResponse(200, resp)
			return resp, nil
		},
	)

	httpmock.RegisterResponder("GET", utils.GetHost()+"/v1/accounts/"+viper.GetString("account_id")+"/account_environments/"+viper.GetString("account_environment_id")+"/flags_usage",
		func(req *http.Request) (*http.Response, error) {
			resp, _ := httpmock.NewJsonResponse(200, respUsage)
			return resp, nil
		},
	)

	httpmock.RegisterResponder("POST", utils.GetHost()+"/v1/accounts/"+viper.GetString("account_id")+"/flags",
		func(req *http.Request) (*http.Response, error) {
			resp, _ := httpmock.NewJsonResponse(200, TestFlag)
			return resp, nil
		},
	)

	httpmock.RegisterResponder("PATCH", utils.GetHost()+"/v1/accounts/"+viper.GetString("account_id")+"/flags/"+TestFlag.Id,
		func(req *http.Request) (*http.Response, error) {
			resp, _ := httpmock.NewJsonResponse(200, TestFlagEdit)
			return resp, nil
		},
	)

	httpmock.RegisterResponder("DELETE", utils.GetHost()+"/v1/accounts/"+viper.GetString("account_id")+"/flags/"+TestFlag.Id,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(204, ""), nil
		},
	)
}
