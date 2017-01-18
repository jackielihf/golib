

package web

import "github.com/google/uuid"
import "github.com/gin-gonic/gin"

const USER_TRACK_KEY string = "USER_TRACK"
const longAfter int = 60 * 60 * 24 * 365 * 1000

func UserTrack() gin.HandlerFunc {
    return func(c *gin.Context) {
        if _, err := c.Cookie(USER_TRACK_KEY); err != nil {
            uuid := uuid.New().String()
            c.SetCookie(USER_TRACK_KEY, uuid, longAfter, "", "", false, false)    
        }
        c.Next()
    }
}





