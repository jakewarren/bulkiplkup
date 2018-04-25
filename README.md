# bulkiplkup
[![Build Status](https://travis-ci.org/jakewarren/bulkiplkup.svg?branch=master)](https://travis-ci.org/jakewarren/bulkiplkup/)
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
go get github.com/jakewarren/bulkiplkup
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

Team Cymru for hosting their excellent IP to ASN mapping service - http://www.team-cymru.com/IP-ASN-mapping.html

## Changes

All notable changes to this project will be documented in the [changelog].

The format is based on [Keep a Changelog](http://keepachangelog.com/) and this project adheres to [Semantic Versioning](http://semver.org/).

## License

MIT © 2018 Jake Warren

[changelog]: https://github.com/jakewarren/bulkiplkup/blob/master/CHANGELOG.md
