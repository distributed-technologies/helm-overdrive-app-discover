/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/distributed-technologies/helm-overdrive-app-discover/pkg/discover"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const discoverDesc = `
The discover cmd builds ArgoCD application resources based on any file, it finds in the current folder and subfolders,
that contains 'apiVersion: argocd-discover/v1alpha1' as the first line of the file.

This is done by walking through the folder structure reading the first line of any '*.yaml' file,
and checking if it matches 'apiVersion: argocd-discover/v1alpha1' if it does, it reads the rest of the content,
which describes the chart that should be used and where it should be deployed.
`

// discoverCmd represents the discover command
var discoverCmd = &cobra.Command{
	Use:   "discover",
	Short: "Discovers any files, in the current folder or subfolder, that contains 'apiVersion: argocd-discover/v1alpha1'.",
	Long:  discoverDesc,
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
