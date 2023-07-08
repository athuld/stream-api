package controllers

import (
	"net/http"
	"streamapi/domain"
	"streamapi/utils/errors"

	"github.com/gin-gonic/gin"
)

func HomeHello(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Hello world"})
}

func AddStreamData(c *gin.Context) {
	var data domain.Data
	if err := c.ShouldBindJSON(&data); err != nil {
		err := errors.NewBadRequestError("Invalid json")
		c.JSON(err.Status, err)
		return
	}

	if err := data.AddDataToDB(); err != nil {
		c.JSON(err.Status, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Data inserted properly"})

}

func GetFileData(c *gin.Context) {
	hash := c.Query("hash")
	ipAddress := c.Query("ip_address")
	action := c.Query("action")

	data, err := domain.GetFileDataFromDB(hash,ipAddress,action)
	if err != nil {
		c.JSON(err.Status, err)
		return
	}
	c.JSON(http.StatusOK, data)
}

func SearchData(c *gin.Context) {
	query := c.Query("query")

	if query == "" {
		c.JSON(http.StatusOK, gin.H{})
		return
	}

	data, err := domain.SearchDataFromDB(query)

	if err != nil {
		c.JSON(err.Status, err)
		return
	}
	c.JSON(http.StatusOK, data)

}
