// json web token



package web

import "os"
import "time"
import "strconv"
import "errors"
import "sync"

import "github.com/gin-gonic/gin"
import "github.com/dgrijalva/jwt-go"
import "github.com/google/uuid"

// cookie key name
const JWT_COOKIE_KEY string = "JWT_TOKEN"
const JWT_OBJ_KEY string = "user_context"
// singleton
var instance *Jwt
var once sync.Once

func getInstance() *Jwt {
    once.Do(func(){
        instance = new(Jwt)
        instance.Init()
    })
    return instance
}

// middleware
func JwtMiddleWare() gin.HandlerFunc{
    return getInstance().Ware
}

// sign and set cookie
func JwtSetCookie(c *gin.Context, payload string) {
    if payload != "" {
        obj := getInstance()
        if token, err := obj.sign(payload); err == nil {
            c.SetCookie(JWT_COOKIE_KEY, token, obj.Expire, "", obj.Domain, false, false)            
        } 
    }
}

// clear cookie
func JwtClearCookie(c *gin.Context) {
    c.SetCookie(JWT_COOKIE_KEY, "", 0, "", "", false, false)    
}



// Jwt 类
type Jwt struct {
    Secret string
    Domain string
    Expire int
    secretBytes []byte
}

func (obj *Jwt) Init() {
    //secret
    if obj.Secret == "" {
        obj.Secret = os.Getenv("JWT_SECRET")    
    }
    if obj.Secret == "" {
        obj.Secret = uuid.New().String()
    }
    obj.secretBytes = []byte(obj.Secret)
    //domain
    if obj.Domain == "" {
        obj.Domain = os.Getenv("JWT_DOMAIN")
    }
    //expire    
    if obj.Expire < 1 {
        obj.Expire, _ = strconv.Atoi(os.Getenv("JWT_EXPIRE"))
    }
    if obj.Expire < 1 {
        obj.Expire = 60 * 60 * 24  // 1 day   
    }
}

func (obj *Jwt) Ware(c *gin.Context) {
    // get token
    if token, err := c.Cookie(JWT_COOKIE_KEY); err == nil && token != "" {
        if json, err := obj.validate(token); err == nil {
            c.Set(JWT_OBJ_KEY, json)
            c.Next()    
            return
        }                
    }
    c.JSON(401, gin.H{"message": "invalid token"})
    c.Abort()
    return
}


// custom claim
type MyClaims struct {
    Payload string `json:"payload"`
    jwt.StandardClaims
}

// encode payload
func (obj *Jwt) sign(payload string) (string, error){    
    // Create the Claims
    claims := MyClaims{
        payload,
        jwt.StandardClaims{
            ExpiresAt: time.Now().Unix() + int64(obj.Expire),
        },
    }
    // get token
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    ss, err := token.SignedString(obj.secretBytes)
    return ss, err
}


// decode token
func (obj *Jwt) validate(tokenString string) (string, error){  
    // parse
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        return obj.secretBytes, nil
    })
    // validate
    if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
        if payload, ok := claims["payload"].(string); ok {
            return payload, nil
        }else{
            return "", errors.New("payload mismatch string type")
        }
    } else {
        return "", err
    }
}






