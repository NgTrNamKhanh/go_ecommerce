package main

import (
	"log"
	"os"
	"github.com/NgTrNamKhanh/go_ecommerce/controllers"
	"github.com/NgTrNamKhanh/go_ecommerce/database"
	"github.com/NgTrNamKhanh/go_ecommerce/middleware"
	"github.com/gin-gonic/gin"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/google"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	app := controllers.NewApplication(database.ProductData(database.Client, "Products"), database.UserData(database.Client, "Users"))

	router := gin.New()
	router.Use(gin.Logger())

	goth.UseProviders(
		google.New(os.Getenv("GOOGLE_CLIENT_ID"), os.Getenv("GOOGLE_CLIENT_SECRET"), "http://localhost:8000/auth/google/callback"),
	)

	router.POST("/signup", controllers.SignUp())
	router.POST("/login", controllers.Login())

	router.GET("/auth/google/login", controllers.StartGoogleAuth())
	router.GET("/auth/google/callback", controllers.CompleteGoogleAuth())

	productPortal := router.Group("/product", middleware.Authentication())
	{
		productPortal.POST("/add", controllers.AddProductAdmin())
		productPortal.GET("/products", controllers.GetAllProducts())
		productPortal.GET("/search", controllers.SearchProductByQuery())
	}

	addressPortal := router.Group("/address", middleware.Authentication())
	{
		addressPortal.POST("/add", controllers.AddAdress())
		addressPortal.PUT("/edit/home/:id", controllers.EditHomeAddress())
		addressPortal.PUT("/edit/work/:id", controllers.EditWorkAddress())
		addressPortal.DELETE("/:id", controllers.DeleteAddress())
	}

	cartPortal := router.Group("/cart", middleware.Authentication())
	{
		cartPortal.POST("/addtocart", app.AddToCart())
		cartPortal.POST("/removeitem", app.RemoveItem())
		cartPortal.GET("/items", app.GetItemFromCart())
		cartPortal.POST("/cartcheckout", app.BuyFromCart())
		cartPortal.POST("/instantbuy", app.InstantBuy())
	}

	log.Fatal(router.Run(":" + port))

}
