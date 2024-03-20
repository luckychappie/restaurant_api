package routes

import (
	"golang/restaurant_api/controllers"

	"github.com/gin-gonic/gin"
)

func MenuRoutes(r *gin.Engine) {
	r.GET("/menus", controllers.GetMenus())
	r.GET("/menus/:menu_id", controllers.GetMenu())
	r.POST("/menus", controllers.CreateMenu())
	r.PATCH("/menus/:menu_id", controllers.UpdateMenu())
}
