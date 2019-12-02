# gh-open [![Build Status](https://action-badges.now.sh/ffflorian/go-tools)](https://github.com/ffflorian/go-tools/actions/) [![Dependabot Status](https://api.dependabot.com/badges/status?host=github&repo=ffflorian/go-tools)](https://dependabot.com)

Open a GitHub repository in your browser.

## Installation

Run `go get github.com/ffflorian/go-tools/gh-open`.

## Usage

```
Open a GitHub repository in your browser.

Usage:
  gh-open [options] [directory]

Options:
  --timeout, -t      Set a custom timeout for HTTP requests
  --print, -p        just print the URL
  --branch, -b       open the branch tree (and not the PR)
  --debug, -d        enable debug mode
  --version, -v      output the version number
  --help, -h         output usage information
```

Hint: You can also enable the debug mode by setting the environment variable `DEBUG` to "gh-open".

Example:

```
DEBUG="gh-open" gh-open
```

## Test

```
go test
```
