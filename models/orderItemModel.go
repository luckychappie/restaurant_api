package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrderItem struct {
	ID               primitive.ObjectID `bson:"_id"`
	Order_item_id    string             `json:"order_item_id"`
	Order_id         string             `json:"order_id"`
	Number_of_guests *int               `json:"number_of_guests" validate:"required"`
	Table_number     *int               `json:"table_number" validate:"required"`
	Created_at       time.Time          `json:"created_at"`
	Updated_at       time.Time          `json:"updated_at"`
	Table_id         string             `json:"table_id" validate:"required"`
	Unit_price       *float64           `josn:"unit_price"`
	Quantity         *int               `json:"quantity"`
	Food_id          *string            `json:"food_id"`
}
