package model

import (
	"net/http"

	"github.com/dchest/captcha"
	"github.com/gin-gonic/gin"
)

func Getcaptchaimg(c *gin.Context) {
	captchaId := c.Param("captchaId")
	c.Header("Content-Type", "image/png")
	if err := captcha.WriteImage(c.Writer, captchaId, 240, 80); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate captcha image"})
	}
}

func Createcaptchaid(c *gin.Context) {
	captchaId := captcha.New()
	c.JSON(http.StatusOK, gin.H{
		"captchaId": captchaId,
	})
}

func Verifycaptcha(captchaId string, value string) bool {
	return captcha.VerifyString(captchaId, value)
}
