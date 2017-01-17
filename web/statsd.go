
// statsd middleware
package web

import "time"
import "github.com/gin-gonic/gin"
import "github.com/jackielihf/golib"

var client = golib.Statsd()

func StatsdIncrement(key string) gin.HandlerFunc{
    return func (c *gin.Context) {
        client.Increment(key)    
        c.Next()
    }
}


func StatsdGauge(key string) gin.HandlerFunc{
    return func (c *gin.Context) {
        start := time.Now()
        // handle request
        c.Next()
        end := time.Now()
        latency := end.Sub(start)
        client.Gauge(key, latency.Seconds()*1000)    
    }
}