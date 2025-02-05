package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/NgTrNamKhanh/go_ecommerce/database"
	"github.com/NgTrNamKhanh/go_ecommerce/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	generate "github.com/NgTrNamKhanh/go_ecommerce/tokens"
)


var UserCollection *mongo.Collection = database.UserData(database.Client, "Users")
var ProductCollection *mongo.Collection = database.ProductData(database.Client, "Products")
var Validate = validator.New()
func HashPassword(password string) string {
	bytes , err := bcrypt.GenerateFromPassword([]byte(password), 14)

	if err != nil {
		log.Panic(err)
	}
	return string(bytes)
}

func VerifyPassword(userPassword string, givenPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(givenPassword), []byte(userPassword))
	valid := true
	msg := ""
	if err!=nil{
		msg = "Email or password not correct"
		valid = false
	}
	return valid, msg
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

		if validationErr != nil{
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr})
			return
		}

		count, err := UserCollection.CountDocuments(ctx, bson.M{"email":user.Email})
		if err != nil{
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		}

		if count> 0{
			c.JSON(http.StatusBadRequest, gin.H{"error": "user already exist"})
		}

		count, err = UserCollection.CountDocuments(ctx, bson.M{"phone":user.Phone})

		defer cancel()
		if err != nil{
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return 
		}

		if count >0 {
			c.JSON(http.StatusBadRequest, gin.H{"error":"This phone already in use"})
			return 
		}
		passsword := HashPassword(*user.Password)
		user.Password = &passsword

		user.Created_At , _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.Updated_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.ID = primitive.NewObjectID()
		user.User_ID =  user.ID.Hex() 
		token,refreshtoken, _ := generate.TokenGenerator(*user.Email, *user.First_Name, *user.Last_Name, user.User_ID)
		user.Token = &token
		user.Refresh_Token = &refreshtoken
		user.UserCart = make([]models.ProductUser, 0)
		user.Address_Details = make([]models.Address, 0)
		user.Order_Status = make([]models.Order, 0)
		
		_, inserterr := UserCollection.InsertOne(ctx, user)
		if inserterr != nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error": "the user did not get created"})
			return
		}

		defer cancel()

		c.JSON(http.StatusCreated, "Successfully signed in ")
	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context){
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var founduser models.User
		var user models.User
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error":err})
			return
		}
		err := UserCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&founduser)
		defer cancel()

		if err!= nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error":"email or password incorrect"})
			return
		}

		PassswordIsValid, msg := VerifyPassword(*user.Password, *founduser.Password)

		defer cancel()

		if !PassswordIsValid {
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			fmt.Println(msg)
			return
		}
		token, refreshToken, err := generate.TokenGenerator(*founduser.Email, *founduser.First_Name,*founduser.Last_Name,founduser.User_ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error":"something went wrong"})
			return
		}
		defer cancel()
		err = generate.UpdateAllTokens(token, refreshToken, founduser.User_ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update tokens"})
			return
		}
		founduser.Token = &token
		founduser.Refresh_Token = &refreshToken
		c.JSON(http.StatusFound, founduser)
	}
}

