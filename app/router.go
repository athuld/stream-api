package app

import "streamapi/controllers"

func mapUrls() {
	router.GET("/", controllers.HomeHello)
	router.POST("/add", controllers.AddStreamData)
	router.GET("/get", controllers.SearchData)
	router.GET("/file", controllers.GetFileData)
}
