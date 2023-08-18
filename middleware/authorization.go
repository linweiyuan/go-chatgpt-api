package middleware

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/linweiyuan/go-chatgpt-api/api"
)

const (
	emptyAccessTokenErrorMessage      = "Please provide a valid access token or api key in 'Authorization' header."
	accessTokenHasExpiredErrorMessage = "The accessToken for account %s has expired."
)

type AccessToken struct {
	HTTPSAPIOpenaiComProfile struct {
		Email         string `json:"email"`
		EmailVerified bool   `json:"email_verified"`
	} `json:"https://api.openai.com/profile"`
	HTTPSAPIOpenaiComAuth struct {
		UserID string `json:"user_id"`
	} `json:"https://api.openai.com/auth"`
	Iss   string   `json:"iss"`
	Sub   string   `json:"sub"`
	Aud   []string `json:"aud"`
	Iat   int      `json:"iat"`
	Exp   int      `json:"exp"`
	Azp   string   `json:"azp"`
	Scope string   `json:"scope"`
}

//goland:noinspection SpellCheckingInspection
func Authorization() gin.HandlerFunc {
	return func(c *gin.Context) {
		authorization := c.GetHeader(api.AuthorizationHeader)
		if authorization == "" {
			authorization = c.GetHeader(api.XAuthorizationHeader)
		}

		if authorization == "" {
			if c.Request.URL.Path == "/" {
				c.Header("Content-Type", "text/plain")
			} else if strings.HasSuffix(c.Request.URL.Path, "/login") ||
				strings.HasPrefix(c.Request.URL.Path, "/chatgpt/public-api") ||
				(strings.HasPrefix(c.Request.URL.Path, "/imitate") && os.Getenv("IMITATE_ACCESS_TOKEN") != "") {
				c.Header("Content-Type", "application/json")
			} else if c.Request.URL.Path == "/favicon.ico" {
				c.Abort()
				return
			} else {
				c.AbortWithStatusJSON(http.StatusUnauthorized, api.ReturnMessage(emptyAccessTokenErrorMessage))
				return
			}

			c.Next()
		} else {
			if expired := isExpired(c); expired {
				c.AbortWithStatusJSON(http.StatusUnauthorized, api.ReturnMessage(fmt.Sprintf(accessTokenHasExpiredErrorMessage, c.GetString(api.EmailKey))))
				return
			}

			c.Set(api.AuthorizationHeader, authorization)
		}
	}
}

//goland:noinspection GoUnhandledErrorResult
func isExpired(c *gin.Context) bool {
	accessToken := c.GetHeader(api.AuthorizationHeader)
	split := strings.Split(accessToken, ".")
	if len(split) == 3 {
		rawDecodedText, _ := base64.RawStdEncoding.DecodeString(split[1])
		var accessToken AccessToken
		json.Unmarshal(rawDecodedText, &accessToken)

		c.Set(api.EmailKey, accessToken.HTTPSAPIOpenaiComProfile.Email)

		exp := int64(accessToken.Exp)
		expTime := time.Unix(exp, 0)
		now := time.Now()
		if now.After(expTime) {
			return true
		}

		return false
	}

	// apiKey
	return false
}
