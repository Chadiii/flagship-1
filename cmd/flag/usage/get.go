/*
Copyright © 2022 Flagship Team flagship@abtasty.com

*/
package usage

import (
	"log"

	"github.com/flagship-io/flagship/utils"
	httprequest "github.com/flagship-io/flagship/utils/httpRequest"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// getCmd represents get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "List all flag usage statistics inside your codebase",
	Long:  `List all flag usage statistics inside your codebase in your account`,
	Run: func(cmd *cobra.Command, args []string) {
		body, err := httprequest.HTTPFlagUsage()
		if err != nil {
			log.Fatalf("error occurred: %v", err)
		}
		utils.FormatItem([]string{"Id", "FlagKey", "Repository", "FilePath", "Branch", "Line", "CodeLineHighlight", "Code"}, body, viper.GetString("output_format"), cmd.OutOrStdout())
	},
}

func init() {
	UsageCmd.AddCommand(getCmd)
}
