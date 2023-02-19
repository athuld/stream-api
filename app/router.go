package app

import "streamapi/controllers"

func mapUrls(){
  router.GET("/",controllers.HomeHello)
}
