package routes

import (
	"golang/restaurant_api/controllers"

	"github.com/gin-gonic/gin"
)

func TableRoutes(r *gin.Engine) {
	r.GET("/tables", controllers.GetTables())
	r.GET("/tables/:table_id", controllers.GetTable())
	r.POST("/tables", controllers.CreateTable())
	r.PATCH("/tables/:table_id", controllers.UpdateTable())
}
