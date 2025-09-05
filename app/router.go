package app

import "streamapi/controllers"

func mapUrls() {
	router.GET("/", controllers.HomeHello)
	router.POST("/add", controllers.AddStreamData)
	router.GET("/get", controllers.SearchData)
	router.GET("/recent/search", controllers.GetRecentSearchData)
	router.GET("/file", controllers.GetFileData)
	router.DELETE("/delete", controllers.DeleteData)
}
