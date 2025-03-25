package database

import (
	"context"
	"errors"
	"go-ecommerce/models"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrCantFindProduct    = errors.New("can't find the product")
	ErrCantDecodeProducts = errors.New("can't find the product")
	ErrUserIdIsNotValid   = errors.New("this user is not valid")
	ErrCantUpdateUser     = errors.New("cannot add this product to the cart")
	ErrCantRemoveItemCart = errors.New("cannot remove this item from the cart")
	ErrCantGetItem        = errors.New("was unable to get the item from the cart")
	ErrCantBuyCartItem    = errors.New("cannot update the purchase")
)

func AddProductToCart(
	ctx context.Context,
	product_collection, user_collection *mongo.Collection,
	product_id primitive.ObjectID,
	user_id string) error {
	searchFromDB, err := product_collection.Find(ctx, bson.M{"_id": product_id})
	if err != nil {
		log.Println(err)
		return ErrCantFindProduct
	}

	var productCart []models.ProductUser
	if err = searchFromDB.All(ctx, &productCart); err != nil {
		log.Println(err)
		return ErrCantDecodeProducts
	}

	id, err := primitive.ObjectIDFromHex(user_id)
	if err != nil {
		log.Println(err)
		return ErrUserIdIsNotValid
	}

	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "usercart", Value: bson.D{{Key: "$each", Value: productCart}}}}}}

	_, err = user_collection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Println(err)
		return ErrCantUpdateUser
	}

	return nil
}

func RemoveCartItem(
	ctx context.Context,
	product_collection, user_collection *mongo.Collection,
	product_id primitive.ObjectID,
	user_id string) error {
	id, err := primitive.ObjectIDFromHex(user_id)
	if err != nil {
		log.Println(err)
		return ErrUserIdIsNotValid
	}

	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	update := bson.M{"$pull": bson.M{"usercart": bson.M{"_id": product_id}}}

	_, err = user_collection.UpdateMany(ctx, filter, update)
	if err != nil {
		log.Println(err)
		return ErrCantRemoveItemCart
	}

	return nil
}

func BuyItemFromCart(
	ctx context.Context,
	user_collection *mongo.Collection,
	user_id string) error {
	id, err := primitive.ObjectIDFromHex(user_id)
	if err != nil {
		log.Println(err)
		return ErrUserIdIsNotValid
	}

	var getCartItems models.User
	var orderCart models.Order

	orderCart.Order_Id = primitive.NewObjectID()
	orderCart.Ordered_At = time.Now()
	orderCart.Order_Cart = make([]models.ProductUser, 0)
	orderCart.Payment_Method.COD = true

	unwind := bson.D{{Key: "$unwind", Value: bson.D{primitive.E{Key: "path", Value: "$usercart"}}}}
	grouping := bson.D{{Key: "$group", Value: bson.D{primitive.E{Key: "_id", Value: "$_id"}, {Key: "total", Value: bson.D{primitive.E{Key: "$sum", Value: "$usercart.price"}}}}}}

	currentResult, err := user_collection.Aggregate(ctx, mongo.Pipeline{unwind, grouping})
	ctx.Done()
	if err != nil {
		panic(err)
	}

	var getUserCart []bson.M
	if err = currentResult.All(ctx, &getUserCart); err != nil {
		panic(err)
	}

	var total_price int32
	for _, user_item := range getUserCart {
		price := user_item["total"]
		total_price = price.(int32)
	}

	price := uint64(total_price)
	orderCart.Price = &price

	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "orders", Value: orderCart}}}}

	_, err = user_collection.UpdateMany(ctx, filter, update)
	if err != nil {
		log.Println(err)
	}

	err = user_collection.FindOne(ctx, bson.D{primitive.E{Key: "_id", Value: id}}).Decode(&getCartItems)
	if err != nil {
		log.Println(err)
	}

	filter_two := bson.D{primitive.E{Key: "_id", Value: id}}
	update_two := bson.M{"$push": bson.M{"orders.$[].order_list": bson.M{"$each": getCartItems.User_Cart}}}

	_, err = user_collection.UpdateOne(ctx, filter_two, update_two)
	if err != nil {
		log.Println(err)
	}

	userCart_empty := make([]models.ProductUser, 0)

	filter_three := bson.D{primitive.E{Key: "_id", Value: id}}
	update_three := bson.D{{Key: "$set", Value: bson.D{primitive.E{Key: "usercart", Value: userCart_empty}}}}

	_, err = user_collection.UpdateOne(ctx, filter_three, update_three)
	if err != nil {
		log.Println(err)
		return ErrCantBuyCartItem
	}

	return nil
}

func InstantBuy(
	ctx context.Context,
	product_collection, user_collection *mongo.Collection,
	product_id primitive.ObjectID,
	user_id string) error {
	id, err := primitive.ObjectIDFromHex(user_id)
	if err != nil {
		log.Println(err)
		return ErrUserIdIsNotValid
	}

	var product_details models.ProductUser
	var orders_detail models.Order

	orders_detail.Order_Id = primitive.NewObjectID()
	orders_detail.Ordered_At = time.Now()
	orders_detail.Order_Cart = make([]models.ProductUser, 0)
	orders_detail.Payment_Method.COD = true

	err = product_collection.FindOne(ctx, bson.D{primitive.E{Key: "_id", Value: product_id}}).Decode(&product_details)
	if err != nil {
		log.Println(err)
	}

	orders_detail.Price = product_details.Price

	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "orders", Value: orders_detail}}}}

	_, err = user_collection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Println(err)
	}

	filter_two := bson.D{primitive.E{Key: "_id", Value: id}}
	update_two := bson.M{"$push": bson.M{"orders.$[].order_list": product_details}}

	_, err = user_collection.UpdateOne(ctx, filter_two, update_two)
	if err != nil {
		log.Println(err)
	}

	return nil
}
