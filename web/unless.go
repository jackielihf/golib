
package web

import "github.com/gin-gonic/gin"
import "regexp"


/* 
usage:

r := gin.Default()
conditionalMiddleWare := web.Unless("/login").Then(otherMiddleWare)
r.Use(conditionalMiddleWare)  // If HTTP url don't match "/login", otherMiddleWare will be executed.
*/ 

type UnlessWare struct {
    apiRegexp *regexp.Regexp  
}

// if HTTP url matches apiRegexp, then skip fn, otherwise invoke fn.
// fn can be a middleware or a web controller.
func (that *UnlessWare) Then(fn gin.HandlerFunc) gin.HandlerFunc{

    // return middleware
    return func (c *gin.Context) {
        method := c.Request.Method
        path := c.Request.URL.Path
        // skip when matching
        if that.apiRegexp.MatchString(method + ";" + path) {
            c.Next()
        }else{
            // go into fn
            fn(c)
        }
    }
}

// new a Unless ware
// @Params reg [string] regexp of HTTP Path
func Unless(reg string) *UnlessWare{
    // new a unless ware
    ware := new(UnlessWare)
    ware.apiRegexp = regexp.MustCompile(reg)
    return ware
}




