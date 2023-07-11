package token

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/linweiyuan/funcaptcha"
	"github.com/linweiyuan/go-chatgpt-api/api"

	http "github.com/bogdanfinn/fhttp"
)

var (
	bx string
)

//goland:noinspection SpellCheckingInspection
func init() {
	bx = os.Getenv("GO_CHATGPT_API_BX")
}

func GetArkoseToken(c *gin.Context) {
	token, err := funcaptcha.GetOpenAITokenWithBx(bx)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, api.ReturnMessage("Failed to get arkose token."))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}
