# my-timezone [![Build Status](https://github.com/ffflorian/wire-bots/workflows/Build/badge.svg)](https://github.com/ffflorian/go-tools/actions/) [![Dependabot Status](https://api.dependabot.com/badges/status?host=github&repo=ffflorian/go-tools)](https://dependabot.com)

Calculate the exact time based on your location.

## Installation

Run `go get github.com/ffflorian/go-tools/my-timezone`.

## Usage

```
Calculate the exact time based on your location.

For the `location` argument you can either use coordinates (e.g. 52.5502,13.4304)
or a city name (e.g. "Berlin, Germany").

Usage:
  my-timezone [options] [location]

Options:
  --server, -s       set the NTP server (default is "pool.ntp.org")
  --timeout, -t      set a custom timeout for HTTP requests (default is 2000ms)
  --version, -v      output the version number
  --help, -h         output usage information
  --offline, -o      enable offline mode (disables city matching)
  --debug, -d        enable debug mode
```

Hint: You can also enable the debug mode by setting the environment variable `DEBUG` to "my-timezone".

Example:

```
DEBUG="my-timezone" my-timezone

### Example

```
$ my-timezone 52.550215,13.430428 -o   # Berlin, Germany
Your personal time in "52.550215,13.430428": 21:24:49

$ my-timezone "Marburg, Germany"
Your personal time in "Marburg, Germany": 21:06:57
```
```

## Test

```
go test
```
