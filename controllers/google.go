package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/NgTrNamKhanh/go_ecommerce/models"
	"github.com/gin-gonic/gin"
	"github.com/markbates/goth/gothic"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func StartGoogleAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		gothic.BeginAuthHandler(c.Writer, c.Request)
	}
}

func CompleteGoogleAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		googleuser, err := gothic.CompleteUserAuth(c.Writer, c.Request)
		if err != nil {
			log.Printf("Google OAuth2 callback failed: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to complete auth"})
			return
		}
		var user models.User
		user.Email = &googleuser.Email
		user.First_Name = &googleuser.FirstName
		user.Last_Name = &googleuser.LastName
		user.Created_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.Updated_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.ID = primitive.NewObjectID()
		user.User_ID = user.ID.Hex()
		user.Token = &googleuser.AccessToken
		user.Refresh_Token = &googleuser.RefreshToken
		user.UserCart = make([]models.ProductUser, 0)
		user.Address_Details = make([]models.Address, 0)
		user.Order_Status = make([]models.Order, 0)

		_, inserterr := UserCollection.InsertOne(ctx, user)
		if inserterr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "the user did not get created"})
			return
		}

		defer cancel()

		c.JSON(http.StatusCreated, "Successfully signed in ")
		// Handle user authentication or registration with user details
		c.JSON(http.StatusOK, gin.H{
			"message": "Google login successful",
			"user":    user,
		})

	}
}
