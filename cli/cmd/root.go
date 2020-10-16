package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Configuration file provided
var cfgFile string

// Used as the main binary identifier, ENV variables prefix and home directories
var appIdentifier = "dlt4eu"

var rootCmd = &cobra.Command{
	Use:   appIdentifier,
	Short: "Digital identity platform for the DLT4EU project",
}

// Execute provides the main entry point for the application
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("error: %s\n", err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file")
}

func initConfig() {
	// ENV
	viper.SetEnvPrefix(appIdentifier)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// Configuration file
	viper.AddConfigPath(fmt.Sprintf("/etc/%s", appIdentifier))
	viper.AddConfigPath(fmt.Sprintf("$HOME/%s", appIdentifier))
	viper.AddConfigPath(fmt.Sprintf("$HOME/.%s", appIdentifier))
	viper.AddConfigPath(".")
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	}

	// Read configuration file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			fmt.Printf("failed to read configuration file: %s\n", err.Error())
		}
	}
}
