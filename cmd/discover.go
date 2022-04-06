/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/distributed-technologies/helm-overdrive-app-discover/pkg/discover"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// discoverCmd represents the discover command
var discoverCmd = &cobra.Command{
	Use:   "discover",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		folder := viper.GetString("folder")

		err := discover.Discover(folder)
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(discoverCmd)

	discoverCmd.Flags().String("folder", "./", "Folder to find apps in (recursive)")

	viper.BindPFlags(discoverCmd.Flags())
}
