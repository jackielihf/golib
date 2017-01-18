# golib
Some common, handy tools in Go language. 
You can easily build up your web applications(based on [gin](https://github.com/gin-gonic/gin)) by using these tools.

Toolbox:
* web api router based on gin
* pluggable logger
* statsd client

**middleware**
* JWT
* cors
* logger (log http request)
* statsd (counting for http request)
* unless (for skipping handler)
* usertrack
* respond 


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

* **app_name**:       application name
* **log_format**:       json, custom formatter(default)
* **log_enable_file**:  output to local file or not
* **log_dir**:          local log dir, default: ".", filename=log_dir + "/info.log"

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
Logger middleware can record the information of every request, such as http path, method, elapsed time, content-type, content-length etc.
  
```
import "github.com/jackielihf/golib/web"
func main(){
    r := gin.New()
    r.Use(gin.Recovery())
    
    // middleware
    r.Use(web.Logger())
    // ...
}
```
**log example**

```
 2017-01-17T18:01:40+08:00 [info]  $ 401 | GET /example?msg=bbb |     294.866µs | 28B | ::1 | application/json | Mozilla/5.0 (Macintosh; Intel Mac OS X 10_9_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/54.0.2840.98 Safari/537.36 |
```

### JWT (middleware)
JSON Web Token(JWT) can be used to encrypt user information which set into Cookie, so as to implement user login context.
#### API
* func JwtMiddleWare() gin.HandlerFunc

    Return a middleware of processing JWT. It read and decrypts the JWT cookie(cookie name: JWT_TOKEN) for plaintext(JSON string). If JWT is valid, it parse the plaintext into an object and set into the header of the request(header name: **user_ctx**), then pass to the next handler, or it returns 401 to client.

* func JwtSetCookie(c *gin.Context, v interface{})   
    Encrypt object v and set cookie. Use it in LOGIN handler.
    
* func JwtClearCookie(c *gin.Context)  
    Clear the JWT cookie. Use it in LOGOUT handler. 
           
* func JwtUserCtx(c *gin.Context, v interface{}) error  
    Read the JWT cookie, and parse it into an object with reference type.
    
#### env variables
* **jwt_secret**: A secret string for encrypting. Default: a random uuid string.
* **jwt_domain**: JWT cookie's domain. Default: empty string.
* **jwt_expire**: JWT cookie's expire time (seconds). Default: 60 * 60 * 24 seconds.

```
// app.go
package main
import "github.com/gin-gonic/gin"
import "github.com/jackielihf/golib/web"

func main() {
    // engine
    r := gin.New()
    r.Use(gin.Recovery())
    
    // middleware
    r.Use(web.Logger())
    r.Use(web.Unless("/login").Then(web.JwtMiddleWare()))   // JWT middleware
    
    // router
    config := GetApiConfig()
    web.CreateApi(r, config)
    //listen    
    r.Run(":8080")
}

//api.go
func login(c *gin.Context) {
    username := c.Query("username")
    obj := map[string]string{"username": username}
    web.JwtSetCookie(c, obj)       // encrypt obj, and set cookie
    c.JSON(200, gin.H{
        "user": username,
    })
    c.Abort()
}
func handler(c *gin.Context) {
    msg := c.Query("msg")
    obj := map[string]string{}
    err := web.JwtUserCtx(c, &obj)  // read JWT cookie, and decrypt JWT into obj
    c.JSON(200, gin.H{
        "username": obj["username"]
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

### unless (middleware)
* func Unless(reg string) *UnlessWare
* func (that *UnlessWare) Then(fn gin.HandlerFunc) gin.HandlerFunc

A little tool that can be used to skip middleware handler when HTTP path matching the given regex string.

```
...
// JWT middleware. When path matching "/login", it will skip the JwtMiddleWare's handler.
r.Use(web.Unless("/login").Then(web.JwtMiddleWare()))   
...
```

### cors (middleware)
Support Cors cross domain.

```
...
// middleware
r.Use(web.Cors)
```
### usertrack (middleware)

Generate an uuid for every client. It is useful for counting PV, UV.

```
...
r.Use(web.UserTrack)
```


# History
* v0.2.0 add example
* v0.1.2 add logger, statsd, respond
* v0.1.1 add jwt singleton, unless
* v0.1.0 add golib/web/router, cors, usertrack
