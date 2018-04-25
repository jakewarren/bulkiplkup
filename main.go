package main

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"text/tabwriter"

	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
	"github.com/jakewarren/ipisp"
	"github.com/ogier/pflag"
)

func main() {

	filePath := pflag.StringP("output", "o", "", "output file name")
	csvOutput := pflag.BoolP("csv", "c", false, "output in CSV format")
	jsonOutput := pflag.BoolP("json", "j", false, "output in JSON format")
	displayHelp := pflag.BoolP("help", "h", false, "display help")

	pflag.Parse()

	// override the default usage display
	if *displayHelp {
		displayUsage()
		os.Exit(0)
	}

	//human-friendly CLI output
	log.SetHandler(cli.New(os.Stderr))

	//set the logging level
	log.SetLevel(log.WarnLevel)

	r := openStdinOrFile()

	scanner := bufio.NewScanner(r)

	ipList := make([]net.IP, 0)
	for scanner.Scan() {
		ip := net.ParseIP(scanner.Text())
		if IsRoutable(ip) && ip != nil {
			ipList = append(ipList, ip)
		} else {
			log.Errorf("non-routable IP: %s", ip)
		}

	}

	var client ipisp.Client

	if len(ipList) <= 10 {
		client, _ = ipisp.NewDNSClient()
		log.Trace("using DNS client")
	} else {
		client, _ = ipisp.NewWhoisClient()
		log.Trace("using whois client")
	}

	resp, lkupErr := client.LookupIPs(ipList)
	checkError("Error during lookup: ", lkupErr)

	resp = cleanResponses(resp)

	var f *os.File
	if *filePath == "" {
		f = os.Stdout
	} else {
		var err error
		f, err = os.Create(*filePath)
		checkError("Cannot create file", err)
		defer f.Close()
	}

	if *csvOutput {
		writeCSV(resp, f)
	} else if *jsonOutput {
		writeJSON(resp, f)
	} else {
		writeHuman(resp, f)
	}

}

func cleanResponses(resp []ipisp.Response) []ipisp.Response {
	output := make([]ipisp.Response, 0)

	for _, i := range resp {
		if i.IP != nil {
			output = append(output, i)
		}
	}
	return output
}

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
func writeCSV(resp []ipisp.Response, f *os.File) {
	w := csv.NewWriter(f)
	defer w.Flush()

	w.Write([]string{"IP", "Location", "ASN", "ISP", "Range"})

	for _, i := range resp {
		w.Write([]string{i.IP.String(), i.Country, i.ASN.String(), i.Name.String(), i.Range.String()})
	}
}

func writeHuman(resp []ipisp.Response, f *os.File) {
	var w *tabwriter.Writer
	w = tabwriter.NewWriter(f, 0, 0, 1, ' ', tabwriter.Debug)
	defer w.Flush()

	fmt.Fprintf(w, "IP\tLOC\tASN\tISP\tRange\n")

	for _, i := range resp {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", i.IP, i.Country, i.ASN, i.Name, i.Range)
	}

}

// print custom usage instead of the default provided by pflag
func displayUsage() {

	fmt.Printf("Usage: bulkiplkup [<flags>] [FILE]\n\n")
	fmt.Printf("Optional flags:\n\n")
	pflag.PrintDefaults()
}

func openStdinOrFile() io.Reader {
	var err error
	r := os.Stdin
	if len(pflag.Args()) > 1 {
		r, err = os.Open(os.Args[1])
		if err != nil {
			panic(err)
		}
	}
	return r
}

func checkError(message string, err error) {
	if err != nil {
		log.Errorf("%s: %s", message, err)
	}
}

//IPRange stores an IP range
type IPRange struct {
	from, to net.IP
}

// BogonRanges is a subset of the more static IPv4 Bogon/Reserved/Private ranges.
// In other words, these ranges are such fucking bogon's that they aren't even
// out in public.
var BogonRanges = []IPRange{
	{from: net.ParseIP("0.0.0.0"), to: net.ParseIP("0.255.255.255")},
	{from: net.ParseIP("10.0.0.0"), to: net.ParseIP("10.255.255.255")},
	{from: net.ParseIP("100.64.0.0"), to: net.ParseIP("10.127.255.255")},
	{from: net.ParseIP("127.0.0.0"), to: net.ParseIP("127.255.255.255")},
	{from: net.ParseIP("169.254.0.0"), to: net.ParseIP("169.254.255.255")},
	{from: net.ParseIP("172.16.0.0"), to: net.ParseIP("172.31.255.255")},
	{from: net.ParseIP("192.0.0.0"), to: net.ParseIP("192.0.0.255")},
	{from: net.ParseIP("192.0.2.0"), to: net.ParseIP("192.0.2.255")},
	{from: net.ParseIP("192.88.99.0"), to: net.ParseIP("192.88.99.255")},
	{from: net.ParseIP("192.168.0.0"), to: net.ParseIP("192.168.255.255")},
	{from: net.ParseIP("198.18.0.0"), to: net.ParseIP("198.19.255.255")},
	{from: net.ParseIP("198.51.100.0"), to: net.ParseIP("198.51.100.255")},
	{from: net.ParseIP("203.0.113.0"), to: net.ParseIP("203.0.113.255")},
	{from: net.ParseIP("224.0.0.0"), to: net.ParseIP("239.255.255.255")},
	{from: net.ParseIP("240.0.0.0"), to: net.ParseIP("255.255.255.255")},
}

// IsRoutable returns true if the IP is a publicly routable address
func IsRoutable(ip net.IP) bool {
	for _, rr := range BogonRanges {
		if rr.Contains(ip) {
			return false
		}
	}
	return true
}

// Contains checks if a given IP is in the IPRange
func (r *IPRange) Contains(ip net.IP) bool {
	return (bytes.Compare(ip, r.from) >= 0 && bytes.Compare(ip, r.to) <= 0)
}
