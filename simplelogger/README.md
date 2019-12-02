# simplelogger [![Dependabot Status](https://api.dependabot.com/badges/status?host=github&repo=ffflorian/go-tools)](https://dependabot.com)

This is a simple logger.

## Installation

Run `go get github.com/ffflorian/go-tools/simplelogger`.

## Usage

```go
var logger = simplelogger.New("my-app", false, true)
                            //    |       |      |
                            //    > your log prefix
                            //            |      |
                            //            > initially enabled?
                            //                   |
                            //                   > check environment variables?


// do something

logger.Enabled = true

logger.Log("Hello, world!")
```

```go
var logger = &simplelogger.Logger{
	Enabled: false
	Prefix:  "my-app"
}

// do something

logger.Enabled = true

logger.Log("Hello, world!")
```

### Output

<pre>
<b>my-app</b> Hello, world!
</pre>

## Test

```
go test
```
