# golib
Some common, handy tools in Go language. 
You can easily build up your web applications(based gin) by using these tools.

## Current version
v0.1.2

## Installation
go get github.com/jackielihf/golib

# Documentation
## package golib
## package web
### router
For a web application, you can put all the api definitions together, and create routes in batch by using MODULE router.

* func CreateApi(router *gin.Engine, config []Api) *gin.Engine
* func CreateGroup(router *gin.Engine, prefix string, config []Api) *gin.Engine


api definitions
`
func GetApiConfig() []web.Api{
    config := []web.Api {
        web.Api{"GET", "/login", web.ApiHandlers{login}},
        web.Api{"GET", "/example", web.ApiHandlers{handler}},
    }    
    return config
}
`
simple web application
`
import "github.com/gin-gonic/gin"
import "github.com/jackielihf/golib/web"

func main() {
    r := gin.New()
    r.Use(gin.Recovery())

    // router
    config := GetApiConfig()
    web.CreateApi(r, config)

    config2 := GetApiGroupConfig()
    web.CreateGroup(r, "/api", config2)
    
    // listen
    r.Run(":8080")
}


# History
* v0.1.2 add logger, statsd, respond
* v0.1.1 add jwt singleton, unless
* v0.1.0 add golib/web/router, cors, usertrack
