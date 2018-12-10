package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/aelindeman/goname"
	"github.com/kyoh86/xdg"
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
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if viper.GetBool("verbose") {
			log.SetLevel(log.DebugLevel)
		}
		ValidateGlobalConfig()
	},
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
	client.BaseURL = viper.GetString("api-url")

	loginErr := client.Login()
	if loginErr != nil {
		log.WithError(loginErr).Fatal("could not authenticate")
	}

	apiClient = client
	return client
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", `Path to config file (default "$XDG_CONFIG_HOME/namedns/.namedns.yaml")`)
	rootCmd.PersistentFlags().StringP("username", "u", "", "API username")
	rootCmd.PersistentFlags().StringP("api-key", "k", "", "API key")
	rootCmd.PersistentFlags().StringP("api-url", "", goname.NameAPIBaseURL, "API base URL")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Display debugging output")
	viper.SetEnvPrefix("namedns")
	viper.BindPFlag("username", rootCmd.PersistentFlags().Lookup("username"))
	viper.BindPFlag("api-key", rootCmd.PersistentFlags().Lookup("api-key"))
	viper.BindPFlag("api-url", rootCmd.PersistentFlags().Lookup("api-url"))
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find XDG config directories.
		for _, xdgConfigDir := range xdg.AllConfigDirs() {
			viper.AddConfigPath(strings.Join([]string{xdgConfigDir, "namedns"}, "/"))
		}

		// Find home directory.
		home, homeErr := homedir.Dir()
		if homeErr != nil {
			log.WithError(homeErr).Warning("could not find home directory")
		} else {
			viper.AddConfigPath(home)
		}

		// Look for one in the current directory.
		wd, wdErr := os.Getwd()
		if wdErr != nil {
			log.WithError(wdErr).Warning("could not find current directory")
		} else {
			viper.AddConfigPath(wd)
		}

		// Search config with name ".namedns" (without extension).
		viper.SetConfigName(".namedns")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if configReadErr := viper.ReadInConfig(); configReadErr == nil {
		log.WithField("file", viper.ConfigFileUsed()).Debug("using config from file")
	}
}
