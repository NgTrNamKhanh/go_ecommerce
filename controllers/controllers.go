package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/NgTrNamKhanh/Go-E-Commerce/models"
	"github.com/gin-gonic/gin"
)

func HashPassword(password string) string {

}

func VerifyPassword(userPassword string, givenPassword string) (bool, string) {

}

func SignUp() gin.HandlerFunc {
	

	return func(c *gin.Context){
		var ctx,cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var user models.User
		
		if err := c.BindJSON(&user); err != nil{
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		validationErr := Validate.Struct(user)

		if validationErr := nil{
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr})
			return
		}

		count, err := UserCollection.CountDocuments(ctx, bson.M{"email":user.Email})
		if err != nil{
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		}

		if count> 0{
			c.JSON(http.StatusBadRequest, gin.H{"error: user already exist" })
		}

		count, err := UserCollection.CountDocuments(ctx, bson.M{"phone":user.Phone})

		defer cancel()
		if err != nil{
			log.Panic(err)
			c.JSON(htpt.StatusInternalServerError, gin.H{"error": err})
			return 
		}

		if count >0 {
			c.JSON(http.StatusBadRequest, gin.H{"error":"This phone already in use"})
			return 
		}
		
	}
}

func Login() gin.handleFunc {

}

func ProductViewerAdmin() gin.handleFunc {

}

func searchProduct() gin.handleFunc {

}

func searchProductByQuery() gin.handleFunc {

}