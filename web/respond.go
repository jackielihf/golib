
package web

// respond with JSON
import "fmt"
import "github.com/gin-gonic/gin"


func respond(c *gin.Context, statusCode int, data gin.H) {
    if statusCode < 100 {
        statusCode = 200
    }
    if data == nil {
        data = gin.H{
            "code": 0,
            "msg": StatusCode[statusCode],
        }
    }
    c.JSON(statusCode, data)
    c.Abort()
}

func respondError(c *gin.Context, statusCode int, data gin.H) {
    if statusCode < 100 {
        statusCode = 500
    }
    respond(c, statusCode, data)
}

func Success(c *gin.Context, data map[string]interface{}) {
    respond(c, 200, data)
}

func BadRequest(c *gin.Context, data map[string]interface{}) {
    respondError(c, 400, data)
}

func Unauthorized(c *gin.Context, data map[string]interface{}) {
    respondError(c, 401, data)
}

func Forbidden(c *gin.Context, data map[string]interface{}) {
    respondError(c, 403, data)
}

func NotFound(c *gin.Context, data map[string]interface{}) {
    respondError(c, 404, data)
}

func ServerError(c *gin.Context, data map[string]interface{}) {
    respondError(c, 500, data)
}

//cache
func CacheControl(c *gin.Context, max int) {
    c.Header("Cache-Control", fmt.Sprintf("max-age=%d", max))
}







