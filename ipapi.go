package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/apex/log"
)

func EnrichIPAPI(ipList []net.IP) {
	var resultData Data
	resultData.IPDetails = make(map[string]IPInfo)

	chunkSize := 100 // the API accepts a maximum of 100 IPs per request
	numRequests := (len(ipList) + chunkSize - 1) / chunkSize
	for i := 0; i < len(ipList); i += chunkSize {
		end := i + chunkSize
		if end > len(ipList) {
			end = len(ipList)
		}

		// convert from []net.IP to []string
		var ipStrings []string
		for _, ip := range ipList[i:end] {
			ipStrings = append(ipStrings, ip.String())
		}

		data := Payload{
			Ips: ipStrings,
		}
		payloadBytes, err := json.Marshal(data)
		if err != nil {
			log.Fatalf("Error marshaling JSON: %v", err)
		}
		body := bytes.NewReader(payloadBytes)

		log.WithField("number of IPs", len(data.Ips)).Debug("Requesting IP info from the API")

		req, err := http.NewRequest("POST", fmt.Sprintf("https://api.ipapi.is?key=%s", *apiKey), body)
		if err != nil {
			log.Fatalf("Error creating request: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{
			Timeout: 30 * time.Second,
		}
		resp, err := client.Do(req)
		if err != nil {
			log.Fatalf("Error making request: %v", err)
		}
		defer resp.Body.Close()

		jsonData, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("Error reading response body: %v", err)
		}

		var tempData Data
		err = json.Unmarshal(jsonData, &tempData)
		if err != nil {
			fmt.Println("Received JSON data that could not be parsed:")
			fmt.Println(string(jsonData))
			log.Fatalf("Error parsing JSON: %v", err)
		}

		for ip, ipInfo := range tempData.IPDetails {
			resultData.IPDetails[ip] = ipInfo
		}

		// Sleep only if there are more API requests to be made
		if numRequests > 1 && i+chunkSize < len(ipList) {
			log.Debug("sleeping for 30 seconds")
			time.Sleep(30 * time.Second)
		}
	}

	// Calculate the is_suspicious field. This field is set to true if other fields indicate the IP is suspicious and warrants further review
	for ip, ipInfo := range resultData.IPDetails {

		// check if the company abuse score contains a 'High' or 'Very High' score
		if strings.Contains(ipInfo.Company.AbuserScore, "High") || strings.Contains(ipInfo.Company.AbuserScore, "Very High") {
			ipInfo.IsSuspicious = true
		}

		if ipInfo.IsAbuser || ipInfo.IsProxy || ipInfo.IsTor || ipInfo.IsDatacenter {
			ipInfo.IsSuspicious = true
		}

		if ipInfo.IsVPN != false {
			ipInfo.IsSuspicious = true
		}

		// Store the modified ipInfo back into the map
		resultData.IPDetails[ip] = ipInfo
	}

	if *jsonOutput {
		var jsonData []byte
		var err error

		// Create a slice to hold the IP details
		ipDetailsSlice := make([]IPInfo, 0, len(resultData.IPDetails))
		for _, ipInfo := range resultData.IPDetails {
			ipDetailsSlice = append(ipDetailsSlice, ipInfo)
		}

		// Marshal the IP details slice to JSON
		jsonData, err = json.MarshalIndent(ipDetailsSlice, "", "  ")
		if err != nil {
			log.WithError(err).Fatal("error marshaling JSON output")
		}
		fmt.Println(string(jsonData))
	} else {
		w := csv.NewWriter(os.Stdout)
		err := w.Write([]string{"ip", "country code", "asn", "asn_name", "asn_type", "asn_abuse_score", "company_name", "company_type", "company_abuse_score", "is_abuser", "is_vpn", "is_proxy", "is_tor", "is_datacenter", "is_crawler", "is_mobile", "is_suspicious"})
		if err != nil {
			log.WithError(err).Fatal("error writing record to csv")
		}

		for ip, result := range resultData.IPDetails {
			var record []string
			record = append(record, ip)
			record = append(record, result.Location.Country)
			record = append(record, fmt.Sprintf("%d", result.ASN.ASN))
			record = append(record, result.ASN.Descr)
			record = append(record, result.ASN.Type)
			record = append(record, result.ASN.AbuserScore)
			record = append(record, result.Company.Name)
			record = append(record, result.Company.Type)
			record = append(record, result.Company.AbuserScore)
			record = append(record, fmt.Sprintf("%t", result.IsAbuser))

			// Sometimes the is_vpn field is enriched with the VPN provider, otherwise it's just a boolean
			switch value := result.IsVPN.(type) {
			case string:
				record = append(record, value)
			case bool:
				record = append(record, fmt.Sprintf("%t", value))
			default:
				record = append(record, "")
			}

			record = append(record, fmt.Sprintf("%t", result.IsProxy))
			record = append(record, fmt.Sprintf("%t", result.IsTor))
			record = append(record, fmt.Sprintf("%t", result.IsDatacenter))
			record = append(record, fmt.Sprintf("%t", result.IsCrawler))
			record = append(record, fmt.Sprintf("%t", result.IsMobile))
			record = append(record, fmt.Sprintf("%t", result.IsSuspicious))

			if err := w.Write(record); err != nil {
				log.WithError(err).Fatal("error writing record to csv")
			}
		}

		w.Flush()

		if err := w.Error(); err != nil {
			log.WithError(err).Fatal("error writing output")
		}
	}
}

type StringOrSlice []string

func (s *StringOrSlice) UnmarshalJSON(data []byte) error {
	if data[0] == '[' {
		var slice []string
		if err := json.Unmarshal(data, &slice); err != nil {
			return err
		}
		*s = StringOrSlice(slice)
		return nil
	}

	var single string
	if err := json.Unmarshal(data, &single); err != nil {
		return err
	}
	*s = StringOrSlice([]string{single})
	return nil
}

type Company struct {
	Name        string `json:"name"`
	AbuserScore string `json:"abuser_score"`
	Domain      string `json:"domain"`
	Type        string `json:"type"`
	Network     string `json:"network"`
	Whois       string `json:"whois"`
}

type Abuse struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	Email   string `json:"email"`
	Phone   string `json:"phone"`
}

type Datacenter struct {
	Datacenter string `json:"datacenter"`
	Domain     string `json:"domain"`
	Network    string `json:"network"`
}

type ASN struct {
	ASN         int           `json:"asn"`
	AbuserScore string        `json:"abuser_score"`
	Route       string        `json:"route"`
	Descr       string        `json:"descr"`
	Country     string        `json:"country"`
	Active      bool          `json:"active"`
	Org         string        `json:"org"`
	Domain      string        `json:"domain"`
	Abuse       StringOrSlice `json:"abuse"`
	Type        string        `json:"type"`
	Created     string        `json:"created"`
	Updated     string        `json:"updated"`
	RIR         string        `json:"rir"`
	Whois       string        `json:"whois"`
}

type Location struct {
	Continent     string  `json:"continent"`
	Country       string  `json:"country"`
	CountryCode   string  `json:"country_code"`
	State         string  `json:"state"`
	City          string  `json:"city"`
	Latitude      float64 `json:"latitude"`
	Longitude     float64 `json:"longitude"`
	Zip           string  `json:"zip"`
	Timezone      string  `json:"timezone"`
	LocalTime     string  `json:"local_time"`
	LocalTimeUnix int64   `json:"local_time_unix"`
	IsDST         bool    `json:"is_dst"`
	Accuracy      int     `json:"accuracy,omitempty"`
}

type IPInfo struct {
	IP           string      `json:"ip"`
	RIR          string      `json:"rir"`
	IsBogon      bool        `json:"is_bogon"`
	IsMobile     bool        `json:"is_mobile"`
	IsCrawler    bool        `json:"is_crawler"`
	IsDatacenter bool        `json:"is_datacenter"`
	IsTor        bool        `json:"is_tor"`
	IsProxy      bool        `json:"is_proxy"`
	IsVPN        interface{} `json:"is_vpn"`
	IsAbuser     bool        `json:"is_abuser"`
	IsSuspicious bool        `json:"is_suspicious"` // This is a calculated field based on other fields, indicating the IP is suspicious and warrants further review
	Company      Company     `json:"company"`
	Abuse        Abuse       `json:"abuse"`
	Datacenter   Datacenter  `json:"datacenter"`
	ASN          ASN         `json:"asn"`
	Location     Location    `json:"location"`
}

type Data struct {
	TotalElapsedMS float64 `json:"total_elapsed_ms"`
	IPDetails      map[string]IPInfo
}

func (d *Data) UnmarshalJSON(data []byte) error {
	var temp struct {
		TotalElapsedMS float64 `json:"total_elapsed_ms"`
	}
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}
	d.TotalElapsedMS = temp.TotalElapsedMS

	var ipDetails map[string]json.RawMessage
	if err := json.Unmarshal(data, &ipDetails); err != nil {
		return err
	}

	d.IPDetails = make(map[string]IPInfo)
	for ip, rawData := range ipDetails {
		if ip == "total_elapsed_ms" {
			continue
		}
		var ipInfo IPInfo
		if err := json.Unmarshal(rawData, &ipInfo); err != nil {
			return err
		}
		d.IPDetails[ip] = ipInfo
	}

	return nil
}

type Payload struct {
	Ips []string `json:"ips"`
}
