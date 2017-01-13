

package web

import "strings"

import "github.com/gin-gonic/gin"

const (
    CORS_ALLOW_METHODS string = "GET,HEAD,PUT,PATCH,POST"
    CORS_ALLOW_HEADERS string = ""
)

func Cors(c *gin.Context) {
    
    method := strings.ToUpper(c.Request.Method)
    if method == "OPTIONS" {
        c.Header("Access-Control-Allow-Origin", "*")
        c.Header("Access-Control-Allow-Methods", CORS_ALLOW_METHODS)
        c.Header("Access-Control-Allow-Headers", CORS_ALLOW_HEADERS)
        c.Status(200)
        c.Abort()
    }else{
        c.Next()
    }
}

