# gin-wavefront

[![Go Report Card](https://goreportcard.com/badge/github.com/retgits/gin-wavefront?style=flat-square)](https://goreportcard.com/report/github.com/retgits/gin-wavefront)
[![Godoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/retgits/gin-wavefront)
![GitHub](https://img.shields.io/github/license/retgits/gin-wavefront?style=flat-square)
![GitHub tag (latest SemVer)](https://img.shields.io/github/v/tag/retgits/gin-wavefront?sort=semver&style=flat-square)

> gin-wavefront is a [Gin middleware](https://github.com/gin-gonic) to emit metrics to [Wavefront](https://www.wavefront.com/).

## Prerequisites

To use this Gin middleware, you'll need to have

* [Go (at least Go 1.12)](https://golang.org/dl/)
* [A Wavefront account](https://www.wavefront.com/sign-up/)
* [A Wavefront API key](https://docs.wavefront.com/wavefront_api.html)

## Installation

Using `go get`

```bash
go get github.com/retgits/gin-wavefront
```

## Usage

To start, you'll need to initialize the Wavefront emitter:

```go
wfconfig := &ginwavefront.WavefrontConfig{
    Server:        "https://<INSTANCE>.wavefront.com",
    Token:         "my-api-key",
    BatchSize:     10000,
    MaxBufferSize: 50000,
    FlushInterval: 1,
    Source:        "my-app",
    MetricPrefix:  "my.awesome.app",
    PointTags:     make(map[string]string),
}
wfemitter, err := ginwavefront.WavefrontEmitter(wfconfig)
if err != nil {
    fmt.Println(err.Error())
}
```

Now you can use the `wfemitter` as a middleware function in Gin

```go
r := gin.New()
r.Use(wfemitter)
```

A complete sample app can be found in the [examples](./examples) folder

## Contributing

[Pull requests](https://github.com/retgits/gin-wavefront/pulls) are welcome. For major changes, please open [an issue](https://github.com/retgits/gin-wavefront/issues) first to discuss what you would like to change.

Please make sure to update tests as appropriate.

## License

See the [LICENSE](./LICENSE) file in the repository
