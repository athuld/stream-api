package controllers

import (
	"net/http"
	"os"
	"strconv"
	"streamapi/domain"
	"streamapi/structs"
	"streamapi/utils/errors"
	"streamapi/utils/logger"

	"github.com/gin-gonic/gin"
)

func HomeHello(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Hello world"})
}

func AddStreamData(c *gin.Context) {
	refSecret := os.Getenv("REFERER_SECRET")
	reqRefSecret := c.Query("ref_secret")

	if refSecret == "" {
		c.JSON(http.StatusOK, gin.H{})
		return
	}

	if refSecret != reqRefSecret {
		c.JSON(http.StatusForbidden, gin.H{"message": "you are not authorised to do this request"})
		return
	}
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
	refSecret := os.Getenv("REFERER_SECRET")
	reqRefSecret := c.Query("ref_secret")

	if refSecret == "" {
		logger.Debug.Println("=== Query param missing ===")
		c.JSON(http.StatusOK, gin.H{})
		return
	}

	if refSecret != reqRefSecret {
		c.JSON(http.StatusForbidden, gin.H{"message": "you are not authorised to do this request"})
		return
	}
	var queryParams structs.QueryParams
	queryParams.Hash = c.Query("hash")
	queryParams.IPAddress = c.Query("ip_address")
	queryParams.Action = c.Query("action")

	data, err, cacheHeaders := domain.GetFileDataFromDB(c,queryParams)
	if err != nil {
		c.JSON(err.Status, err)
		return
	}
	
	if cacheHeaders.CacheHit {
		c.Header("Cache-Status", "HIT")
		c.Header("X-Cache-Key", cacheHeaders.CacheKey)
		c.JSON(http.StatusOK, data)
		return
	}

	c.Header("Cache-Status", "MISS")
	c.JSON(http.StatusOK, data)
}

func SearchData(c *gin.Context) {
	refSecret := os.Getenv("REFERER_SECRET")
	reqRefSecret := c.Query("ref_secret")

	if refSecret == "" {
		c.JSON(http.StatusOK, gin.H{})
		return
	}

	if refSecret != reqRefSecret {
		c.JSON(http.StatusForbidden, gin.H{"message": "you are not authorised to do this request"})
		return
	}
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

func GetRecentSearchData(c *gin.Context) {
	refSecret := os.Getenv("REFERER_SECRET")
	reqRefSecret := c.Query("ref_secret")

	if refSecret == "" {
		c.JSON(http.StatusOK, gin.H{})
		return
	}

	if refSecret != reqRefSecret {
		c.JSON(http.StatusForbidden, gin.H{"message": "you are not authorised to do this request"})
		return
	}

	var pageNumber, limitQuery int
	page := c.Query("page")
	limit := c.Query("limit")
	sort := c.Query("sort")
	var convErr error
	if page == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "'page' is a required query parameter"})
		return
	} else {
		pageNumber, convErr = strconv.Atoi(page)
	}

	if limit == "" {
		limitQuery = 10
	} else {
		limitQuery, convErr = strconv.Atoi(limit)
		if convErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "int conversion error"})
			return
		}
	}

	data, err := domain.GetRecentFilesFromDB(pageNumber, limitQuery, sort)

	if err != nil {
		c.JSON(err.Status, err)
		return
	}
	if len(*data) == 0 {
		c.JSON(http.StatusOK, []domain.Data{})
		return
	}

	c.JSON(http.StatusOK, data)
}

func DeleteData(c *gin.Context) {
	refSecret := os.Getenv("REFERER_SECRET")
	hash := c.Query("hash")
	reqRefSecret := c.Query("ref_secret")
	if refSecret != reqRefSecret {
		c.JSON(http.StatusForbidden, gin.H{"message": "you are not authorised to do this request"})
		return
	}
	err := domain.DeleteFileDataFromDB(hash)
	if err != nil {
		c.JSON(err.Status, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "data deleted from db"})
}
