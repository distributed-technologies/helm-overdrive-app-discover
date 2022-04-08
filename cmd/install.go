/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"path/filepath"
	"strings"

	"github.com/distributed-technologies/helm-overdrive-app-discover/pkg/discover"
	"github.com/distributed-technologies/helm-overdrive-app-discover/pkg/logging"
	"github.com/distributed-technologies/helm-overdrive/pkg/template"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const longDesc = `
The install cmd takes an app_file and the helm-overdrive.yaml config, using these it renders the helm chart from the app.yaml file.
This is all done using the helm-overdrive template pkg.
`

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "The install cmd will take a app.yaml and helm-overdrive config and render the chart",
	Long:  longDesc,
	Run: func(cmd *cobra.Command, args []string) {

		app_file := viper.GetString("app_file")

		files, err := discover.GetFiles(app_file)
		if err != nil {
			panic(err)
		}

		path := files[0]
		file := filepath.Base(path)
		var app discover.App

		app.GetValuesFromYamlFile(path)

		logging.Debug("app: %v\n", app)

		additional_resources_folder := viper.GetString("additional_resources")
		base_folder := viper.GetString("base_folder")
		chart_name := app.Source.Chart_name
		chart_version := app.Source.Chart_version
		env_folder := viper.GetString("env_folder")
		global_file := viper.GetString("global_file")
		helm_repo := app.Source.Helm_repo
		values_file := viper.GetString("values_file")

		app_folder := path[strings.Index(path, base_folder)+(len(base_folder)+1) : len(path)-(len(file)+1)]

		err = template.Template(
			additional_resources_folder,
			app_folder,
			base_folder,
			chart_name,
			chart_version,
			env_folder,
			global_file,
			helm_repo,
			values_file)

		if err != nil {
			panic(err)
		}

	},
}

func init() {
	rootCmd.AddCommand(installCmd)
	cobra.OnInitialize(initInstallCmdConfig)

	installCmd.Flags().String("app_file", "./", "The folder of the app you want to install")
	installCmd.Flags().StringVarP(&cfgFile, "config", "c", "", "Point to the helm-overdrive config (default is $HOME/.helm-overdrive.yaml)")

	viper.BindPFlags(installCmd.Flags())

	initInstallCmdConfig()
}

// initConfig reads in config file and ENV variables if set.
func initInstallCmdConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else if viper.GetString("config") != "" {
		// Use config file from environment 'HO_CONFIG'
		viper.SetConfigFile(viper.GetString("config"))
	} else {
		// Look in these paths for a config file
		viper.AddConfigPath("./") // Checks running dir
		viper.SetConfigType("yaml")
		viper.SetConfigName("helm-overdrive")
	}

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		logging.Debug("Using config file: %s\n", viper.ConfigFileUsed())
	}
}
