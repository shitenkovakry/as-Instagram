package mongo

import (
	"context"
	"instagram/errordefs"
	models "instagram/models/likes"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (likes *LikesManager) findLikeByFilter(filter *bson.M) (*models.Like, error) {
	collectionLike := likes.db.Collection(likesCollection)
	result := collectionLike.FindOne(context.Background(), filter)

	err := result.Err()
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, errordefs.ErrNoDocuments
	}

	if err != nil {
		return nil, errors.Wrap(err, "can not find by ID")
	}

	var foundLike *models.Like
	if err := result.Decode(&foundLike); err != nil {
		return nil, errors.Wrap(err, "can not decode found like")
	}

	return foundLike, nil
}

func (likes *LikesManager) InsertForLike(idPhoto int, idUser int) error {
	collectionLikes := likes.db.Collection(likesCollection)

	filter := &bson.M{
		"user_id":  idUser,
		"photo_id": idPhoto,
	}

	_, err := likes.findLikeByFilter(filter)
	if err == nil {
		return nil
	}

	opts := options.InsertOne()

	likeToAdd := &models.Like{
		IDUser:  idUser,
		IdPhoto: idPhoto,
	}

	if _, err := collectionLikes.InsertOne(context.Background(), likeToAdd, opts); err != nil {
		return errors.Wrap(err, "can not insert like")
	}

	return nil
}

type counterLikes []struct {
	Count int32 `bson:"count"`
}

func (likes *LikesManager) CountLikes(idPhoto int) (int, error) {
	collectionLikes := likes.db.Collection(likesCollection)

	// Define the MongoDB pipeline stages
	pipeline := mongo.Pipeline{
		{
			{
				"$group", bson.D{
					{"_id", "$photo_id"},
					{"count",
						bson.D{
							{
								"$count", bson.D{},
							},
						},
					},
				},
			},
		},
		{
			{
				"$project", bson.D{
					{"_id", 0},
					{"count", "$count"},
				},
			},
		},
	}

	cursor, err := collectionLikes.Aggregate(context.Background(), pipeline)
	if err != nil {
		return 0, err
	}

	defer cursor.Close(context.Background())

	var myresult counterLikes

	if err := cursor.All(context.Background(), &myresult); err != nil {
		return 0, err
	}

	if len(myresult) == 1 {
		return int(myresult[0].Count), nil
	}

	return 0, errordefs.ErrIncorrect
}

func (likes *LikesManager) DeleteLike(photoID int, userID int) error {
	collectionLikes := likes.db.Collection(likesCollection)

	filter := &bson.M{
		"photo_id": photoID,
		"user_id":  userID,
	}

	_, err := likes.findLikeByFilter(filter)
	if err != nil {
		return errors.Wrap(err, "can not find deleted like")
	}

	if _, err := collectionLikes.DeleteOne(context.Background(), filter); err != nil {
		return errors.Wrap(err, "can not delete like")
	}

	return nil
}
