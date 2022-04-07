package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var isDebug bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "helm-overdrive-app-discover",
	Short: "Templating multiple environments together",
	Long: `Helm-overdrive is a tool that allows the templating og multiple yaml resources on top of each other.	`,

	// Run: func(cmd *cobra.Command, args []string) {
	// },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Persistent Flags will be available to this command and all subcommands to this
	// rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is $HOME/.helm-overdrive.yaml)")
	rootCmd.PersistentFlags().BoolVar(&isDebug, "debug", false, "enable debug logs")

	viper.BindPFlags(rootCmd.Flags())
	viper.BindPFlags(rootCmd.PersistentFlags())
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.SetEnvPrefix("HO") // Standing for 'helm-overdrive'
	viper.AutomaticEnv()     // read in environment variables that match

}
