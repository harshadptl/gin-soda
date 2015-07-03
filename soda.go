package soda

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/spaolacci/murmur3"
	"github.com/wunderlist/ttlcache"

	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

var cache *ttlcache.Cache

type response struct {
	Header      http.Header
	Data        []byte
	Status      int
	ContentType string
}

//SetRespCache ...
//save resp to cache
func SetRespCache(c *gin.Context, data []byte) error {
	if c.Request.Method != "GET" {
		return errors.New("request type is not GET")
	}

	if c.Writer.Status() != 200 {
		return errors.New("response code not 200")
	}

	key := keyFromRequest(c.Request)

	resp := &response{}
	resp.Header = c.Writer.Header()
	resp.Data = data
	resp.Status = c.Writer.Status()
	resp.ContentType = c.ContentType()

	value, err := json.Marshal(resp)
	if err != nil {
		cache.Set(key, string(value))
	}
	return nil
}

//SodaMiddleware ...
//returns a middleware which overrides the call
//and returns data from cache
func SodaMiddleware() gin.HandlerFunc {
	cache = ttlcache.NewCache(time.Second * 60)
	return func(c *gin.Context) {
		//pass if request method not GET
		if c.Request.Method != "GET" {
			return
		}

		key := keyFromRequest(c.Request)

		value, keyExists := cache.Get(key)
		if keyExists {
			if value != "" {
				return
			}

			var resp response
			err := json.Unmarshal([]byte(value), resp)
			if err != nil {
				cache.Set(key, "")
				return
			}

			for key, value := range resp.Header {
				c.Header(key, strings.Join(value, ", "))
			}
			c.Data(resp.Status, resp.ContentType, resp.Data)

			c.Abort()
		}
	}
}

func keyFromRequest(request *http.Request) string {
	requestURI := request.RequestURI
	hash := murmur3.Sum64([]byte(requestURI))
	key := fmt.Sprint(hash)
	return key
}
