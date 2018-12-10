package cmd

import (
	"fmt"
	"strings"

	"github.com/aelindeman/goname"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:     "list [domain ...]",
	Short:   "List DNS records for a domain",
	Aliases: []string{"ls"},
	Run:     listDNSRecords,
}

var listOutputFormat string

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.PersistentFlags().StringVarP(&listOutputFormat, "format", "f", "", "output format")
}

func listDNSRecords(cmd *cobra.Command, args []string) {
	domains := parseRequestedDomains(args)
	log.WithField("domains", domains).Debug("fetching DNS records for domains")

	for _, domain := range domains {
		records, recordsErr := GetClient().ListDNSRecords(domain)
		if recordsErr != nil || records.Result.Failed() {
			log.WithError(recordsErr).WithFields(log.Fields{
				"domain": domain,
			}).Error("error fetching DNS records")
			continue
		}

		switch listOutputFormat {
		case "json":
			// printJsonRecordsList(domain, records.Records)
			fmt.Println("not implemented")
		case "yaml":
			// printYamlRecordsList(domain, records.Records)
			fmt.Println("not implemented")
		default:
			log.WithFields(log.Fields{
				"format": listOutputFormat,
			}).Warning("unrecognized output format")
			fallthrough
		case "basic", "":
			printBasicRecordsList(domain, records.Records)
		}
	}
}

func parseRequestedDomains(args []string) (domains []string) {
	if len(args) > 0 {
		log.WithField("args", args).Debug("using domain list from args")
		domains = args
	} else {
		log.Debug("querying account for domain list")
		domainList, domainListErr := GetClient().ListDomains()
		if domainListErr != nil {
			log.WithError(domainListErr).Error("error fetching account's domains")
		}

		for d := range domainList.Domains {
			domains = append(domains, d)
		}
	}

	return domains
}

func printBasicRecordsList(domain string, records []goname.DNSRecordResponse) {
	fmt.Println("#", domain)
	for _, record := range records {
		row := strings.Join([]string{
			record.RecordID,
			record.Name,
			record.Type,
			strings.TrimLeft(strings.Join([]string{record.Priority, record.Content}, " "), " "),
			record.TTL,
		}, " ")
		fmt.Println(row)
	}
}
