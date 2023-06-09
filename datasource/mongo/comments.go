package mongo

import (
	"context"
	"instagram/errordefs"
	models3 "instagram/models/comments"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (comments *СommentsManager) obtainNextIDForComment() (int, error) {
	nextID := 0
	collectionComments := comments.db.Collection(commentsCollection)

	// Определение этапов агрегации.
	pipeline := []bson.M{
		{"$sort": bson.M{"comment_id": -1}},
		{"$limit": 1},
		{"$project": bson.M{"_id": 0, "comment_id": 1}},
	}

	// Создание объекта Aggregation.
	agg, err := collectionComments.Aggregate(context.Background(), pipeline)
	if err != nil {
		panic(err)
	}

	// Получение результатов агрегации.
	var resultForComment models3.Comments

	if err := agg.All(context.Background(), &resultForComment); err != nil {
		panic(err)
	}

	// Вывод результата.
	if len(resultForComment) > 0 {
		nextID = resultForComment[0].IdComment
	}

	return nextID + 1, nil
}

func (comments *СommentsManager) findCommentByFilter(filter *bson.M) (*models3.Comment, error) {
	collectionComment := comments.db.Collection(commentsCollection)
	result := collectionComment.FindOne(context.Background(), filter)

	err := result.Err()
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, errordefs.ErrNoDocuments
	}

	if err != nil {
		return nil, errors.Wrap(err, "can not find by ID")
	}

	var foundComment *models3.Comment
	if err := result.Decode(&foundComment); err != nil {
		return nil, errors.Wrap(err, "can not decode found comment")
	}

	return foundComment, nil
}

func (comments *СommentsManager) InsertForComment(commentOfUser *models3.Comment) (*models3.Comment, error) {
	collectionComments := comments.db.Collection(commentsCollection)

	nextID, err := comments.obtainNextIDForComment()
	if err != nil {
		return nil, errors.Wrap(err, "can not find next id for comment")
	}

	commentOfUser.IdComment = nextID

	opts := options.InsertOne()

	result, err := collectionComments.InsertOne(context.Background(), commentOfUser, opts)
	if err != nil {
		return nil, errors.Wrap(err, "can not insert comment")
	}

	filter := &bson.M{
		"_id": result.InsertedID,
	}

	insertedComment, err := comments.findCommentByFilter(filter)
	if err != nil {
		return nil, errors.Wrap(err, "can not find inserted comment")
	}

	return insertedComment, nil
}

func (comments *СommentsManager) DeleteComment(commentID int) (*models3.Comment, error) {
	collectionComments := comments.db.Collection(commentsCollection)

	filter := &bson.M{
		"comment_id": commentID,
	}

	deletedComment, err := comments.findCommentByFilter(filter)
	if err != nil {
		return nil, errors.Wrap(err, "can not find deleted comment")
	}

	if _, err := collectionComments.DeleteOne(context.Background(), filter); err != nil {
		return nil, errors.Wrap(err, "can not delete comment")
	}

	return deletedComment, nil
}

func (comments *СommentsManager) UpdateComment(commentID int, newComment string) (*models3.Comment, error) {
	collectionComments := comments.db.Collection(commentsCollection)

	filter := &bson.M{
		"comment_id": commentID,
	}

	upd := &bson.M{
		"$set": &bson.M{
			"comment": newComment,
		},
	}

	_, err := collectionComments.UpdateOne(context.Background(), filter, upd)
	if err != nil {
		return nil, errors.Wrap(err, "can not update comment by user")
	}

	updatedComment, err := comments.findCommentByFilter(filter)
	if err != nil {
		return nil, errors.Wrap(err, "can not find updated by comment user")
	}

	return updatedComment, nil
}

func (comments *СommentsManager) ReadComments(userID int, photoID int) (models3.Comments, error) {
	collectionComments := comments.db.Collection(commentsCollection)
	filter := &bson.M{
		"user_id":  userID,
		"photo_id": photoID,
	}

	cursor, err := collectionComments.Find(context.Background(), filter)
	if err != nil {
		return nil, errors.Wrapf(err, "can not read commentss")
	}

	var comms models3.Comments

	if err := cursor.All(context.Background(), &comms); err != nil {
		return nil, errors.Wrap(err, "can not read cursor")
	}

	return comms, nil
}
