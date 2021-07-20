package middleware

import (
	"errors"
	"flying-star/internal/casbin"
	"flying-star/internal/config"
	"flying-star/internal/module/user/rdb"
	"flying-star/utils"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strings"
)

var (
	ErrRoleCheckFail = errors.New("没有权限!")
)

const (
	TokenName = "token"
)

func CrossDomain() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Vary", "Origin")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Max-Age", "1728000")
		c.Header("Access-Control-Allow-Headers", "Accept,Origin,X-Requested-With,Content-Type,token,sign,app_id,timestamp")
		c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
		c.Header("Access-Control-Allow-Origin", "*")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
		} else {
			c.Next()
		}
	}
}

func CrossDomainForDebug() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Vary", "Origin")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Max-Age", "1728000")
		c.Header("Access-Control-Allow-Headers", "Accept,Origin,X-Requested-With,Content-Type,token,sign,app_id,timestamp")
		c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
		c.Header("Access-Control-Allow-Origin", "*")
		if c.Request.Method == http.MethodGet {
			log.Println("get ", c.Request.URL.RequestURI())
			c.Next()
		} else if c.Request.Method == http.MethodPost || c.Request.Method == http.MethodPut {
			raw, err := utils.CopyBody(c.Request)
			if err == nil {
				log.Println("----------")
				log.Println(string(raw))
				log.Println("----------")
			}
			c.Next()
		} else if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
		} else {
			log.Println(c.Request.Method, c.Request.RequestURI)
			c.Next()
		}
	}
}

// GetToken 获取token
func GetToken(c *gin.Context) string {
	if len(c.Query(TokenName)) > 0 {
		return c.Query(TokenName)
	} else if len(c.PostForm(TokenName)) > 0 {
		return c.PostForm(TokenName)
	} else if len(c.GetHeader(TokenName)) > 0 {
		return c.GetHeader(TokenName)
	} else {
		if t, err := c.Cookie(TokenName); err == nil && len(t) > 0 {
			return t
		}
	}
	return ""
}

// AdminAuth 管理员权限检查
func AdminAuth(secret string, withCasbin bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := GetToken(c)
		if len(token) > 0 {
			id, role, err := ParseAdminToken(token, secret)
			if err != nil || id == "" {
				utils.OutAuthOutdatedError(c)
				c.Abort()
				return
			}
			c.Set("uid", id)
			c.Set("role", role)
			c.Set("token", token)
			if withCasbin {
				// 通过casbin验证权限
				roles := strings.Split(role, ",")
				pass := false
				for _, s := range roles {
					if ok, _ := casbin.Casbin.Enforce(s, c.FullPath(), c.Request.Method); ok {
						pass = true
					}
				}
				if pass {
					c.Next()
					return
				} else {
					utils.OutRBACError(c)
					c.Abort()
					return
				}
			}
			c.Next()
			return
		}
		utils.OutAuthNeedError(c)
		c.Abort()
		return
	}
}

// DocAuth 文档权限检查
func DocAuth(prefix string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}

// GetLoginUid 获取登陆ID
func GetLoginUid(c *gin.Context) string {
	token := GetToken(c)
	userRdb := rdb.Info{}
	info, err := userRdb.Get(token, config.Info.Server.Token)
	if err != nil {
		return ""
	}
	return info.Uid
}

// GetLoginUserInfo 获取登陆用户信息
func GetLoginUserInfo(c *gin.Context) (info rdb.Info, err error) {
	token := GetToken(c)
	userRdb := rdb.Info{}
	return userRdb.Get(token, config.Info.Server.Token)
}
