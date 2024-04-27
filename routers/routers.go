package routes

import (
	"github.com/gin-gonic/gin"
	"keyvaluestore/controller"
)

func SecretRoute(router *gin.Engine) {
	router.POST("/keys", controller.CreateSecret)
	router.GET("/keys", controller.GetSecrets)
	router.GET("/keys/:key", controller.GetSecret)
	router.DELETE("/keys/:key", controller.DelSecret)
	router.PUT("/keys/:key", controller.PutSecret)
}
