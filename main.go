package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
	"github.com/ogier/pflag"
)

var csvOutput *bool
var jsonOutput *bool
var verbose *bool
var apiKey *string

func main() {

	apiKey = pflag.StringP("api-key", "k", "", "API key for ipapi.is. Also accepts the IPAPI_KEY environment variable.")
	csvOutput = pflag.BoolP("csv", "c", true, "output in CSV format")
	jsonOutput = pflag.BoolP("json", "j", false, "output in JSON format")
	verbose = pflag.BoolP("verbose", "v", false, "verbose output")
	displayHelp := pflag.BoolP("help", "h", false, "display help")

	pflag.Parse()

	// override the default usage display
	if *displayHelp {
		displayUsage()
		os.Exit(0)
	}

	//human-friendly CLI output
	log.SetHandler(cli.New(os.Stderr))
	//set the default logging level
	log.SetLevel(log.WarnLevel)

	if *verbose {
		log.SetLevel(log.DebugLevel)
	}

	r := openStdinOrFile()

	scanner := bufio.NewScanner(r)

	ipList := make([]net.IP, 0)
	for scanner.Scan() {
		input := strings.TrimSpace(scanner.Text())
		ip := net.ParseIP(input)
		if ip == nil {
			log.Errorf("error parsing IP: %s", input)
		}
		if IsRoutable(ip) && ip != nil {
			ipList = append(ipList, ip)
		} else {
			log.Warnf("warning: non-routable IP: %s", ip)
		}

	}

	// check if the API key was provided as a flag or environment variable
	if *apiKey == "" {
		*apiKey = os.Getenv("IPAPI_KEY")
	}

	if *apiKey == "" {
		// if no API key was provided, fall back to the Team CYMRU service
		EnrichCymru(ipList)
	} else {
		// if an API key was provided, use the ipapi.is service
		EnrichIPAPI(ipList)
	}

}

// print custom usage instead of the default provided by pflag
func displayUsage() {

	fmt.Printf("Usage: bulkiplkup [<flags>] [FILE]\n\n")
	fmt.Printf("Optional flags:\n\n")
	pflag.PrintDefaults()
}
