package cmd

import (
	"os"

	"github.com/distributed-technologies/helm-overdrive-app-discover/pkg/logging"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var isDebug bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "helm-overdrive",
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
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is $HOME/.helm-overdrive.yaml)")
	rootCmd.PersistentFlags().BoolVar(&isDebug, "debug", false, "enable debug logs")

	viper.BindPFlags(rootCmd.Flags())
	viper.BindPFlags(rootCmd.PersistentFlags())
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.SetEnvPrefix("HO") // Standing for 'helm-overdrive'
	viper.AutomaticEnv()     // read in environment variables that match

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
