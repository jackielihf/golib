

package web

import "github.com/google/uuid"
import "github.com/gin-gonic/gin"

const key string = "USER_TRACK"
const longAfter int = 60 * 60 * 24 * 365 * 1000


func UserTrack(c *gin.Context) {
    if _, err := c.Cookie(key); err != nil {
        uuid := uuid.New().String()
        c.SetCookie(key, uuid, longAfter, "", "", false, false)    
    }
    c.Next()
}






