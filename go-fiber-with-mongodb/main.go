package main

import (
	"context"
	"fmt"
	"time"

	"github.com/mandarinkb/go-fiber-with-mongodb/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func connect(uri string) (*mongo.Client, context.Context,
	context.CancelFunc, error) {
	credential := options.Credential{
		Username: "root",
		Password: "mandarinkb",
	}
	// ctx will be used to set deadline for process, here
	// deadline will of 30 seconds.
	ctx, cancel := context.WithTimeout(context.Background(),
		30*time.Second)

	// mongo.Connect return mongo.Client method
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri).SetAuth(credential))
	return client, ctx, cancel, err
}

// This is a user defined method that accepts
// mongo.Client and context.Context
// This method used to ping the mongoDB, return error if any.
func ping(client *mongo.Client, ctx context.Context) error {

	// mongo.Client has Ping to ping mongoDB, deadline of
	// the Ping method will be determined by cxt
	// Ping method return error if any occurred, then
	// the error can be handled.
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return err
	}
	fmt.Println("connected successfully")
	return nil
}

// This is a user defined method to close resources.
// This method closes mongoDB connection and cancel context.
func close(client *mongo.Client, ctx context.Context,
	cancel context.CancelFunc) {

	// CancelFunc to cancel to context
	defer cancel()

	// client provides a method to close
	// a mongoDB connection.
	defer func() {

		// client.Disconnect method also has deadline.
		// returns error if any,
		if err := client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
}
func main() {

	// Get Client, Context, CancelFunc and
	// err from connect method.
	client, ctx, cancel, err := connect("mongodb://localhost:27017")
	if err != nil {
		panic(err)
	}

	// Release resource when the main
	// function is returned.
	defer close(client, ctx, cancel)

	// Ping mongoDB with Ping method
	ping(client, ctx)

	accountCollection := client.Database("test").Collection("accounts")

	// acc, err := repository.NewAccountRepositry(accountCollection).FindOneByID(context.TODO(), "62e540d134efc63741a25fbf")
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// js, _ := json.Marshal(acc)
	// fmt.Println(string(js))

	// repository.NewAccountRepositry(accountCollection).CreateOneAccount(context.TODO(), model.MockAccount)
	// var listAccount []interface{}
	// for i := 0; i < 5; i++ {
	// 	listAccount = append(listAccount, model.MockAccount)
	// }
	// repository.NewAccountRepositry(accountCollection).CreateManyAccount(context.TODO(), listAccount)

	// update := bson.M{
	// 	"$inc": bson.M{
	// 		fmt.Sprintf("earn_promotion.privilege.%s", "AOM20220414"): 1,
	// 	},
	// }
	// repository.NewAccountRepositry(accountCollection).UpdateOneByWalletID(context.TODO(), "7a270b3b-32ff-413b-9d13-0383775fa4f8", update)

	update := bson.M{
		"$inc": bson.M{
			fmt.Sprintf("earn_promotion.privilege.%s", "TOP20220414"): 1,
		},
	}
	err = repository.NewAccountRepositry(accountCollection).UpdateManyByWalletID(context.TODO(), "51ec6f65-82a0-4b65-960b-74805cce2e15", update)
	if err != nil {
		fmt.Println(err)
	}

	// replacement := bson.M{
	// 	"firstName":    "John",
	// 	"lastName":     "Doe",
	// 	"age":          30,
	// 	"emailAddress": "johndoe@email.com",
	// }
	// err = repository.NewAccountRepositry(accountCollection).ReplaceOneByWalletID(context.TODO(), "7a270b3b", replacement)
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// updator := bson.M{"$set": bson.M{"citizen_id": "1234567890"}}
	// acc, err := repository.NewAccountRepositry(accountCollection).FindOneAndUpdateByWalletID(context.TODO(), "51ec6f65-82a0-4b65-960b-74805cce2e15", updator)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// js, _ := json.Marshal(acc)
	// fmt.Println(string(js))

	// acc, err := repository.NewAccountRepositry(accountCollection).FindOneAndReplace(context.TODO(), "johndoe@email.com", model.MockAccount)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// js, _ := json.Marshal(acc)
	// fmt.Println(string(js))

	// err = repository.NewAccountRepositry(accountCollection).DeleteManyByWalletID(context.TODO(), "7a270b3b-32ff-413b-9d13-0383775fa4f8")
	// if err != nil {
	// 	fmt.Println(err)
	// }
}
