package cmd

import (
  "fmt"

  "github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
  Use:   "version",
  Short: "Print the version of namedns",
  Run: func(cmd *cobra.Command, args []string) {
    fmt.Println("namedns", Version)
  },
}

// Version should be set to actual version string in `go build -ldflags "-X github.com/aelindeman/namedns/cmd.Version="`
var Version string

func init() {
  rootCmd.AddCommand(versionCmd)
}
