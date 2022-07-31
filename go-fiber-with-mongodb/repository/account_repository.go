package repository

import (
	"context"
	"fmt"

	"github.com/mandarinkb/go-fiber-with-mongodb/domain"
	"github.com/mandarinkb/go-fiber-with-mongodb/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type accountRepository struct {
	collection *mongo.Collection
}

func NewAccountRepositry(collection *mongo.Collection) domain.AccountRepository {
	return &accountRepository{collection}
}

// method that allows you to insert a single document
func (a *accountRepository) CreateOneAccount(ctx context.Context, account interface{}) error {
	result, err := a.collection.InsertOne(ctx, account)
	if err != nil {
		return err
	}

	fmt.Printf("insert success id : %v \n", result.InsertedID)
	return nil
}

// method to insert multiple documents
func (a *accountRepository) CreateManyAccount(ctx context.Context, accounts []interface{}) error {
	result, err := a.collection.InsertMany(ctx, accounts)
	if err != nil {
		return err
	}

	for _, v := range result.InsertedIDs {
		fmt.Printf("insert success id : %v \n", v)
	}

	return nil
}

// find one by _id ObjectId
func (a *accountRepository) FindOneByID(ctx context.Context, id string) (*model.Account, error) {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	selector := bson.M{"_id": objectId}
	var result model.Account
	if err := a.collection.FindOne(ctx, selector).Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

// method that returns all documents that match a search filter
func (a *accountRepository) FindByWalletID(ctx context.Context, walletID string) ([]model.Account, error) {

	selector := bson.M{"wallet_id": walletID}

	cursor, err := a.collection.Find(ctx, selector)
	if err != nil {
		return nil, err
	}

	var result []model.Account
	if err = cursor.All(context.TODO(), &result); err != nil {
		return nil, err
	}
	return result, nil
}

// method that returns only the first document that matches the filter
func (a *accountRepository) FindOneByWalletID(ctx context.Context, walletID string) (*model.Account, error) {
	selector := bson.M{"wallet_id": walletID}
	var result model.Account
	if err := a.collection.FindOne(ctx, selector).Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

// which updates the fields of a single document with a specified ObjectID
// ex updator bson.M{"$set": bson.M{}}
func (a *accountRepository) UpdateOneByID(ctx context.Context, id string, updator interface{}) error {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	result, err := a.collection.UpdateByID(context.TODO(), objectId, updator)
	if err != nil {
		return err
	}
	fmt.Println("Number of documents updated:", result.ModifiedCount)
	return nil
}

// which updates the fields of a single document
// ex updator bson.M{"$set": bson.M{}}
func (a *accountRepository) UpdateOneByWalletID(ctx context.Context, walletID string, updator interface{}) error {
	selector := bson.M{"wallet_id": walletID}
	return a.updateOne(ctx, selector, updator)
}

func (a *accountRepository) updateOne(ctx context.Context, selector interface{}, updator interface{}) error {
	result, err := a.collection.UpdateOne(ctx, selector, updator)
	if err != nil {
		return err
	}
	fmt.Println("Number of documents updated:", result.ModifiedCount)
	return nil
}

// which updates the fields of a multiple document
// ex updator bson.M{"$set": bson.M{}}
func (a *accountRepository) UpdateManyByWalletID(ctx context.Context, walletID string, updator interface{}) error {
	selector := bson.M{"wallet_id": walletID}
	result, err := a.collection.UpdateMany(ctx, selector, updator)
	if err != nil {
		return err
	}
	fmt.Println("Number of documents updated:", result.ModifiedCount)
	return nil
}

// function to overwrite the data in a document that matches a selector
func (a *accountRepository) ReplaceOneByWalletID(ctx context.Context, walletID string, replacement interface{}) error {
	selector := bson.M{"wallet_id": walletID}
	result, err := a.collection.ReplaceOne(context.TODO(), selector, replacement)
	if err != nil {
		return err
	}
	fmt.Println("Number of documents replace:", result.ModifiedCount)
	return nil
}

// ex updator bson.M{"$set": bson.M{}}
func (a *accountRepository) FindOneAndUpdateByWalletID(ctx context.Context, walletID string, updator interface{}) (*model.Account, error) {
	var result model.Account
	selector := bson.M{"wallet_id": walletID}
	upset := true
	after := options.After

	opt := options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
		Upsert:         &upset,
	}
	singleResult := a.collection.FindOneAndUpdate(ctx, selector, updator, &opt)

	if err := singleResult.Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil

}

func (a *accountRepository) FindOneAndReplace(ctx context.Context, email string, updator interface{}) (*model.Account, error) {
	var result model.Account
	selector := bson.M{"emailAddress": email}
	upset := true
	after := options.After

	opt := options.FindOneAndReplaceOptions{
		ReturnDocument: &after,
		Upsert:         &upset,
	}
	singleResult := a.collection.FindOneAndReplace(ctx, selector, updator, &opt)

	if err := singleResult.Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil

}

func (a *accountRepository) DeleteOneByWalletID(ctx context.Context, walletID string) error {
	selector := bson.M{"wallet_id": walletID}
	result, err := a.collection.DeleteOne(ctx, selector)
	if err != nil {
		return err
	}
	fmt.Println("Number of documents deleted:", result.DeletedCount)
	return nil
}

func (a *accountRepository) DeleteManyByWalletID(ctx context.Context, walletID string) error {
	selector := bson.M{"wallet_id": walletID}
	result, err := a.collection.DeleteMany(ctx, selector)
	if err != nil {
		return err
	}
	fmt.Println("Number of documents deleted:", result.DeletedCount)
	return nil
}

func (a *accountRepository) EarnPromotion(ctx context.Context, walletID string, key string, count int) (*model.Account, error) {
	var result model.Account
	selector := bson.M{"wallet_id": walletID}
	updator := bson.M{
		"$inc": bson.M{
			fmt.Sprintf("earn_promotion.privilege.%s", key): count,
		},
	}
	upset := true
	after := options.After

	opt := options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
		Upsert:         &upset,
	}
	singleResult := a.collection.FindOneAndUpdate(ctx, selector, updator, &opt)

	if err := singleResult.Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil

}
