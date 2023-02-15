# bulkiplkup
[![Build Status](https://github.com/jakewarren/bulkiplkup/workflows/lint/badge.svg)](https://github.com/jakewarren/bulkiplkup/actions)
[![GitHub release](http://img.shields.io/github/release/jakewarren/bulkiplkup.svg?style=flat-square)](https://github.com/jakewarren/bulkiplkup/releases])
[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat-square)](https://github.com/jakewarren/bulkiplkup/blob/master/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/jakewarren/bulkiplkup)](https://goreportcard.com/report/github.com/jakewarren/bulkiplkup)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg?style=shields)](http://makeapullrequest.com)
> perform a bulk lookup of IP addresses



## Install
### Option 1: Binary

Download the latest release from [https://github.com/jakewarren/bulkiplkup/releases/latest](https://github.com/jakewarren/bulkiplkup/releases/latest)

### Option 2: From source

```
go install github.com/jakewarren/bulkiplkup@latest
```

## Example

```
❯ echo "8.8.8.8" | bulkiplkup 
IP      |LOC |ASN     |ISP            |Range
8.8.8.8 |US  |AS15169 |Google LLC, US |8.8.8.0/24
```

## Usage

`bulkiplkup` reads newline separated IP addresses from a file or STDIN.

```
❯ bulkiplkup -h
Usage: bulkiplkup [<flags>] [FILE]

Optional flags:

  -c, --csv=false: output in CSV format
  -h, --help=false: display help
  -j, --json=false: output in JSON format
  -o, --output="": output file name
```

## Acknowledgements

* Team Cymru for hosting their excellent IP to ASN mapping service - http://www.team-cymru.com/IP-ASN-mapping.html
* https://github.com/42wim/dt/ for the inspiration of output format
* https://github.com/ammario/ipisp/ Golang IP to ISP library utilizing team cymru's IP to ASN service

## Changes

All notable changes to this project will be documented in the [changelog].

The format is based on [Keep a Changelog](http://keepachangelog.com/) and this project adheres to [Semantic Versioning](http://semver.org/).

## License

MIT © 2018 Jake Warren

[changelog]: https://github.com/jakewarren/bulkiplkup/blob/master/CHANGELOG.md
