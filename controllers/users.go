package controllers

import (
	"net/http"
	"real-chat-backend/models"
	"real-chat-backend/utils/token"

	"github.com/gin-gonic/gin"
)

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
