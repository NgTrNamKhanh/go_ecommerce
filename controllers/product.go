package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/NgTrNamKhanh/go_ecommerce/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AddProductAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var products models.Product
		defer cancel()
		if err := c.BindJSON(&products); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		products.Product_ID = primitive.NewObjectID()
		_, anyerr := ProductCollection.InsertOne(ctx, products)
		if anyerr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "not inserted"})
			return
		}
		defer cancel()
		c.JSON(http.StatusOK, "Sucessfully added")

	}
}
// func EditProductAdmin() gin.HandlerFunc {
// 	return func(c *gin.Context){
// 		// product_id := c.Query("id")
// 		// if product_id ==""{
// 		// 	c.Header("Content-Type","application/json")
// 		// 	c.JSON(http.StatusNotFound, gin.H{"error": "invalid"})
// 		// 	c.Abort()
// 		// 	return
// 		// }

// 		// productt_id, err := primitive.ObjectIDFromHex(product_id)

// 		// if err != nil {
// 		// 	c.IndentedJSON(500, "internal server error")
// 		// }

// 		var editaddress models.Product

// 		if err := c.BindJSON(&editaddress); err !=nil{
// 			c.IndentedJSON(http.StatusBadRequest, err.Error())
// 		}
// 		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
// 		defer cancel()

// 		filter := bson.D{primitive.E{Key: "_id", Value: productt_id}}

// 		update:= bson.D{{Key: "$set", Value: bson.D{primitive.E{Key: "address.0.house_name", Value: editaddress.House},{Key: "address.0.street_name", Value: editaddress.Street}, {Key: "address.0.city_name", Value: editaddress.City}, {Key: "address.0.pin_code",Value: editaddress.Pincode}}}}

// 		_,err = UserCollection.UpdateOne(ctx, filter, update)

// 		if err != nil {
// 			c.IndentedJSON(500, "something went wrong")
// 			return 
// 		}
// 		defer cancel()
// 		ctx.Done()
// 		c.IndentedJSON(200, "successfully updated")

// 	}
// }

func GetAllProducts() gin.HandlerFunc {
	return func(c *gin.Context) {
		var productlist []models.Product
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		cursor, err := ProductCollection.Find(ctx, bson.D{{}})
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, "something went wrong")
		}

		cursor.All(ctx, &productlist)

		if err != nil {
			log.Panicln(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		defer cursor.Close(ctx)

		if err := cursor.Err(); err != nil {
			log.Println(err)
			c.IndentedJSON(400, "invalid")

		}
		defer cancel()
		c.IndentedJSON(200, productlist)
	}
}

func SearchProductByQuery() gin.HandlerFunc {
	return func(c *gin.Context) {
		var searchProducts []models.Product
		queryParam := c.Query("name")

		//check if empty

		if queryParam == "" {
			log.Println("query is empty")
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"error": "Invalid search index"})
			c.Abort()
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		searchquerydb, err := ProductCollection.Find(ctx, bson.M{"product_name": bson.M{"$regex": queryParam}})

		if err != nil {
			c.IndentedJSON(404, "something went wrong")
			return
		}

		err = searchquerydb.All(ctx, &searchProducts)
		if err != nil {
			log.Println(err)
			c.IndentedJSON(400, "invalid")
			return
		}

		defer searchquerydb.Close(ctx)

		if err := searchquerydb.Err(); err != nil {
			log.Println(err)
			c.IndentedJSON(400, "something went wrong")
			return
		}

		defer cancel()

		c.IndentedJSON(200, searchProducts)

	}
}
