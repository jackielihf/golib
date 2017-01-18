

package main


import "os"
import "fmt"

import "github.com/gin-gonic/gin"
import "github.com/jackielihf/golib/web"


func main() {
    // engine
    r := gin.New()
    r.Use(gin.Recovery())

    r.GET("/ping", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "message": "pong",
        })
    })

    // middleware
    r.Use(web.Cors())
    r.Use(web.UserTrack())
    r.Use(web.Logger())
    r.Use(web.Unless("/login").Then(web.JwtMiddleWare()))
    
    r.Use(web.StatsdIncrement("example.go_web_starter.count"))
    r.Use(web.StatsdGauge("example.go_web_starter.elapse"))
    
    // router
    config := GetApiConfig()
    web.CreateApi(r, config)

    config2 := GetApiGroupConfig()
    web.CreateGroup(r, "/api", config2)
    
    // listen
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }
    fmt.Println("web server listening on: " + port) 
    r.Run(":" + port)
    
}