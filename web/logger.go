
package web

import "fmt"
import "time"
import "bytes"
import "encoding/json"
import "github.com/jackielihf/golib"
import "github.com/gin-gonic/gin"


var log = golib.Logger()

// middleware
func Logger(headers ...string) gin.HandlerFunc{


    return func(c *gin.Context) {
        // Start timer
        start := time.Now()
        // process request
        c.Next()
        log.Info(format(c, start, headers))
    }

}


func format(c *gin.Context, start time.Time, headerKeys []string) string {
    var b *bytes.Buffer = new(bytes.Buffer)
    url := c.Request.URL.String()
    end := time.Now()
    latency := end.Sub(start)
    clientIP := c.ClientIP()
    method := c.Request.Method
    statusCode := c.Writer.Status() 
    size := c.Writer.Size()
    //header
    userAgent := c.Request.UserAgent()
    contentType := c.Request.Header.Get("content-type")
    var headers = map[string]string{}
    for _, key := range headerKeys {
        value := c.Request.Header.Get(key)
        headers[key] = value
    }
    headerString := ""
    if len(headers) > 0 {
        if headerBytes, err := json.Marshal(headers); err == nil {
            headerString = string(headerBytes)
        }    
    }
    //to line
    fmt.Fprintf(b, "%3d | %s %s | %13v | %dB | %s | %s | %s | %s", statusCode, method, url, latency, size, clientIP, contentType, userAgent, headerString)
    
    return b.String()

}



