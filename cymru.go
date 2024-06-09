package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"text/tabwriter"

	"github.com/apex/log"
	"github.com/jakewarren/ipisp"
)

func EnrichCymru(ipList []net.IP) {
	var client ipisp.Client

	// if looking up more than 10 IPs, use the whois client
	if len(ipList) <= 10 {
		client, _ = ipisp.NewDNSClient()
		log.Debug("using DNS client")
	} else {
		client, _ = ipisp.NewWhoisClient()
		log.Debug("using whois client")
	}

	resp, lkupErr := client.LookupIPs(ipList)
	checkError("Error during lookup: ", lkupErr)

	resp = cleanResponses(resp)

	f := os.Stdout

	if *csvOutput {
		writeCSV(resp, f)
	} else if *jsonOutput {
		writeJSON(resp, f)
	} else {
		writeHuman(resp, f)
	}
}

// writeCSV outputs the response as CSV
func writeCSV(resp []ipisp.Response, f *os.File) {
	w := csv.NewWriter(f)
	defer w.Flush()

	w.Write([]string{"IP", "Location", "ASN", "ISP", "Range"})

	for _, i := range resp {
		w.Write([]string{i.IP.String(), i.Country, i.ASN.String(), i.Name.String(), i.Range.String()})
	}
}

// writeHuman outputs the response as a pretty tabular output
func writeHuman(resp []ipisp.Response, f *os.File) {
	w := tabwriter.NewWriter(f, 0, 0, 1, ' ', tabwriter.Debug)
	defer w.Flush()

	fmt.Fprintf(w, "IP\tLOC\tASN\tISP\tRange\n")

	for _, i := range resp {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", i.IP, i.Country, i.ASN, i.Name, i.Range)
	}

}

// cleanResponses removes any responses with nil IP entries
func cleanResponses(resp []ipisp.Response) []ipisp.Response {
	output := make([]ipisp.Response, 0)

	for _, i := range resp {
		if i.IP != nil {
			output = append(output, i)
		}
	}
	return output
}

// writeJSON outputs the responses as JSON
func writeJSON(resp []ipisp.Response, f *os.File) {
	type Record struct {
		IP      string
		Country string
		ASN     string
		ISP     string
		Range   string
	}

	records := make([]Record, 0)
	for _, i := range resp {
		records = append(records, Record{i.IP.String(), i.Country, i.ASN.String(), i.Name.String(), i.Range.String()})
	}

	rec, _ := json.MarshalIndent(records, "", "    ")

	fmt.Fprint(f, string(rec))

}
