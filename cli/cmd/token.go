package cmd

import (
	"fmt"

	"github.com/bryk-io/dlt4eu/service"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var tokenCmd = &cobra.Command{
	Use:   "token",
	Short: "Generate an admin credential to access the API",
	Aliases: []string{"admin", "credential"},
	RunE: func(cmd *cobra.Command, args []string) error {
		// Load service handler settings.
		conf := new(service.Config)
		if err := viper.UnmarshalKey("service", conf); err != nil {
			return err
		}

		// Get service handler instance.
		handler, err := service.New(conf)
		if err != nil {
			return err
		}

		// Get credentials
		token, err := handler.AdminToken(args[0])
		if err != nil {
			return err
		}
		fmt.Println(token)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(tokenCmd)
}
