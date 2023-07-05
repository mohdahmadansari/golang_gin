package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (ctrl *Controllers) Welcome(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"success": 1, "message": "Welcome to go API with mysql."})
}
