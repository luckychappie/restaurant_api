package controllers

import (
	"context"
	"golang/restaurant_api/database"
	"golang/restaurant_api/models"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/bluesuncorp/validator.v5"
)

var fooCollection *mongo.Collection = database.OpenCollection(database.Client, "food")
var validate = validator.New("validate", validator.BakedInValidators)

func GetFoods() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		recordPerPage, err := strconv.Atoi(c.Query("recordPerPage"))
		if err != nil || recordPerPage < 1 {
			recordPerPage = 10
		}

		page, err := strconv.Atoi(c.Query("page"))
		if err != nil || page < 1 {
			page = 1
		}

		startIndex := (page - 1) * recordPerPage

		matchStage := bson.D{{"$match", bson.D{{}}}}
		groupStage := bson.D{{"$group", bson.D{{"_id", bson.D{{"_id", "null"}}}, {"total_count", bson.D{{"$sum", "1"}}}, {"data", bson.D{{"$push", "$$ROOT"}}}}}}
		projectStage := bson.D{
			{
				"$project", bson.D{
					{"_id", 0},
					{"total_count", 1},
					{"food_items", bson.D{{"$slice", []interface{}{"$data", startIndex, recordPerPage}}}}}}}

		result, err := fooCollection.Aggregate(ctx, mongo.Pipeline{
			matchStage, groupStage, projectStage})

		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}

		var allFood []bson.M
		if err = result.All(ctx, &allFood); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		c.JSON(http.StatusOK, allFood)
	}
}

func GetFood() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		foodId := c.Param("food_id")
		var food *models.Food

		err := fooCollection.FindOne(ctx, bson.M{"food_id": foodId}).Decode(&food)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}

		c.JSON(http.StatusOK, food)
	}
}

func CreateFood() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		var food *models.Food
		var menu *models.Menu

		if err := c.ShouldBind(&food); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}

		if err := validate.Struct(food); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}

		err := menuCollection.FindOne(ctx, bson.M{"menu_id": food.Menu_id}).Decode(&menu)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Menu not found"})
			return
		}

		food.ID = primitive.NewObjectID()
		food.Food_id = food.ID.Hex()
		var price = ToFixed(*food.Price, 2)
		food.Price = &price
		food.Crated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		food.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		result, err := fooCollection.InsertOne(ctx, food)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		defer cancel()
		c.JSON(http.StatusOK, gin.H{"data": result})

	}
}

func UpdateFood() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		foodId := c.Param("food_id")
		var food *models.Food
		var menu *models.Menu
		var updateObject primitive.D

		if err := c.BindJSON(&food); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}

		err := validate.Struct(food)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}

		if food.Name != nil {
			updateObject = append(updateObject, bson.E{"name", food.Name})
		}

		if food.Price != nil {
			updateObject = append(updateObject, bson.E{"price", food.Price})
		}
		if food.Food_image != nil {
			updateObject = append(updateObject, bson.E{"food_image", food.Food_image})
		}
		if *food.Menu_id != "" {
			err := menuCollection.FindOne(ctx, bson.M{"menu_id": food.Menu_id}).Decode(&menu)
			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "Menu not found"})
				return
			}
			updateObject = append(updateObject, bson.E{"menu_id", food.Menu_id})
		}

		food.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		updateObject = append(updateObject, bson.E{"updated_at", food.Updated_at})

		upsert := true
		opt := options.UpdateOptions{
			Upsert: &upsert,
		}

		result, updateErr := fooCollection.UpdateOne(ctx, bson.M{"food_id": foodId}, bson.M{"$set": updateObject}, &opt)
		defer cancel()
		if updateErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": updateErr})
			return
		}

		c.JSON(http.StatusOK, result)
	}
}

func Round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func ToFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(Round(num*output)) / output
}
