package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func HomeHello(c *gin.Context){
  c.JSON(http.StatusOK,gin.H{"message":"Hello world"})
}
