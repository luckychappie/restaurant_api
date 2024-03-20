package controllers

import (
	"context"
	"golang/restaurant_api/database"
	"golang/restaurant_api/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var invoiceCollection = database.OpenCollection(database.Client, "invoice")

type InvoiceViewFormat struct {
	Invocie_id       string
	Payment_method   string
	Order_id         string
	Payment_status   *string
	Payment_due      interface{}
	Table_number     interface{}
	Payment_due_date time.Time
	Order_details    interface{}
}

func GetInvoices() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		var invoices []models.Invoice

		cursor, err := invoiceCollection.Find(context.TODO(), bson.M{})
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}

		if err = cursor.All(ctx, &invoices); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}

		c.JSON(http.StatusOK, invoices)
	}
}

func GetInvoice() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		invoiceId := c.Param("invoice_id")
		var invoice *models.Invoice

		err := invoiceCollection.FindOne(ctx, bson.M{"invoice_id": invoiceId}).Decode(&invoice)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}

		var invoiceView InvoiceViewFormat

		// allOrderItem, err := ItemsByOrder(invoice.Order_id)
		invoiceView.Order_id = invoice.Order_id
		invoiceView.Payment_due_date = invoice.Payment_due_date
		invoiceView.Payment_method = "null"

		if invoice.Payment_method != nil {
			invoiceView.Payment_method = *invoice.Payment_method
		}
		invoiceView.Invocie_id = invoice.Invoice_id
		invoiceView.Payment_status = *&invoice.Payment_status
		// invoiceView.Payment_due = allOrderItem[0]["payment_due"]
		// invoiceView.Table_number = allOrderItem[0]["table_number"]
		// invoiceView.Order_details = allOrderItem[0]["order_items"]

		c.JSON(http.StatusOK, invoiceView)
	}
}

func UpdateInvoice() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		invoiceId := c.Param("invoiceId")
		var invoice models.Invoice
		var updateObject primitive.D

		if err := c.BindJSON(invoice); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}

		invoice.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		updateObject = append(updateObject, bson.E{"updated_at", invoice.Updated_at})
		upsert := true
		opt := options.UpdateOptions{
			Upsert: &upsert,
		}
		paymentStatus := "PENDING"
		if invoice.Payment_status == nil {
			invoice.Payment_status = &paymentStatus
		}

		result, err := invoiceCollection.UpdateOne(ctx, bson.M{"invoice_id": invoiceId}, bson.M{"$set": updateObject}, &opt)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		c.JSON(http.StatusOK, result)
	}
}

func CreateInvoice() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		var invoice models.Invoice
		var order models.Order

		if err := c.BindJSON(&invoice); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}

		err := orderCollection.FindOne(ctx, bson.M{"ordre_id": invoice.Order_id}).Decode(&order)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Order not found"})
			return
		}

		invoice.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		invoice.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		invoice.Payment_due_date, _ = time.Parse(time.RFC3339, time.Now().AddDate(0, 0, 1).Format(time.RFC3339))
		invoice.ID = primitive.NewObjectID()
		invoice.Invoice_id = invoice.ID.Hex()

		paymentStatus := "PENDING"
		if invoice.Payment_status == nil {
			invoice.Payment_status = &paymentStatus
		}

		if err = validate.Struct(invoice); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}

		result, err := invoiceCollection.InsertOne(ctx, invoice)
		defer cancel()
		if err = validate.Struct(invoice); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}
		c.JSON(http.StatusOK, result)
	}
}
