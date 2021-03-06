package middleware

import (
	"CourseSelectionSystem/db"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

var JwtKey = []byte("CourseSelectionSystem")

type MyClaims struct {
	Num  string `json:"num,omitempty"`
	Role int    `json:"role,omitempty"`
	jwt.StandardClaims
}

//SetToken 生成token
func SetToken(num string, role int, expireTime time.Time) (string, error) {
	SetClaims := MyClaims{
		Num:  num,
		Role: role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    "CourseSelectionSystem",
		},
	}
	reqClaim := jwt.NewWithClaims(jwt.SigningMethodHS256, SetClaims)
	token, err := reqClaim.SignedString(JwtKey)
	if err != nil {
		return "", err
	}
	return token, nil
}

// CheckToken 验证token
func CheckToken(token string) (*MyClaims, error) {
	setToken, _ := jwt.ParseWithClaims(token, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return JwtKey, nil
	})
	if key, _ := setToken.Claims.(*MyClaims); setToken.Valid {
		return key, nil
	}
	return nil, errors.New("token过期，验证失败！")
}

// TokenFmtCheck token格式校验
func TokenFmtCheck(c *gin.Context) (*MyClaims, error) {
	tokenHeader := c.Request.Header.Get("Authorization")
	if tokenHeader == "" {
		c.JSON(http.StatusOK, gin.H{
			"code":    -1,
			"stata":   "false",
			"message": "token不存在",
		})
		return nil, errors.New("token不存在")
	}
	checkToken := strings.SplitN(tokenHeader, " ", 2)
	if len(checkToken) != 2 && checkToken[0] != "Bearer" {
		c.JSON(http.StatusOK, gin.H{
			"code":    -1,
			"stata":   "false",
			"message": "token格式错误",
		})
		return nil, errors.New("token格式错误")
	}
	key, err := CheckToken(checkToken[1])
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    -1,
			"stata":   "false",
			"message": err.Error(),
		})
		return key, err
	}
	if time.Now().Unix() > key.ExpiresAt {
		c.JSON(http.StatusOK, gin.H{
			"code":    -1,
			"stata":   "false",
			"message": "token过期",
		})
		return key, errors.New("token过期")
	}
	return key, nil
}

// StuJwt 学生jwt中间件
func StuJwt() gin.HandlerFunc {
	return func(c *gin.Context) {
		key, err := TokenFmtCheck(c)
		if err != nil {
			c.Abort()
			return
		}
		stuNum := c.PostForm("StuNum")
		if stuNum != key.Num {
			c.JSON(http.StatusOK, gin.H{
				"code":    -1,
				"status":  "false",
				"message": "权限不足",
			})
			c.Abort()
			return
		}
		v, dberr := db.Rdb.HGET("users", key.Num)
		if dberr != nil || v != "1" {
			c.JSON(http.StatusOK, gin.H{
				"code":    -1,
				"status":  "false",
				"message": "请重新登录",
			})
			c.Abort()
			return
		}
		if key.Role != 1 {
			c.JSON(http.StatusOK, gin.H{
				"code":    -1,
				"stata":   "false",
				"message": "权限不足",
			})
			c.Abort()
			return
		}
		c.Set("Num", key.Num)
		c.Next()
	}
}

//TechJwt 教师jwt中间件
func TechJwt() gin.HandlerFunc {
	return func(c *gin.Context) {
		key, err := TokenFmtCheck(c)
		if err != nil {
			c.Abort()
			return
		}
		_, dberr := db.Rdb.HGET("users", key.Num)
		if dberr != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":    -1,
				"status":  "false",
				"message": "请重新登录",
			})
			c.Abort()
			return
		}
		if key.Role != 2 {
			c.JSON(http.StatusOK, gin.H{
				"code":    -1,
				"stata":   "false",
				"message": "权限不足",
			})
			c.Abort()
			return
		}
		c.Set("Num", key.Num)
		c.Next()
	}
}

//AdminJwt 管理员jwt中间件
func AdminJwt() gin.HandlerFunc {
	return func(c *gin.Context) {
		key, err := TokenFmtCheck(c)
		if err != nil {
			c.Abort()
			return
		}
		_, dberr := db.Rdb.HGET("users", key.Num)
		if dberr != nil {
			fmt.Println(dberr)
			c.JSON(http.StatusOK, gin.H{
				"code":    -1,
				"status":  "false",
				"message": "请重新登录",
			})
			c.Abort()
			return
		}
		if key.Role != 3 {
			c.JSON(http.StatusOK, gin.H{
				"code":    -1,
				"stata":   "false",
				"message": "权限不足",
			})
			c.Abort()
			return
		}
		c.Set("Num", key.Num)
		c.Next()
	}
}
