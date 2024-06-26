package models

import "time"

type User struct {
	ID              primitive.ObjectID `json: "_id" bson: "_id`
	First_Name      *string            `json: "first_name" validate:"required, min=2, max=30"` 
	Last_Name       *string            `json: "last" validate:"required, min=2, max= 30"`
	Password        *string            `json: "password" validate:"required, min=6"`
	Email           *string            `json: "email" validate:"email, required"`
	Phone           *string            `json: "phone" validate:"required"`
	Token           *string            `json: "token" `
	Refresh_Token   *string            `json: "refresh_token"`
	Created_At      time.Time          `json: "create_at"`
	Updated_At      time.Time          `json: "updated_at"`
	User_ID         string             `json: "user_id"`
	UserCart        []ProductUser      `json: "usercart" bson:"usercart"`
	Address_Details []Address          `json: "address" bson:"address"`
	Order_Status    []Order            `json: "orders" bson:"orders"`
}

type Product struct {
	Product_ID   primitive.ObjectID `bson: "_id`
	Product_Name *string            `json: "product_name"`
	Price        *uint              `json: "price"`
	Rating       *uint              `json: "email"`
	Image        *string            `json: "image"`
}

type ProductUser struct {
	Product_ID   primitive.ObjectID `bson: "_id"`
	Product_Name *string            `json: "product_name" bson: "product_name"`
	Price        *uint              `json: "price" bson: "product_name`
	Rating       *uint              `json: "rating" bson: "rating`
	Image        *string            `json: "image" bson: "image`
}

type Address struct {
	Address_ID primitive.ObjectID `bson: "_id"`
	House      *string `json: "house_name" bson: "house_name"`
	Street     *string`json: "street" bson: "street"`
	City       *string`json: "city" bson: "city"`
	Pincode    *string`json: "pin_code" bson: "pin_code"`
}

type Order struct {
	Order_ID       primitive.ObjectID `bson: "_id"`
	Order_Cart     []ProductUser  `json: "order_list" bson: "order_list"`
	Ordered_At     time.Time `json: "ordered_at" bson: "ordered_at"`
	Price          int `json: "price" bson: "price"`
	Discount       *int `json: "discount" bson: "discount"`
	Payment_Method Payment `json: "payment_method" bson: "payment_method"`
}
type Payment struct {
	Digital bool
	COD     bool
}
