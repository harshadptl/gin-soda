# gin-soda


```golang

package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/harshadptl/gin-soda"
)

func main() {
	api := gin.Default()
	api.Use(soda.SodaMiddleware())
	api.GET("/dummy", GetDummyEndpoint)
	api.Run(":5000")
}


func GetDummyEndpoint(c *gin.Context) {
	resp := map[string]string{"hello": "world"}
	c.JSON(200, resp)
	data, _ := json.Marshal(resp)
	soda.SetRespCache(c, data)
}

```