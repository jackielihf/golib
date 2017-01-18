
package main




import "github.com/gin-gonic/gin"
import "github.com/jackielihf/golib/web"
import "fmt"



func login(c *gin.Context) {
    username := c.Query("username")
    obj := map[string]string{"username": username}
    web.JwtSetCookie(c, obj)    
    c.JSON(200, gin.H{
        "message": username,
    })
    c.Abort()
}

func handler(c *gin.Context) {
    msg := c.Query("msg")

    userCtx := c.Request.Header.Get("user_ctx")
    fmt.Println(userCtx)

    obj := map[string]string{}
    err := web.JwtUserCtx(c, &obj)
    fmt.Println(obj, err)

    web.CacheControl(c, 60)
    web.Success(c, gin.H{
        "username": obj["username"],
        "msg": msg,
    })
}



func GetApiConfig() []web.Api{
    config := []web.Api {
        web.Api{"GET", "/login", web.ApiHandlers{login}},
        web.Api{"GET", "/example", web.ApiHandlers{handler}},
    }    
    return config
}


func GetApiGroupConfig() []web.Api{
    config := []web.Api {
        web.Api{"GET", "/group", web.ApiHandlers{handler}},
    }    
    return config
}

