/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package authorization

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func logout() string {
	return "logout from session by deleting token"
}

// logoutCmd represents the logout command
var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "this authorization logout",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(logout())
		configMap := viper.AllSettings()
		delete(configMap, "token")
		encodedConfig, _ := json.MarshalIndent(configMap, "", " ")
		err := viper.ReadConfig(bytes.NewReader(encodedConfig))
		if err != nil {
			fmt.Println(err)
		}
		viper.WriteConfigAs("config.yaml")
	},
}

func init() {

	// Here you will define your flags and configuration settings.
	AuthorizationCmd.AddCommand(logoutCmd)
	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
