package middleware

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
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

func Authorization() gin.HandlerFunc {
	return func(c *gin.Context) {
		authorization := c.GetHeader(api.AuthorizationHeader)
		if authorization == "" {
			authorization = c.GetHeader(api.XAuthorizationHeader)
		}

		if authorization == "" {
			if c.Request.URL.Path == "/" {
				c.Header("Content-Type", "text/plain")
			} else if strings.HasSuffix(c.Request.URL.Path, "/login") || strings.HasPrefix(c.Request.URL.Path, "/chatgpt/public-api") {
				c.Header("Content-Type", "application/json")
			} else {
				c.AbortWithStatusJSON(http.StatusUnauthorized, api.ReturnMessage(emptyAccessTokenErrorMessage))
				return
			}

			c.Next()
		} else {
			if expired, email := isExpired(authorization); expired {
				c.AbortWithStatusJSON(http.StatusUnauthorized, api.ReturnMessage(fmt.Sprintf(accessTokenHasExpiredErrorMessage, email)))
				return
			}

			c.Set(api.AuthorizationHeader, authorization)
		}
	}
}

//goland:noinspection GoUnhandledErrorResult
func isExpired(accessToken string) (bool, string) {
	// accessToken
	split := strings.Split(accessToken, ".")
	if len(split) == 3 {
		rawDecodedText, _ := base64.RawStdEncoding.DecodeString(split[1])
		var accessToken AccessToken
		json.Unmarshal(rawDecodedText, &accessToken)

		exp := int64(accessToken.Exp)
		expTime := time.Unix(exp, 0)
		now := time.Now()
		if now.After(expTime) {
			return true, accessToken.HTTPSAPIOpenaiComProfile.Email
		}

		return false, ""
	}

	// apiKey
	return false, ""
}
