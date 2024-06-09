# bulkiplkup
[![Build Status](https://github.com/jakewarren/bulkiplkup/workflows/lint/badge.svg)](https://github.com/jakewarren/bulkiplkup/actions)
[![GitHub release](http://img.shields.io/github/release/jakewarren/bulkiplkup.svg?style=flat-square)](https://github.com/jakewarren/bulkiplkup/releases])
[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat-square)](https://github.com/jakewarren/bulkiplkup/blob/master/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/jakewarren/bulkiplkup)](https://goreportcard.com/report/github.com/jakewarren/bulkiplkup)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg?style=shields)](http://makeapullrequest.com)
> perform a bulk lookup of IP addresses

This tool assists with enriching a large amount of IP addresses with additonal information. If an [ipapi.is](https://ipapi.is/) API key is provided, additional information will be fetched. The data from ipapi.is is extended with an additional field called `is_suspicious` which is set to true if the IP is a known abuser, VPN, proxy, Tor exit node, datacenter, or the company's abuse score is 'High' or 'Very High'. 

> [!NOTE]
> The output format of the ipapi.is is opioninated and was designed to facilitate threat hunts against log data. To receive all available information use the JSON output, which can then be filtered as needed. 

## Install
### Option 1: Binary

Download the latest release from [https://github.com/jakewarren/bulkiplkup/releases/latest](https://github.com/jakewarren/bulkiplkup/releases/latest)

### Option 2: From source

```
go install github.com/jakewarren/bulkiplkup@latest
```

## Example
### Enriched with ipapi.is
```
❯ echo "8.8.8.8" | bulkiplkup 
ip,country code,asn,asn_name,asn_type,asn_abuse_score,company_name,company_type,company_abuse_score,is_abuser,is_vpn,is_proxy,is_tor,is_datacenter,is_crawler,is_mobile,is_suspicious
8.8.8.8,United States,15169,"GOOGLE, US",hosting,0 (Very Low),Google LLC,hosting,0.0039 (Low),true,true,false,false,true,false,false,true
```

### Enriched with Team Cymru's IP to ASN mapping service
```
❯ echo "8.8.8.8" | bulkiplkup 
IP      |LOC |ASN     |ISP            |Range
8.8.8.8 |US  |AS15169 |Google LLC, US |8.8.8.0/24
```

## Usage

`bulkiplkup` reads newline separated IP addresses from a file or STDIN.

To enrich IPs with [ipapi.is](https://ipapi.is/), provide an API key in the `IPAPI_KEY` environment variable or as a parameter. If the key is not avaiable the program will fall back to Team Cymru's IP to ASN mapping service.

```
❯ bulkiplkup -h
Usage: bulkiplkup [<flags>] [FILE]

Optional flags:

  -k, --api-key="": API key for ipapi.is. Also accepts the IPAPI_KEY environment variable.
  -c, --csv=true: output in CSV format
  -h, --help=false: display help
  -j, --json=false: output in JSON format
  -v, --verbose=false: verbose output
```

## Acknowledgements

* Team Cymru for hosting their excellent IP to ASN mapping service - http://www.team-cymru.com/IP-ASN-mapping.html
* https://github.com/ammario/ipisp/ Golang IP to ISP library utilizing team cymru's IP to ASN service

## Changes

All notable changes to this project will be documented in the [changelog].

The format is based on [Keep a Changelog](http://keepachangelog.com/) and this project adheres to [Semantic Versioning](http://semver.org/).


[changelog]: https://github.com/jakewarren/bulkiplkup/blob/master/CHANGELOG.md
