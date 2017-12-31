package cmd

import (
	"fmt"
	"os"

	"github.com/aelindeman/goname"
	homedir "github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile   string
	apiClient *goname.GoName
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "namedns",
	Short: "A command-line utility to manipulate Name.com DNS records",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// ValidateGlobalConfig ensures required flags are set.
func ValidateGlobalConfig() error {
	if viper.GetString("username") == "" {
		return fmt.Errorf("username is not set")
	}
	if viper.GetString("api-key") == "" {
		return fmt.Errorf("API key is not set")
	}
	return nil
}

// GetClient sets up and returns the Name.com API client.
func GetClient() *goname.GoName {
	if apiClient != nil {
		return apiClient
	}

	configErr := ValidateGlobalConfig()
	if configErr != nil {
		log.WithError(configErr).Fatal("configuration error")
	}

	client := goname.New(viper.GetString("username"), viper.GetString("api-key"))

	loginErr := client.Login()
	if loginErr != nil {
		log.WithError(loginErr).Fatal("could not authenticate")
	}

	apiClient = client
	return client
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "path to config file (default: $HOME/.namedns.yaml)")
	rootCmd.PersistentFlags().StringP("username", "u", "", "API username")
	rootCmd.PersistentFlags().StringP("api-key", "k", "", "API key")
	viper.BindPFlag("username", rootCmd.PersistentFlags().Lookup("username"))
	viper.BindPFlag("api-key", rootCmd.PersistentFlags().Lookup("api-key"))
	viper.BindPFlag("staging", rootCmd.PersistentFlags().Lookup("staging"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			log.WithError(err).Fatal("could not find home directory")
		}

		// Search config in home directory with name ".namedns" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".namedns")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.WithField("file", viper.ConfigFileUsed()).Debug("using config from file")
	}
}
