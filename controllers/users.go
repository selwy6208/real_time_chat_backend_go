package controllers

import (
	"net/http"
	"real-chat-backend/models"
	"real-chat-backend/utils/token"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetUSerDetails(c *gin.Context) {
	_, err := token.ExtractTokenID(c)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.GetQuery("userID")
	userUIntId, _ := strconv.ParseUint(userID, 10, 32)
	ui := uint(userUIntId)
	u, _err := models.GetUserByID(ui)

	if _err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": _err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success", "data": u})
}
func GetUsers(c *gin.Context) {

	_, err := token.ExtractTokenID(c)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u, err := models.GetUsers()

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success", "data": u})
}
