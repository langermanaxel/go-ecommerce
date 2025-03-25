package controllers

import (
	"context"
	"fmt"
	"go-ecommerce/database"
	"go-ecommerce/models"
	"go-ecommerce/tokens"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var user_collection *mongo.Collection = database.UserData(database.Client, "Users")
var product_collection *mongo.Collection = database.ProductData(database.Client, "Products")
var Validate = validator.New()

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}
	return string(bytes)
}

func VerifyPassword(user_password, given_password string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(given_password), []byte(user_password))
	valid := true
	msg := ""
	if err != nil {
		msg = "login or password is incorrect"
		valid = false
	}
	return valid, msg
}

func Signup() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var context, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var user models.User
		if err := ctx.BindJSON(&user); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := Validate.Struct(user)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}

		count, err := user_collection.CountDocuments(context, bson.M{"email": user.Email})
		if err != nil {
			log.Panic(err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}

		if count > 0 {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "user alredy exists"})
			return
		}

		count, err = user_collection.CountDocuments(context, bson.M{"phone": user.Phone})
		if err != nil {
			log.Panic(err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}

		if count > 0 {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "user alredy exists"})
			return
		}

		password := HashPassword(*user.Password)
		user.Password = &password

		user.Created_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.Updated_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.ID = primitive.NewObjectID()
		user.User_Id = user.ID.Hex()

		token, refresh_token, _ := tokens.TokenGenerator(*user.Email, *user.First_Name, *user.Last_Name, user.User_Id)

		user.Token = &token
		user.Refresh_Token = &refresh_token
		user.User_Cart = make([]models.ProductUser, 0)
		user.Address_Details = make([]models.Address, 0)
		user.Order_Status = make([]models.Order, 0)

		_, err = user_collection.InsertOne(context, user)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "user not created"})
			return
		}

		ctx.JSON(http.StatusCreated, gin.H{"success": "sign in successfully!"})
	}
}

func Login() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var context, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var user models.User
		if err := ctx.BindJSON(&user); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var found_user models.User

		err := user_collection.FindOne(context, bson.M{"email": user.Email}).Decode(&found_user)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}

		passwordIsValid, msg := VerifyPassword(*user.Password, *found_user.Password)
		if !passwordIsValid {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			fmt.Println(msg)
			return
		}

		token, refresh_token, _ := tokens.TokenGenerator(*found_user.Email, *found_user.First_Name, *found_user.Last_Name, found_user.User_Id)
		tokens.UpdateAllTokens(token, refresh_token, found_user.User_Id)

		ctx.JSON(http.StatusFound, found_user)
	}
}

func ProductViewerAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var products models.Product

		if err := c.BindJSON(&products); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		products.Product_Id = primitive.NewObjectID()

		_, err := product_collection.InsertOne(ctx, products)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "not inserted"})
			return
		}

		c.JSON(http.StatusOK, "successfully added")
	}
}

func SearchProduct() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var productList []models.Product
		var context, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		cursor, err := product_collection.Find(context, bson.D{{}})
		if err != nil {
			ctx.IndentedJSON(http.StatusInternalServerError, "something went wrong, please try after some time")
			return
		}

		err = cursor.All(context, &productList)
		if err != nil {
			log.Println(err)
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		defer cursor.Close(context)

		if err := cursor.Err(); err != nil {
			log.Println(err)
			ctx.IndentedJSON(http.StatusBadRequest, "invalid")
			return
		}

		ctx.IndentedJSON(http.StatusOK, productList)
	}
}

func SearchProductByQuery() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var searchProducts []models.Product

		queryParam := ctx.Query("name")
		if queryParam == "" {
			log.Println("query is empty")
			ctx.Header("content-type", "application/json")
			ctx.JSON(http.StatusNotFound, gin.H{"Error": "Invalid Search Index"})
			ctx.Abort()
			return
		}

		var context, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		searchQueryDB, err := product_collection.Find(context, bson.M{"product_name": bson.M{"$regex": queryParam}})
		if err != nil {
			ctx.IndentedJSON(http.StatusNotFound, "something went wrong while fetching the data")
			return
		}

		err = searchQueryDB.All(context, &searchProducts)
		if err != nil {
			log.Println(err)
			ctx.IndentedJSON(http.StatusBadRequest, "invalid")
			return
		}

		defer searchQueryDB.Close(context)

		if err := searchQueryDB.Err(); err != nil {
			log.Println(err)
			ctx.IndentedJSON(http.StatusBadRequest, "invalid")
			return
		}

		ctx.IndentedJSON(http.StatusOK, searchProducts)
	}
}
