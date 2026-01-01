package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"cortex/backend/invites"
)

func RedeemInvite(c *gin.Context) {
	code := c.Param("code")

	if !invites.ValidateInvite(code) {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "invalid or expired invite",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ok": true,
	})
}
