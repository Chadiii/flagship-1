/*
Copyright © 2022 Flagship Team flagship@abtasty.com

*/
package flag

import (
	"github.com/flagship-io/flagship/cmd/flag/usage"
	"github.com/spf13/cobra"
)

var (
	FlagID               string
	DataRaw              string
	FlagName             string
	FlagType             string
	FlagDefaultValue     string
	FlagDescription      string
	FlagPredefinedValues string
)

// FlagCmd represents the flag command
var FlagCmd = &cobra.Command{
	Use:   "flag [create|edit|get|list|delete|usage]",
	Short: "Manage your flags",
	Long:  `Manage your flags in your account`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	FlagCmd.AddCommand(usage.UsageCmd)
}
