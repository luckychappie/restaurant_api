package routes

import (
	"golang/restaurant_api/controllers"

	"github.com/gin-gonic/gin"
)

func InvoiceRoutes(r *gin.Engine) {
	r.GET("/invoices", controllers.GetInvoices())
	r.GET("/invoices/:invoice_id", controllers.GetInvoice())
	r.POST("/invoices", controllers.CreateInvoice())
	r.PATCH("/invoices/:invoice_id", controllers.UpdateInvoice())
}
