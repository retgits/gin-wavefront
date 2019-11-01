package main

import (
	"github.com/gin-gonic/gin"
	ginwavefront "github.com/retgits/gin-wavefront"
)

func main() {
	// Create an instance of the Gin Wavefront middleware
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
		panic(err.Error())
	}

	r := gin.New()
	r.Use(wfemitter)

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.Run(":8083")
}
