# gh-open [![Dependabot Status](https://api.dependabot.com/badges/status?host=github&repo=ffflorian/go-tools)](https://dependabot.com)

Open a GitHub repository in your browser.

## Installation

Run `go get github.com/ffflorian/go-tools/gh-open`.

## Usage

```
Usage: gh-open [options] [directory]

Open a GitHub repository in your browser. Opens pull requests by default.

Options:
  -p, --print             Just print the URL
  -b, --branch            Open the branch tree (and not the PR)
  -t, --timeout <number>  Set a custom timeout for HTTP requests
  -v, --version           output the version number
  -h, --help              output usage information
```

## Test

```
go test
```
