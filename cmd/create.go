package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/aelindeman/goname"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// createCmd represents the set command
var createCmd = &cobra.Command{
	Use:     "create <domain> <host> <type> <content>",
	Short:   "Create a DNS record in a domain",
	Args:    cobra.MinimumNArgs(4),
	Aliases: []string{"add"},
	Run:     createDNSRecord,
}

var (
	targetDomain string
	setHostname  string
	setType      string
	setPriority  int
	setContent   string
	setTTL       int
)

var validDNSRecordTypes = []string{
	"A", "AAAA", "CNAME", "MX", "NS", "PTR", "SOA", "SRV", "TXT",
}

func init() {
	rootCmd.AddCommand(createCmd)
	createCmd.PersistentFlags().IntVar(&setTTL, "ttl", 3600, "time-to-live")
}

func createDNSRecord(cmd *cobra.Command, args []string) {
	record, parseErr := parseRecordFromArgs(args)
	if parseErr != nil {
		log.WithError(parseErr).Fatal("couldn't parse record from arguments")
	}

	resp, err := GetClient().CreateDNSRecord(targetDomain, record)
	if err != nil {
		log.WithError(err).WithFields(log.Fields{
			"domain": targetDomain,
		}).Fatal("could not create DNS record")
	}

	log.WithFields(log.Fields{
		"recordID":   resp.RecordID,
		"name":       resp.Name,
		"type":       resp.Type,
		"content":    resp.Content,
		"priority":   resp.Priority,
		"createDate": resp.CreateDate,
		"ttl":        resp.TTL,
	}).Info("record created successfully")
}

func parseRecordFromArgs(args []string) (goname.DNSRecordRequest, error) {
	var record goname.DNSRecordRequest

	targetDomain = args[0]
	setHostname = args[1]
	setType = strings.ToUpper(args[2])
	setContent = strings.Join(args[3:], " ")

	log.WithFields(log.Fields{
		"domain":   targetDomain,
		"hostname": setHostname,
		"type":     setType,
		"content":  setContent,
	}).Debug("parsed arguments")

	if validationErr := validateInput(); validationErr != nil {
		return record, validationErr
	}

	record = goname.DNSRecordRequest{
		Hostname: setHostname,
		Type:     setType,
		Content:  setContent,
		TTL:      setTTL,
	}

	if setType == "MX" || setType == "SRV" {
		var err error
		setPriority, setContent, err = splitPriority(setContent)
		if err == nil {
			record.Priority = setPriority
			record.Content = setContent
		}
	}

	return record, nil
}

func splitPriority(content string) (int, string, error) {
	priorityAndContent := strings.SplitN(content, " ", 2)
	ttl, err := strconv.Atoi(priorityAndContent[0])
	return ttl, priorityAndContent[1], err
}

func validateInput() error {
	if strings.HasSuffix(setHostname, targetDomain) {
		return fmt.Errorf("don't suffix the target domain as part of the host field")
	}
	if setTTL < 60 {
		return fmt.Errorf("TTL must be at least 60 seconds")
	}

	for _, valid := range validDNSRecordTypes {
		if setType == valid {
			return nil
		}
	}
	return fmt.Errorf("%s is not a valid DNS record type", setType)
}
