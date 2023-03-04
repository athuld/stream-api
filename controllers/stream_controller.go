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

func SearchData(c *gin.Context) {
	query := c.Query("query")

	data, err := domain.SearchDataFromDB(query)
	if err != nil {
		c.JSON(err.Status, err)
		return
	}
	c.JSON(http.StatusOK, data)

}
