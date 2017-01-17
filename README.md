# golib
Some common, handy tools in Go language. 
You can easily build up your web applications(based on [gin](https://github.com/gin-gonic/gin)) by using these tools.

## Current version
v0.1.2

## Installation
go get github.com/jackielihf/golib

# Documentation
## package golib
### logger
A logger client of a pluggable logger [logrus](https://github.com/sirupsen/logrus).  
support: console, local file, syslog, kafka  etc.  
setting with env variables:   

* app_name          - application name
* log_format        - json, default
* log_enable_file   - output to local file or not
* log_dir           - local log dir

```
import "github.com/jackielihf/golib"
var log = golib.Logger()
log.Info("some info msg")
log.Error("some error msg")
```
### statsd
A statsd client for collecting data.   
env variables:  

* **statsd_addr** - address of statsd server, default: localhost:8125

```
import "github.com/jackielihf/golib"
client := golib.Statsd()
client.Increment("example.count")
client.Gauge("example.delta")
```

## package web
### router
For a web application, you can put all the api definitions together, and create routes in batch by using MODULE router.

* func CreateApi(router *gin.Engine, config []Api) *gin.Engine
* func CreateGroup(router *gin.Engine, prefix string, config []Api) *gin.Engine


**api definitions**

```
func login(c *gin.Context) {
    username := c.Query("username")
    c.JSON(200, gin.H{
        "user": username,
    })
    c.Abort()
}
func handler(c *gin.Context) {
    msg := c.Query("msg")
    c.JSON(200, gin.H{
        "msg": msg,
    })
    c.Abort()  
}
func GetApiConfig() []web.Api{
    config := []web.Api {
        web.Api{"GET", "/login", web.ApiHandlers{login}},
        web.Api{"GET", "/example", web.ApiHandlers{handler}},
    }    
    return config
}
```
**simple web application**

```
import "github.com/gin-gonic/gin"
import "github.com/jackielihf/golib/web"
func main() {
    r := gin.New()
    r.Use(gin.Recovery())

    // router
    config := GetApiConfig()
    web.CreateApi(r, config)  // create api routes

    config2 := GetApiGroupConfig()
    web.CreateGroup(r, "/api", config2) // create an api route group, with prefix '/api'
    
    // listen
    r.Run(":8080")
}
```

### logger (middleware)

### JWT (middleware)

### unless (middleware)
### cors (middleware)
### usertrack (middleware)




# History
* v0.1.2 add logger, statsd, respond
* v0.1.1 add jwt singleton, unless
* v0.1.0 add golib/web/router, cors, usertrack
