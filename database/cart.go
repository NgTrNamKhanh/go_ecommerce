package database

import "errors"

var (
	ErrCantFindProduct    = errors.New("")
	ErrCantDecodeProducts = errors.New("")
	ErrUserIdIsNotValid   = errors.New("")
	ErrCantUpdateUser     = errors.New("")
	ErrRemoveItemCart     = errors.New("")
	ErrCantGetItem        = errors.New("")
	ErrCantBuyCartItem    = errors.New("")
)

func AddProductToCart() {

}

func RemoveCartItem() {

}
func BuyItemFromCart() {

}

func InstantBuy() {

}
