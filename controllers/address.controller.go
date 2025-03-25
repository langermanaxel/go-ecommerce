package controllers

import (
	"context"
	"fmt"
	"go-ecommerce/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func AddAddress() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user_id := ctx.Query("id")
		if user_id == "" {
			ctx.Header("Content-Type", "application/json")
			ctx.JSON(http.StatusNotFound, gin.H{"Error": "Invalid search index"})
			ctx.Abort()
			return
		}

		address, err := primitive.ObjectIDFromHex(user_id)
		if err != nil {
			ctx.IndentedJSON(http.StatusInternalServerError, err)
			return
		}

		var addresses models.Address
		addresses.Address_Id = primitive.NewObjectID()
		if err := ctx.BindJSON(&addresses); err != nil {
			ctx.IndentedJSON(http.StatusNotAcceptable, err.Error())
			return
		}

		var context, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		match_filter := bson.D{{Key: "$match", Value: bson.D{primitive.E{Key: "_id", Value: address}}}}
		unwind := bson.D{{Key: "$unwind", Value: bson.D{primitive.E{Key: "path", Value: "$address"}}}}
		grouping := bson.D{{Key: "$group", Value: bson.D{primitive.E{Key: "_id", Value: "$address_id"}, {Key: "count", Value: bson.D{primitive.E{Key: "$sum", Value: 1}}}}}}

		pointCursor, err := user_collection.Aggregate(context, mongo.Pipeline{match_filter, unwind, grouping})
		if err != nil {
			ctx.IndentedJSON(http.StatusInternalServerError, err)
			return
		}

		var address_info []bson.M
		if err := pointCursor.All(context, &address_info); err != nil {
			panic(err)
		}

		var size int32
		for _, address_nr := range address_info {
			count := address_nr["count"]
			size = count.(int32)
		}

		if size < 2 {
			filter := bson.D{primitive.E{Key: "_id", Value: address}}
			update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "address", Value: addresses}}}}
			_, err := user_collection.UpdateOne(context, filter, update)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			ctx.IndentedJSON(http.StatusBadRequest, "Not allowed")
		}

		context.Done()
	}
}

func EditHomeAddress() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user_id := ctx.Query("id")
		if user_id == "" {
			ctx.Header("Content-Type", "application/json")
			ctx.JSON(http.StatusNotFound, gin.H{"Error": "Invalid search index"})
			ctx.Abort()
			return
		}

		user_index, err := primitive.ObjectIDFromHex(user_id)
		if err != nil {
			ctx.IndentedJSON(http.StatusInternalServerError, err)
			return
		}

		var editaddress models.Address
		if err := ctx.BindJSON(&editaddress); err != nil {
			ctx.IndentedJSON(http.StatusBadRequest, err.Error())
			return
		}

		var context, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		filter := bson.D{primitive.E{Key: "_id", Value: user_index}}
		update := bson.D{{Key: "$set", Value: bson.D{primitive.E{Key: "address.0.house_name", Value: editaddress.House}, {Key: "address.0.street_name", Value: editaddress.Street}, {Key: "address.0.city_name", Value: editaddress.City}, {Key: "address.0.pin_code", Value: editaddress.Pincode}}}}

		_, err = user_collection.UpdateOne(context, filter, update)
		if err != nil {
			ctx.IndentedJSON(http.StatusInternalServerError, err)
			return
		}

		context.Done()
		ctx.IndentedJSON(http.StatusOK, "Successfully Updated the Home address")
	}
}

func EditWorkAddress() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user_id := ctx.Query("id")
		if user_id == "" {
			ctx.Header("Content-Type", "application/json")
			ctx.JSON(http.StatusNotFound, gin.H{"Error": "Invalid search index"})
			ctx.Abort()
			return
		}

		user_index, err := primitive.ObjectIDFromHex(user_id)
		if err != nil {
			ctx.IndentedJSON(http.StatusInternalServerError, err)
			return
		}

		var editaddress models.Address
		if err := ctx.BindJSON(&editaddress); err != nil {
			ctx.IndentedJSON(http.StatusBadRequest, err.Error())
			return
		}

		var context, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		filter := bson.D{primitive.E{Key: "_id", Value: user_index}}
		update := bson.D{{Key: "$set", Value: bson.D{primitive.E{Key: "address.1.house_name", Value: editaddress.House}, {Key: "address.1.street_name", Value: editaddress.Street}, {Key: "address.1.city_name", Value: editaddress.City}, {Key: "address.1.pin_code", Value: editaddress.Pincode}}}}

		_, err = user_collection.UpdateOne(context, filter, update)
		if err != nil {
			ctx.IndentedJSON(http.StatusInternalServerError, err)
			return
		}

		context.Done()
		ctx.IndentedJSON(http.StatusOK, "Successfully Updated the Work address")
	}
}

func DeleteAddress() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user_id := ctx.Query("id")
		if user_id == "" {
			ctx.Header("Content-Type", "application/json")
			ctx.JSON(http.StatusNotFound, gin.H{"Error": "Invalid search index"})
			ctx.Abort()
			return
		}

		adresses := make([]models.Address, 0)
		user_index, err := primitive.ObjectIDFromHex(user_id)
		if err != nil {
			ctx.IndentedJSON(http.StatusInternalServerError, err)
			return
		}

		var context, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		filter := bson.D{primitive.E{Key: "_id", Value: user_index}}
		update := bson.D{{Key: "$set", Value: bson.D{primitive.E{Key: "address", Value: adresses}}}}

		_, err = user_collection.UpdateOne(context, filter, update)
		if err != nil {
			ctx.IndentedJSON(http.StatusNotFound, err)
			return
		}

		context.Done()

		ctx.IndentedJSON(http.StatusOK, "Successfully deleted")
	}
}
