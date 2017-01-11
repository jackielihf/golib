

package web

import "github.com/gin-gonic/gin"
import "fmt"
import "strings"

// gin.HandlerFunc 别名
type ApiHandlers []gin.HandlerFunc 

// Api 数据结构
type Api struct {
    Method string
    Path string
    Handlers ApiHandlers
}


// 创建API路由
func CreateApi(router *gin.Engine, config []Api) *gin.Engine{
    for _, api := range config {    
        Method := strings.ToUpper(api.Method)
        switch Method {
        case "GET":
            router.GET(api.Path, api.Handlers...)
        case "POST":
            router.POST(api.Path, api.Handlers...)
        case "PUT":
            router.PUT(api.Path, api.Handlers...)
        case "DELETE":
            router.DELETE(api.Path, api.Handlers...)
        case "HEAD":
            router.HEAD(api.Path, api.Handlers...)
        case "PATCH":
            router.PATCH(api.Path, api.Handlers...)
        case "OPTIONS":
            router.OPTIONS(api.Path, api.Handlers...)
        default:
            fmt.Println("unsuported method of api: " + Method)
        }
    
    }
    return router
}

// 创建API路由Group
// @param prefix [string] api路径的前缀
func CreateGroup(router *gin.Engine, prefix string, config []Api) *gin.Engine{
    group := router.Group(prefix)
    for _, api := range config {    
        Method := strings.ToUpper(api.Method)
        switch Method {
        case "GET":
            group.GET(api.Path, api.Handlers...)
        case "POST":
            group.POST(api.Path, api.Handlers...)
        case "PUT":
            group.PUT(api.Path, api.Handlers...)
        case "DELETE":
            group.DELETE(api.Path, api.Handlers...)
        case "HEAD":
            group.HEAD(api.Path, api.Handlers...)
        case "PATCH":
            group.PATCH(api.Path, api.Handlers...)
        case "OPTIONS":
            group.OPTIONS(api.Path, api.Handlers...)
        default:
            fmt.Println("unsuported method of api: " + Method)
        }
    
    }
    return router
}


