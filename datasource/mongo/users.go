package mongo

import (
	"context"
	"instagram/errordefs"
	models "instagram/models/users"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (users *UsersManager) obtainNextID() (int, error) {
	nextID := 0
	collectionUsers := users.db.Collection(usersCollection)

	// Определение этапов агрегации.
	pipeline := []bson.M{
		{"$sort": bson.M{"user_id": -1}},
		{"$limit": 1},
		{"$project": bson.M{"_id": 0, "user_id": 1}},
	}

	// Создание объекта Aggregation.
	agg, err := collectionUsers.Aggregate(context.Background(), pipeline)
	if err != nil {
		panic(err)
	}

	// Получение результатов агрегации.
	var result models.UsersRegistration

	if err := agg.All(context.Background(), &result); err != nil {
		panic(err)
	}

	// Вывод результата.
	if len(result) > 0 {
		nextID = result[0].ID
	}

	return nextID + 1, nil
}

func (users *UsersManager) findUserByFilter(filter *bson.M) (*models.UserRegistration, error) {
	collectionUsers := users.db.Collection(usersCollection)
	result := collectionUsers.FindOne(context.Background(), filter)

	err := result.Err()
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, errordefs.ErrNoDocuments
	}

	if err != nil {
		return nil, errors.Wrap(err, "can not find by ID")
	}

	var foundUser *models.UserRegistration
	if err := result.Decode(&foundUser); err != nil {
		return nil, errors.Wrap(err, "can not decode found user")
	}

	return foundUser, nil
}

func (users *UsersManager) Insert(user *models.UserRegistration) (*models.UserRegistration, error) {
	collectionUsers := users.db.Collection(usersCollection)

	nextID, err := users.obtainNextID()
	if err != nil {
		panic(err)
	}

	user.ID = nextID

	opts := options.InsertOne()

	result, err := collectionUsers.InsertOne(context.Background(), user, opts)
	if err != nil {
		return nil, errors.Wrap(err, "can not insert user")
	}

	filter := &bson.M{
		"_id": result.InsertedID,
	}

	insertedUser, err := users.findUserByFilter(filter)
	if err != nil {
		return nil, errors.Wrap(err, "can not find inserted user")
	}

	return insertedUser, nil
}

func (users *UsersManager) UpdateName(userID int, newName string) (*models.UserRegistration, error) {
	collectionUsers := users.db.Collection(usersCollection)

	filter := &bson.M{
		"user_id": userID,
	}

	upd := &bson.M{
		"$set": &bson.M{
			"name": newName,
		},
	}

	_, err := collectionUsers.UpdateOne(context.Background(), filter, upd)
	if err != nil {
		return nil, errors.Wrap(err, "can not update name by user")
	}

	updatedUser, err := users.findUserByFilter(filter)
	if err != nil {
		return nil, errors.Wrap(err, "can not find updated by name user")
	}

	return updatedUser, nil
}

func (users *UsersManager) UpdateEmail(userID int, newEmail string) (*models.UserRegistration, error) {
	collectionUsers := users.db.Collection(usersCollection)

	filter := &bson.M{
		"user_id": userID,
	}

	upd := &bson.M{
		"$set": &bson.M{
			"email": newEmail,
		},
	}

	_, err := collectionUsers.UpdateOne(context.Background(), filter, upd)
	if err != nil {
		return nil, errors.Wrap(err, "can not update email by user")
	}

	updatedUser, err := users.findUserByFilter(filter)
	if err != nil {
		return nil, errors.Wrap(err, "can not find updated by email user")
	}

	return updatedUser, nil
}

func (users *UsersManager) DeleteUser(userID int) (*models.UserRegistration, error) {
	collectionUsers := users.db.Collection(usersCollection)

	filter := &bson.M{
		"user_id": userID,
	}

	deletedUser, err := users.findUserByFilter(filter)
	if err != nil {
		return nil, errors.Wrap(err, "can not find deleted user")
	}

	if _, err := collectionUsers.DeleteOne(context.Background(), filter); err != nil {
		return nil, errors.Wrap(err, "can not delete user")
	}

	collectionComments := users.db.Collection(commentsCollection)

	if _, err := collectionComments.DeleteMany(context.Background(), filter); err != nil {
		return nil, errors.Wrap(err, "can not delete user's comments")
	}

	collectionPhotos := users.db.Collection(photosCollection)

	if _, err := collectionPhotos.DeleteMany(context.Background(), filter); err != nil {
		return nil, errors.Wrap(err, "can not delete user's photos")
	}

	collectionLikes := users.db.Collection(likesCollection)

	if _, err := collectionLikes.DeleteMany(context.Background(), filter); err != nil {
		return nil, errors.Wrap(err, "can not delete user's likes")
	}

	return deletedUser, nil
}
