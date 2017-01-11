

package web

import "strings"

import "github.com/gin-gonic/gin"

const (
    allowMethods string = "GET,HEAD,PUT,PATCH,POST"
    allowHeaders string = ""
)

func Cors(c *gin.Context) {
    
    method := strings.ToUpper(c.Request.Method)
    if method == "OPTIONS" {
        c.Header("Access-Control-Allow-Origin", "*")
        c.Header("Access-Control-Allow-Methods", allowMethods)
        c.Header("Access-Control-Allow-Headers", allowHeaders)
        c.Status(200)
    }else{
        c.Next()
    }
}

