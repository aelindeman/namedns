package cmd

import (
  log "github.com/sirupsen/logrus"
  "github.com/spf13/cobra"
)

// deleteCmd represents the unset command
var deleteCmd = &cobra.Command{
  Use:     "delete <domain> <record id> [record id ...]",
  Short:   "Delete DNS records from a domain",
  Args:    cobra.MinimumNArgs(2),
  Aliases: []string{"remove", "rm"},
  Run:     unsetDNSRecord,
}

func init() {
  rootCmd.AddCommand(deleteCmd)
  ValidateGlobalConfig()
}

func unsetDNSRecord(cmd *cobra.Command, args []string) {
  for _, id := range args[1:] {
    _, err := GetClient().DeleteDNSRecord(args[0], id)
    if err != nil {
      log.WithError(err).WithFields(log.Fields{
        "domain":   args[0],
        "recordID": id,
      }).Error("could not unset DNS record")
      continue
    }
    log.WithField("recordID", id).Info("record deleted successfully")
  }
}
